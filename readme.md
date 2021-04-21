# FetchRewards Points API
This project is a web API used for managing users' point balances while keeping track of which partner's points are spent at specific times.

## Building
Building this project requires `Go` to be installed. Please refer to the documentation [here](https://golang.org/doc/install) for installing the latest version of `Go`. This project uses 1.15 but the latest version should work fine.

You can run this project two different ways, build the binary and run the binary or just hand the main folder to the `go` command.

### Running from binary
```
go build cmd/api
./api
```
### Running with the `go` commaand
```
go run cmd/api
```
#### Specifying the port
You can specify the port in which the server listens by providing an additional integer argument in the run command. The default port is 8090.
```
go run cmd/api 9090
```

## Testing
Run the following command from the project root to run the tests.
```
go test ./...
```
### Local endpoint testing
Once you start the server you may want to test it out. The following commands exercise most of the API's functionality. You can use these as a starting point.
#### Initialize transactions
```
curl -X POST \
  http://localhost:8090/v1/users/1/points/add \
  -d '{ "payer": "DANNON", "points": 1000, "timestamp": "2020-11-02T14:00:00Z" }'

curl -X POST \
  http://localhost:8090/v1/users/1/points/add \
  -d '{ "payer": "UNILEVER", "points": 200, "timestamp": "2020-10-31T11:00:00Z" }'

curl -X POST \
  http://localhost:8090/v1/users/1/points/add \
  -d '{ "payer": "DANNON", "points": -200, "timestamp": "2020-10-31T15:00:00Z" }'

curl -X POST \
  http://localhost:8090/v1/users/1/points/add \
  -d '{ "payer": "MILLER COORS", "points": 10000, "timestamp": "2020-11-01T14:00:00Z" }'

curl -X POST \
  http://localhost:8090/v1/users/1/points/add \
  -d '{ "payer": "DANNON", "points": 300, "timestamp": "2020-10-31T10:00:00Z" }'  
```

#### Spend points
```
curl -X POST \
  http://localhost:8090/v1/users/1/points/spend \
  -d '{
	"points": 5000
}'
```

#### Get payer balances
```
curl -X GET \
  http://localhost:8090/v1/users/1/payers

```