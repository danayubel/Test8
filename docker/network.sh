#!/bin/bash

# Detener y limpiar contenedores
docker-compose down -v
rm -rf ./ledgers-backup

# Iniciar red
docker-compose up -d
