# Projeto de Leilão com Go

Este projeto implementa um sistema de leilões com Go, utilizando MongoDB como banco de dados.

## Funcionalidades

- Criação de leilões
- Realização de lances
- Fechamento automático de leilões expirados

## Como executar o projeto

### Pré-requisitos

- Docker
- Docker Compose

### Passos

1. **Clone o repositório:**

   ```bash
   git clone https://github.com/evandrojr/labs-auction-goexpert.git
   ```

2. **Inicie os serviços com Docker Compose:**

   ```bash
   docker compose up -d
   ```
   ou 

   ```bash
   docker-compose up -d
   ```

   Este comando irá iniciar a aplicação Go e o banco de dados MongoDB.

3. **Acesse a aplicação:**

   A API estará disponível em `http://localhost:8080`.

## Como executar os testes

1. **Inicie o banco de dados de teste:**

   ```bash
   docker compose up -d mongodb
   ```

2. **Execute os testes:**

   ```bash
   go test ./...
   ```
