# Crypto Voting Platform

Live ranking system where the user can add cryptos, upvote your favorates and downvote the other ones.

## Tech Used

- Golang
- MongoDB
- gRPC

## Endpoints

### POST /crypto

#### Request

```bash
curl -XPOST http://ec2-54-89-58-88.compute-1.amazonaws.com:8080/crypto --header "Content-Type: application/json"  --data '{
    "name": "Ethereum",
}'
```

#### Response

```json
{
  "id": "5fd6ac5c6884b412d6ec1475",
  "name": "Ethereum"
}
```

### GET /crypto/:id

#### Request

```bash
curl http://ec2-54-89-58-88.compute-1.amazonaws.com:8080/crypto/6057708661402e9c75b80109
```

#### Response

```json
{
  "id": "6057708661402e9c75b80109",
  "name": "Bitcoin",
  "upvotes": 3,
  "downvotes": 1,
  "score": 2
}
```
### POST /upvote/:id

#### Request

```bash
curl -XPOST http://ec2-54-89-58-88.compute-1.amazonaws.com:8080/upvote/6057708661402e9c75b80109 
```

#### Response

```json
{
  "success": true
}
```


### POST /downvote/:id

#### Request

```bash
curl -XPOST http://ec2-54-89-58-88.compute-1.amazonaws.com:8080/downvote/6057708661402e9c75b80109 
```

#### Response

```json
{
  "success": true
}
```

### GET /cryptos

#### Request

```bash
curl http://ec2-54-89-58-88.compute-1.amazonaws.com:8080/cryptos
```

#### Response

```json
{
  "id": "6057708661402e9c75b80109",
  "name": "Bitcoin",
  "upvotes": 2,
  "downvotes": 1,
  "score": 1
}{
    "id": "6057708f61402e9c75b8010a",
    "name": "Ethereum",
    "upvotes": 12,
    "downvotes": 2,
    "score": 10
}{
    "id": "6057709461402e9c75b8010b",
    "name": "Litcoin",
    "upvotes": 10,
    "downvotes": 1,
    "score": 9
}
```
### DELETE /delete/:id

#### Request

```bash
curl http://ec2-54-89-58-88.compute-1.amazonaws.com:8080/delete/6057708661402e9c75b80109 
```

#### Response

```json
{
  "success": true
}
```