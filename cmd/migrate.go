package cmd

import (
	"bytes"
	"io"
	"log"
	"regexp"

	"github.com/h3poteto/hanazawabot/models/basedb"

	"github.com/pkg/errors"
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

	if len(Assets.Dirs["/cmd/migrate"]) < 1 {
		log.Fatal(errors.New("migration files does not exist"))
	}

	for _, file := range Assets.Dirs["/cmd/migrate"] {
		if match, _ := regexp.MatchString(".sql$", file); match {
			log.Printf("migrate execute: %s ...", file)
			f, err := Assets.Open("/cmd/migrate/" + file)
			if err != nil {
				log.Fatal(err)
			}
			byteQuery := new(bytes.Buffer)
			io.Copy(byteQuery, f)
			query := string(byteQuery.Bytes())
			if len(query) > 0 {
				log.Println(query)
				if _, err := db.Exec(query); err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}
