package main

import "github.com/bnb-chain/token-recover-approver/cmd"

//go:generate wire ./...

func main() {
	cmd.Execute()
}
