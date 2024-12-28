#!/bin/bash

# Limpiar artefactos existentes
rm -rf crypto-config channel-artifacts/*.block channel-artifacts/*.tx

# Generar certificados
cryptogen generate --config=./crypto-config/config.yaml --output=./crypto-config

# Crear bloque génesis
configtxgen -profile SystemChannel -channelID system-channel -outputBlock ./channel-artifacts/genesis.block

# Crear transacción de canal
configtxgen -profile BasicChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID mychannel
