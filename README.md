# CMPE273-Assignment1

TO MAKE RPC CALLS USING curl :-

- Build the server.
go build Server.go
- Run the Server. 
go run Server.go

- Run curl on another terminal. For buying stock
curl  -H "Content-Type: application/json"  -d '{"method":"StockService.Say","params":[{"StockSymbolandPercentage":"GOOG:50%,YHOO:50%","Budget":7224.4}], "id":0}' http://localhost:8080/rpc

- For checking Portfolio
curl  -H "Content-Type: application/json"  -d '{"method":"CheckPortfolioService.Che","params":[{"Tradeid":1}], "id":0}' http://localhost:8080/rpc
 
TO MAKE RPC CALLS USING Client.go :-

- Build the server.
go build Server.go
- Run the server
go run Server.go

- Build the client
go build Client.go

- Run client with proper commanline arguments.
2 arguments will call buyingstock service
./Client GOOG:50%,YHOO:50% 7224.4

1 argument will call checking portfolio service
./Client 1 



