package repo

import (
	"context"

	"github.com/mgrabazey/tic-tac-toe/internal/domain"
)

// GameRepository keeps domain.Game entities.
type GameRepository interface {
	// All returns all domain.Game entities.
	All(ctx context.Context) (domain.Games, error)

	// Get returns a domain.Game by the domain.GameId. Returns errorx.NotFound if the
	// domain.Game couldn't be found.
	Get(ctx context.Context, id domain.GameId) (*domain.Game, error)

	// Create creates a new domain.Game.
	Create(ctx context.Context, game *domain.Game) error

	// Update updates the domain.Game. Returns errorx.NotFound error if the domain.Game
	// couldn't be found.
	Update(ctx context.Context, game *domain.Game) error

	// Delete deletes a domain.Game by the domain.GameId. Returns errorx.NotFound if the
	// domain.Game couldn't be found.
	Delete(ctx context.Context, id domain.GameId) error
}
