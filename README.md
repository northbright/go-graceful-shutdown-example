# go-graceful-shutdown-example

Example of graceful shutdown HTTP server.

## Shutdown the HTTP server
* It follows the example of [Shutdown](https://godoc.org/net/http#Server.Shutdown) to start a goroutine to shutdown server.
* It'll stop all worker goroutines after [os.Interrupt received](https://godoc.org/os/signal#Notify).

## Usage
* Run `go run main.go` to start the server.
* Open one or more browser tabs to visit `http://localhost:8080`.
* Check the terminal output:
  * It uses a 10-second timeout to emulate a long time work for each request.
  * Try to press `Control-C` in the terminal to send SIGINT before or after the work finished(timeout).
