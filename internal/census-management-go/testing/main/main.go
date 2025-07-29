package main

import (
	"net"
	"net/http"

	cmt "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go/testing"
)

func main() {
	server, serverErr := cmt.NewCensusManagementServer()
	if serverErr != nil {
		panic(serverErr)
	}

	listener, listenerErr := net.Listen("tcp", ":0")
	if listenerErr != nil {
		panic(listenerErr)
	}

	defer listener.Close()

	print("census-management/testing/server: ", listener.Addr().String(), "\n")

	err := http.Serve(listener, server)
	if err != nil {
		panic(err)
	}
}
