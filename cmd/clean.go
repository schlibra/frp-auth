package cmd

import (
	"fmt"
	"frp-auth/constant"
	"frp-auth/utils"
	"log"

	"github.com/spf13/cobra"
)

var CleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove all databases(!!!Dangerous!!!)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Loading configuration...")
		cfg := utils.LoadConfig()
		if cfg.Dangerous.AllowClean {
			fmt.Println("Loading database...")
			db := utils.LoadMySQL(cfg)
			fmt.Println("Dropping user table...")
			result, err := db.Exec(constant.DropUserTable)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(result)
			fmt.Println("Dropping token table...")
			result, err = db.Exec(constant.DropTokenTable)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(result)
			fmt.Println("Dropping port table...")
			result, err = db.Exec(constant.DropPortTable)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(result)
		} else {
			fmt.Println("Remove database is not allowed!")
		}
	},
}
