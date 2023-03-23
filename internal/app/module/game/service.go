package game

import (
	"context"
	"fmt"
	"log"

	"github.com/mgrabazey/tic-tac-toe/internal/domain"
	"github.com/mgrabazey/tic-tac-toe/internal/domain/error"
	"github.com/mgrabazey/tic-tac-toe/internal/domain/repo"
)

type CreateRequest struct {
	Board domain.GameBoard
}

type UpdateRequest struct {
	Id    domain.GameId
	Board domain.GameBoard
}

type Service struct {
	r repo.GameRepository
	e *Engine
}

func NewService(repo repo.GameRepository) *Service {
	return &Service{
		r: repo,
		// Not the best solution, but good enough since we only have one strategy.
		e: NewEngine(NewMinimaxStrategy()),
	}
}

func (s *Service) All(ctx context.Context) (domain.Games, error) {
	// Get all games.
	v, err := s.r.All(ctx)
	if err != nil {
		log.Printf("Unable to get games: %v\n", err)
		return nil, err
	}
	return v, nil
}

func (s *Service) Get(ctx context.Context, id domain.GameId) (*domain.Game, error) {
	// Get game by identifier.
	v, err := s.r.Get(ctx, id)
	if err != nil {
		log.Printf("Unable to get game: %v\n", err)
		return nil, err
	}
	return v, nil
}

func (s *Service) Create(ctx context.Context, request *CreateRequest) (*domain.Game, error) {
	// Create new game.
	g := &domain.Game{
		Id:     domain.NewGameId(),
		Status: domain.GameStatusRunning,
	}
	// Check user's move if any and make own.
	var err error
	g.Char, g.Board, err = s.e.Start(request.Board)
	if err != nil {
		log.Printf("Unable to start game: %v\n", err)
		return nil, err
	}
	// Create game
	err = s.r.Create(ctx, g)
	if err != nil {
		log.Printf("Unable to create game: %v\n", err)
		return nil, err
	}
	return g, nil
}

func (s *Service) Update(ctx context.Context, request *UpdateRequest) (*domain.Game, error) {
	// Get game by identifier.
	g, err := s.r.Get(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	// Throw if game is already over.
	if g.Status != domain.GameStatusRunning {
		return nil, errorx.WrapInBadRequest(fmt.Errorf("game is already over"))
	}
	// Check user's move and make own.
	g.Board, err = s.e.Move(g.Board, request.Board, g.Char)
	if err != nil {
		log.Printf("Unable to move: %v\n", err)
		return nil, err
	}
	// Check if game is already over.
	switch g.Board.Winner() {
	case domain.GameBoardCharCross:
		g.Status = domain.GameStatusCrossWon
	case domain.GameBoardCharNought:
		g.Status = domain.GameStatusNoughtWon
	default:
		if g.Board.IsFull() {
			g.Status = domain.GameStatusDraw
		}
	}
	// Update game,
	err = s.r.Update(ctx, g)
	if err != nil {
		log.Printf("Unable to update game: %v\n", err)
		return nil, err
	}
	return g, nil
}

func (s *Service) Delete(ctx context.Context, id domain.GameId) error {
	// Delete game by identifier.
	err := s.r.Delete(ctx, id)
	if err != nil {
		log.Printf("Unable to delete game: %v\n", err)
		return err
	}
	return nil
}
