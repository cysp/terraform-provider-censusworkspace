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

	ln, lnErr := net.Listen("tcp", ":0")
	if lnErr != nil {
		panic(lnErr)
	}

	print("Server is running at: ", ln.Addr().String(), "\n")

	defer ln.Close()

	err := http.Serve(ln, server.Server())
	if err != nil {
		panic(err)
	}
}
