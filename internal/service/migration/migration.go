package migration

import (
	"database/sql"
	"log"
)

type migration interface {
	name() string
	up(tx *sql.Tx) error
}

// Run runs database migrations.
func Run(db *sql.DB) error {
	log.Printf("Run migrations...\n")
	err := prepare(db)
	if err != nil {
		return err
	}

	m, err := load(db)
	if err != nil {
		return err
	}

	for _, i := range []migration{
		&createGamesTable{},
	} {
		if m[i.name()] {
			continue
		}
		log.Printf("Run %q migration...\n", i.name())
		err = up(db, i)
		if err != nil {
			return err
		}
	}
	log.Printf("Database is up to date.\n")
	return nil
}

func prepare(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS "migrations" (name varchar(255) not null)`)
	if err != nil {
		return err
	}
	return nil
}

func load(db *sql.DB) (map[string]bool, error) {
	s, err := db.Query(`SELECT * FROM "migrations"`)
	if err != nil {
		return nil, err
	}
	defer s.Close()

	m := make(map[string]bool)
	for s.Next() {
		var n string
		err = s.Scan(&n)
		if err != nil {
			return nil, err
		}
		m[n] = true
	}
	return m, nil
}

func up(db *sql.DB, migration migration) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			err = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()
	err = migration.up(tx)
	if err != nil {
		return
	}
	_, err = tx.Exec(`INSERT INTO "migrations" ("name") VALUES ($1)`, migration.name())
	if err != nil {
		return
	}
	return
}
