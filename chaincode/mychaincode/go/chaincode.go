package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract gestiona las transacciones de energía.
type SmartContract struct {
	contractapi.Contract
}

// EnergyTransaction describe una transacción de energía.
type EnergyTransaction struct {
	Sender     string `json:"sender"`     // Quién envía la energía.
	Receiver   string `json:"receiver"`   // Quién recibe la energía.
	EnergyType string `json:"energyType"` // Tipo de energía (solar, eólica, etc.).
	Tokens     int    `json:"tokens"`     // Cantidad de tokens.
}

// UserBalance gestiona el balance de tokens de un usuario.
type UserBalance struct {
	Balance int `json:"balance"` // Balance de tokens.
}

// SetTransaction registra una transacción de energía entre dos usuarios.
func (s *SmartContract) SetTransaction(ctx contractapi.TransactionContextInterface, sender string, receiver string, energyType string, kwh int) error {
	// Validar roles usando el MSP del remitente
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("No se pudo obtener el MSPID del remitente: %s", err.Error())
	}

	if clientMSPID != "SellerMSP" && clientMSPID != "ClientMSP" {
		return fmt.Errorf("Solo vendedores o clientes pueden iniciar una transacción de venta de energía")
	}

	// Calcular el costo en tokens (1 kWh = 2 tokens)
	tokens := kwh * 2

	// Obtener el balance actual del remitente
	senderBalance, err := s.getBalance(ctx, sender)
	if err != nil {
		return fmt.Errorf("Error al obtener el balance del remitente: %s", err.Error())
	}

	// Verificar si el remitente tiene suficientes tokens
	if senderBalance.Balance < tokens {
		return fmt.Errorf("El remitente %s no tiene suficientes tokens", sender)
	}

	// Obtener el balance del receptor
	receiverBalance, err := s.getBalance(ctx, receiver)
	if err != nil {
		return fmt.Errorf("Error al obtener el balance del receptor: %s", err.Error())
	}

	// Actualizar balances
	senderBalance.Balance -= tokens
	receiverBalance.Balance += tokens

	// Guardar balances actualizados
	if err := s.updateBalance(ctx, sender, senderBalance); err != nil {
		return fmt.Errorf("Error al actualizar el balance del remitente: %s", err.Error())
	}
	if err := s.updateBalance(ctx, receiver, receiverBalance); err != nil {
		return fmt.Errorf("Error al actualizar el balance del receptor: %s", err.Error())
	}

	// Registrar la transacción en el ledger
	transaction := EnergyTransaction{
		Sender:     sender,
		Receiver:   receiver,
		EnergyType: energyType,
		Tokens:     tokens,
	}
	transactionBytes, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("Error al empaquetar la transacción de energía: %s", err.Error())
	}

	if err := ctx.GetStub().PutState(ctx.GetStub().GetTxID(), transactionBytes); err != nil {
		return fmt.Errorf("Error al registrar la transacción en el ledger: %s", err.Error())
	}

	return nil
}

// getBalance obtiene el balance de un usuario.
func (s *SmartContract) getBalance(ctx contractapi.TransactionContextInterface, user string) (*UserBalance, error) {
	balanceBytes, err := ctx.GetStub().GetState(user)
	if err != nil {
		return nil, fmt.Errorf("No se pudo obtener el estado del usuario: %s", err.Error())
	}
	if balanceBytes == nil {
		return nil, fmt.Errorf("El usuario %s no existe", user)
	}

	var balance UserBalance
	if err := json.Unmarshal(balanceBytes, &balance); err != nil {
		return nil, fmt.Errorf("Error al desempaquetar el balance del usuario: %s", err.Error())
	}

	return &balance, nil
}

// updateBalance actualiza el balance de un usuario en el ledger.
func (s *SmartContract) updateBalance(ctx contractapi.TransactionContextInterface, user string, balance *UserBalance) error {
	balanceBytes, err := json.Marshal(balance)
	if err != nil {
		return fmt.Errorf("Error al empaquetar el balance: %s", err.Error())
	}

	if err := ctx.GetStub().PutState(user, balanceBytes); err != nil {
		return fmt.Errorf("Error al guardar el balance en el ledger: %s", err.Error())
	}

	return nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		fmt.Printf("Error al crear el chaincode: %s\n", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error al iniciar el chaincode: %s\n", err.Error())
	}
}
