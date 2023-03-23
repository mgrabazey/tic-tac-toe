package repo

import (
	"context"
	"database/sql"
	"time"

	"github.com/mgrabazey/tic-tac-toe/internal/domain"
	"github.com/mgrabazey/tic-tac-toe/internal/domain/error"
	"github.com/mgrabazey/tic-tac-toe/internal/domain/repo"
)

type game struct {
	id     string
	board  string
	status string
	char   string
}

func (g *game) scan(scanner func(...any) error) error {
	return scanner(
		&g.id,
		&g.board,
		&g.status,
		&g.char,
	)
}

func (g *game) to() *domain.Game {
	return &domain.Game{
		Id:     domain.MustGameIdFromString(g.id),
		Board:  domain.MustGameBoardFromString(g.board),
		Status: domain.GameStatus(g.status),
		Char:   domain.GameBoardChar(g.char),
	}
}

type gameRepository struct {
	db *sql.DB
}

func NewGameRepository(db *sql.DB) repo.GameRepository {
	return &gameRepository{
		db: db,
	}
}

func (r *gameRepository) All(ctx context.Context) (domain.Games, error) {
	q := `SELECT "id", "board", "status", "char" FROM "games" WHERE "deleted_at" IS NULL ORDER BY "created_at" ASC`
	v, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer v.Close()

	var s domain.Games
	for v.Next() {
		i := &game{}
		err = i.scan(v.Scan)
		if err != nil {
			return nil, err
		}
		s = append(s, i.to())
	}
	return s, nil
}

func (r *gameRepository) Get(ctx context.Context, id domain.GameId) (*domain.Game, error) {
	q := `SELECT "id", "board", "status", "char" FROM "games" WHERE "id" = $1 AND "deleted_at" IS NULL`
	v := r.db.QueryRowContext(ctx, q, id)
	i := &game{}
	err := i.scan(v.Scan)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errorx.NewNotFound()
		}
		return nil, err
	}
	return i.to(), nil
}

func (r *gameRepository) Create(ctx context.Context, game *domain.Game) error {
	q := `INSERT INTO "games" ("id", "board", "status", "char", "created_at", "updated_at") VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, q, game.Id, game.Board.String(), game.Status, game.Char, time.Now(), time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (r *gameRepository) Update(ctx context.Context, game *domain.Game) error {
	q := `UPDATE "games" SET "board" = $1, "status" = $2, "char" = $3, "updated_at" = $4  WHERE "id" = $5 AND "deleted_at" IS NULL`
	v, err := r.db.ExecContext(ctx, q, game.Board.String(), game.Status, game.Char, time.Now(), game.Id)
	if err != nil {
		return err
	}
	n, err := v.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return errorx.NewNotFound()
	}
	return nil
}

func (r *gameRepository) Delete(ctx context.Context, id domain.GameId) error {
	q := `UPDATE "games" SET "deleted_at" = $1 WHERE "id" = $2`
	v, err := r.db.ExecContext(ctx, q, time.Now(), id)
	if err != nil {
		return err
	}
	n, err := v.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return errorx.NewNotFound()
	}
	return nil
}
