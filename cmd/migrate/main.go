package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/dockloud/document/internal/infrastructure/config"
	"github.com/dockloud/document/internal/infrastructure/database/postgres/migrations"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("usage: migrate <up|down|status>")
	}

	conf, err := config.LoadMergedConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", conf.Database.URL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}

	goose.SetBaseFS(migrations.FS)

	switch os.Args[1] {
	case "up":
		err = goose.Up(db, ".")
	case "down":
		err = goose.Down(db, ".")
	case "status":
		err = goose.Status(db, ".")
	default:
		log.Fatalf("unknown command: %s", os.Args[1])
	}

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Migration command completed")
}
