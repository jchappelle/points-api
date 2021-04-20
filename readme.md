# FetchRewards Points API
This project is a web API used for managing users' point balances while keeping track of which partner's points are spent at specific times.

## Building
Building this project requires `Go` to be installed. Please refer to the documentation here for installing the latest version of `Go`. This project uses 1.15 but the latest version should work fine.

You can run this project two different ways, build the binary and run the binary or just hand the main folder to the `go` command.

### Running from binary
```
go build cmd/api
./api
```
### Running with `go` commaand
```
go run cmd/api
```
#### Specifying the port
You can specify the port in which the server listens by providing an additional integer argument in the run command. The default port is 8090.
```
go run cmd/api 9090
```

## Testing
```
go test ./...
```