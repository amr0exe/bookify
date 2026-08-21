package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/amr0exe/bookify/internal/db/models"
)

func main() {
	stmts, err := gormschema.New("postgres").
		Load(
			&models.Account{},
			&models.Consumer{},
			&models.Business{},
			&models.RefreshToken{},
			&models.Service{},
			&models.Appointment{},
		)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load schema: %v\n", err)
		os.Exit(1)
	}

	io.WriteString(os.Stdout, stmts)
}
