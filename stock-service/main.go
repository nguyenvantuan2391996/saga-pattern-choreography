package main

import (
	"fmt"
	"os"
	cmd "saga-pattern-choreography/stock-service/grpc"
)

func main() {
	if err := cmd.RunServerCMD(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
