package cmd

import (
	"fmt"
	"frp-auth/constant"
	"frp-auth/utils"
	"log"

	"github.com/spf13/cobra"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the database",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Loading config...")
		cfg := utils.LoadConfig()
		fmt.Println("Loading database...")
		db := utils.LoadMySQL(cfg)
		fmt.Println("Creating user table...")
		result, err := db.Exec(constant.CreateUserTable)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
		fmt.Println("Creating token table...")
		result, err = db.Exec(constant.CreateTokenTable)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
		fmt.Println("Creating port table...")
		result, err = db.Exec(constant.CreatePortTable)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
		createAdmin()
	},
}
