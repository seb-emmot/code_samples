## Tollcalculator implementation in Go.

either run the unit tests in tollcalculator/ \
`go run test .`

or run the http server \
`go run cmd/server/server.go`

test the api with
`curl -X POST http://localhost:8080/tolls -H "Content-Type: application/json" -d '{"passes": ["2024-09-23T15:04:05Z"], "vehicletype": 6}`
