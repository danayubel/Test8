#!/bin/bash

# Detener y eliminar contenedores, redes y volúmenes
docker-compose down -v
docker system prune -f
docker volume prune -f
