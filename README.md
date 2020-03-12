# Stone Challenge

## Requisitos originais do desafio
**O desafio é criar uma API de transferencia entre contas Internas de um banco.**

=====================
#### API

#### Regras gerais

* Usar formato JSON para leitura e escrita. (ex: `GET /accounts/` retorna json, `POST /accounts/ {name: 'james bond'}`)

#### Rotas esperadas

##### `/accounts/`

A entidade `Account` possui os seguintes atributos:

* `id`
* `name` 
* `cpf` 
* `ballance` 
* `created_at` 

Espera-se as seguintes ações:

- `GET /accounts` - obtém a lista de contas
- `GET /accounts/{account_id}/ballance` - obtém o saldo da conta
- `POST /accounts` - cria um `Account`

*Regras para esta rota*

- `ballance` pode iniciar com algum valor para simplificar 

* * *

##### `/transfers/`

A entidade `Transfer` possui os seguintes atributos:

* `id`
* `account_origin_id`
* `account_destination_id`
* `amount`
* `created_at`

Espera-se as seguintes ações:

- `GET /transfers` - obtém a lista de transferencias
- `POST /transfers` - faz transferencia de um `Account` para outro.

*Regras para esta rota*

- Caso `Account` de origem no tenha saldo, retornar um código de erro apropriado
- Atualizar o `ballance` das contas

# Sobre este Banco

- O código e sua documentação estão em inglês, seguindo o padrão das entidades

## Vantagens desse Banco incrível 💰
- Aqui não há limite de crédito. Ou seja: você não se individa à toa!
- Nossos correntistas são protegidos contra transferências duplicadas
- Você pode nos confiar seu dinheiro desde a criação da sua conta 🥳
- Mantemos um histórico de todas as requisições de transferências de nossos correntistas para fins de compliance 🧮

## Como rodar
`go build -o app cmd/main.go && ./app`

Uma mensagem similar a essa deve aparecer:

`2020/03/12 18:04:39 initializing server on :3000`

## Endpoint /accounts
###### POST
`POST http://localhost:3000/accounts
Content-Type: application/json`

- Exemplo de request:
```json
{
  "name": "Kevin Malone",
  "cpf": "66648111038",
  "balance": 2000
}
```
- Retornos possíveis:
  - Sucesso: `201 Created`
  ```json
  {
    "id": 1
  }
  ```
  - Insucesso: `400 Bad Request`, `500 Internal Server Error`
  
###### GET

`GET http://localhost:3000/accounts`

- Retornos possíveis:
  - Sucesso: `200 OK`
  ```json
  [
     {
        "id":1,
        "name":"Kevin Malone",
        "cpf":"66648111038",
        "balance":2000,
        "created_at":"2020-03-12T16:58:34.267575763-03:00"
     }
  ]
  ```
  - Insucesso: `500 Internal Server Error`

## Endpoint /accounts/{account_id}/balance

`GET http://localhost:3000/accounts/1/balance`

- Retornos possíveis:
  - Sucesso: `200 OK`
  ```json
  {
    "id": 1,
    "balance": 2000
  }
  ```
  - Insucesso: `400 Bad Request`, `404 Not Found`, `500 Internal Server Error`

## Endpoint /transfers

###### POST
`POST http://localhost:3000/transfers
 Content-Type: application/json`

- Exemplo de request:
```json
{
  "account_origin_id": 1,
  "account_destination_id": 2,
  "amount": 1000
}
```
- Retornos possíveis:
  - Sucesso: `201 Created`
  ```json
  {
    "id": 1
  }
  ```
  - Insucesso: `400 Bad Request`, `500 Internal Server Error`

###### GET
`GET http://localhost:3000/transfers`

- Retornos possíveis:
  - Sucesso: `200 OK`
  ```json
  [
     {
         "id": 1,
         "account_origin_id": 1,
         "account_destination_id": 2,
         "amount": 1000,
         "created_at": "2020-03-12T17:04:42.911774963-03:00",
         "status": "Confirmed"
       }
  ]
  ```

  - Insucesso: `500 Internal Server Error`

## Endpoint /transfers/{transfer_id}

`GET http://localhost:3000/transfers/1`

- Retornos possíveis:
  - Sucesso: `200 OK`
  ```json
  {
      "id": 1,
      "account_origin_id": 1,
      "account_destination_id": 2,
      "amount": 1000,
      "created_at": "2020-03-12T17:04:42.911774963-03:00",
      "status": "Confirmed"
  }
  ```
  - Insucesso: `400 Bad Request`, `404 Not Found`, `500 Internal Server Error`

## Regras
- Todos os valores de `balance` e `amount` são representados em centavos
- Não é possível efetuar transferências:
  - Caso a conta de origem não tenha `balance` suficiente para transferir
  - Caso o `account_origin_id` e o `account_destination_id` informados sejam iguais
  - Caso a requisição da transferência tenha mesmos `account_origin_id`, `account_destination_id` e `amount` que uma transferência com status `Confirmed` que tenha acontecido em 10 segundos ou menos
  - Caso o `amount` indicado seja 0
- Todos os requests de criação de transferência criam registros, para futuras auditorias. Só não criarão registro as requisições que tiverem `account_origin_id` e `account_destination_id` que não existem no Banco
- As contas precisam ser criadas com um valor de `balance`, sempre igual ou maior a 0
- O `cpf` informado precisa ter 11 caracteres, todos numéricos

🤓