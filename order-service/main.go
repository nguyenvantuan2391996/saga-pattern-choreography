package main

import (
	"fmt"
	"os"
	cmd "saga-pattern-choreography/payment-service/grpc"
)

func main() {
	if err := cmd.RunServerCMD(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
