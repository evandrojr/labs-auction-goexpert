#!/bin/bash

# Executa os testes usando Docker Compose
echo "Executando os testes..."
docker compose -f docker-compose.test.yml up --build --abort-on-container-exit

if [ $? -ne 0 ]; then
  echo "Os testes falharam. A aplicação não será iniciada."
  docker compose -f docker-compose.test.yml down
  exit 1
fi

echo "Testes executados com sucesso!"
docker compose -f docker-compose.test.yml down

# Inicia a aplicação e o banco de dados usando Docker Compose
echo "Iniciando os serviços..."
docker compose up -d --build
echo "Serviços iniciados com sucesso!"
