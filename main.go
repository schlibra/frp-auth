package main

import (
	"embed"
	"frp-auth/cmd"

	_ "github.com/go-sql-driver/mysql"
)

//go:embed frontend/dist
var distFS embed.FS

func main() {
	cmd.Execute(distFS)
}
