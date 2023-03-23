package main

import (
	"flag"
	"log"
	"os"

	"github.com/mgrabazey/tic-tac-toe/internal/api/transport/http"
	"github.com/mgrabazey/tic-tac-toe/internal/app/module/game"
	"github.com/mgrabazey/tic-tac-toe/internal/pkg/postgres"
	"github.com/mgrabazey/tic-tac-toe/internal/service/migration"
	"github.com/mgrabazey/tic-tac-toe/internal/service/repo"
)

func main() {
	var (
		publicUrl  string
		dbUser     string
		dbPassword string
		dbHost     string
		dbPort     string
		dbName     string
	)
	flag.StringVar(&publicUrl, "public-url", os.Getenv("PUBLIC_URL"), "Public URL of the API")
	flag.StringVar(&dbUser, "db-user", os.Getenv("DB_USER"), "Database user")
	flag.StringVar(&dbPassword, "db-password", os.Getenv("DB_PASSWORD"), "Database password")
	flag.StringVar(&dbHost, "db-host", os.Getenv("DB_HOST"), "Database host")
	flag.StringVar(&dbPort, "db-port", os.Getenv("DB_PORT"), "Database port")
	flag.StringVar(&dbName, "db-name", os.Getenv("DB_NAME"), "Database name")
	flag.Parse()

	if publicUrl == "" || dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" || dbName == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	db, err := postgres.Conn(&postgres.Config{
		User:     dbUser,
		Password: dbPassword,
		Host:     dbHost,
		Port:     dbPort,
		Name:     dbName,
	})
	if err != nil {
		log.Fatalf("unable to connect to database: %v\n", err)
	}

	// Run database migrations.
	err = migration.Run(db)
	if err != nil {
		log.Fatalf("unable to run service: %v\n", err)
	}

	// Tun HTTP application
	err = httpx.Run(publicUrl, game.NewService(repo.NewGameRepository(db)))
	if err != nil {
		log.Fatalf("unable to run service: %v\n", err)
	}
}
