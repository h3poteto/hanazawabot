package cmd

import (
	"log"

	"github.com/h3poteto/hanazawabot/models/basedb"

	"github.com/spf13/cobra"
)

func seedCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "seed",
		Short: "Initial dataset",
		Run:   seed,
	}

	return cmd
}

func seed(cmd *cobra.Command, args []string) {
	db := basedb.SharedInstance().Connection
	defer basedb.SharedInstance().Close()

	_, err := db.Exec("truncate table serifs;")
	if err != nil {
		log.Fatalf("mysql execute error: %v", err)
	}
	_, err = db.Exec("insert into serifs (body, created_at) values (?, now()), (?, now()), (?, now()), (?, now()), (?, now()), (?, now()), (?, now()), (?, now())",
		"うん、うん、これですよ。このふともも・・・",
		"オチ？オチなんてないよぉ！",
		"私のこの楽しい話を邪魔しないでくれるっ！",
		"あたしもう25歳だよ！お前あたしと結婚できんのか！",
		"そ、そんな事言うと脱ぐぞ",
		"ディスってんのかぁぁっ",
		"ねぇじゃねーよ！",
		"弟紹介しようか？")
	if err != nil {
		log.Fatalf("mysql execute error: %v", err)
	}
}
