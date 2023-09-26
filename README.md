# Transaction API

Transaction API is a challenge proposed by the company Wex, which allows users to store and retrieve purchase transactions.

## Requirement

- Go 1.21+

## Getting Started

1. Clone the repository.
2. Go to the root of the project.

### Use Dockerfile

```
docker build -t app .
docker run -p 8082:8082 --rm app
```

## API documentation

- [Create a transaction](#create-a-transaction)
- [Get a transaction](#get-a-transaction)
- [Detailed documentation](#detailed-documentation)

### Create a transaction

`[POST] /transactions`

#### cURL example

```
curl -X POST -H "Content-Type: application/json" -d '{
  "description": "some transaction",
  "transaction_date": "2023-09-01T12:00:00Z",
  "amount": 100.50
}' http://localhost:8082/v1/transactions
```

### Get a transaction

`[GET] /transactions/{id}?country={country}&currency={currency}`

#### cURL example

```
curl -X GET \
  "http://localhost:8082/v1/transactions/9b25d3e4-dfc0-45d8-b600-0920c9c00c43?country=Canada&currency=Dollar"

```

### Detailed documentation

Please check: [link](https://vickiliou.github.io/challenge-wex/swagger.html)
