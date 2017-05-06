package main

import (
	"github.com/h3poteto/hanazawabot/cmd"

	"fmt"
	"os"
)

//go:generate go-bindata -pkg=cmd -o=cmd/bindata.go cmd/migrate/

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
