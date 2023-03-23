package migration

import "database/sql"

type createGamesTable struct{}

func (m *createGamesTable) name() string {
	return "20230321_193200_create_games_table"
}

func (m *createGamesTable) up(tx *sql.Tx) error {
	_, err := tx.Exec(`CREATE TABLE "games"
(
    "id" UUID PRIMARY KEY,
    "board" CHAR(9) NOT NULL,
    "status" VARCHAR(16)  NOT NULL,
    "char" CHAR(1)  NOT NULL,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL,
    "deleted_at" TIMESTAMP
)`)
	if err != nil {
		return err
	}
	return nil
}
