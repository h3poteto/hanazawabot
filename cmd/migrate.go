package cmd

import (
	"log"
	"regexp"

	"../models/basedb"

	"github.com/spf13/cobra"
)

func migrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migration database",
		Run:   migrate,
	}

	return cmd
}

func migrate(cmd *cobra.Command, args []string) {
	db := basedb.SharedInstance().Connection
	defer basedb.SharedInstance().Close()

	sqlFiles, err := AssetDir("cmd/migrate")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range sqlFiles {
		if match, _ := regexp.MatchString(".sql$", file); match {
			log.Printf("migrate execute: %s ...", file)
			byteQuery, err := Asset("cmd/migrate/" + file)
			if err != nil {
				log.Fatal(err)
			}
			query := string(byteQuery)
			if len(query) > 0 {
				log.Println(query)
				if _, err := db.Exec(query); err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}
