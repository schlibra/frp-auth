package cmd

import (
	"database/sql"
	"fmt"
	"frp-auth/constant"
	"frp-auth/model"
	"frp-auth/utils"
	"log"

	"github.com/spf13/cobra"
)

func createAdmin() {
	fmt.Println("Loading config...")
	cfg := utils.LoadConfig()
	fmt.Println("Connecting to database...")
	db := utils.LoadMySQL(cfg)
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	adminPassword, err := utils.GenerateRandomString(6)
	if err != nil {
		log.Fatal(err)
	}
	passwordHash, err := utils.PasswordHash(adminPassword)
	var user model.UserTable
	if err := db.QueryRow(constant.SelectUserTableByUsername, "admin").Scan(&user.ID, &user.Username, &user.Password, &user.Nickname, &user.Admin, &user.Enable, &user.Admin); err == nil {
		if _, err = db.Exec(constant.UpdateUserTableById, passwordHash, user.Nickname, true, true, user.ID); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Admin password reset to: " + adminPassword)
	} else {
		if _, err = db.Exec(constant.InsertUserTable, "admin", adminPassword, "Admin User", true, true); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Created Admin user with password: " + adminPassword)
	}
}

var AdminCmd = &cobra.Command{
	Use:   "admin",
	Short: "create admin user",
	Run: func(cmd *cobra.Command, args []string) {
		createAdmin()
	},
}
