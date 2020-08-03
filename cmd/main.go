package main

import (
	"fmt"

	"github.com/congcongke/httpfileserver/cmd/httpserver"
)

func main() {
	fmt.Println("start http file server")

	cmd := httpserver.NewHttpFileServerCommand()

	cmd.Execute()
}
