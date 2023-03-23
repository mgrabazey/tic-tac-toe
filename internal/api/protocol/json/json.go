package jsonx

import "github.com/mgrabazey/tic-tac-toe/internal/domain"

type Games []*Game

func NewGames(games domain.Games) Games {
	s := make(Games, len(games))
	for n, i := range games {
		s[n] = NewGame(i)
	}
	return s
}

type Game struct {
	Id     string `json:"id"`
	Board  string `json:"board"`
	Status string `json:"status"`
}

func NewGame(game *domain.Game) *Game {
	return &Game{
		Id:     string(game.Id),
		Board:  game.Board.String(),
		Status: string(game.Status),
	}
}

type GameLocation struct {
	Location string `json:"location"`
}

func NewGameLocation(location string) *GameLocation {
	return &GameLocation{
		Location: location,
	}
}

type Error struct {
	Reason string `json:"reason"`
}

func NewError(reason string) *Error {
	return &Error{
		Reason: reason,
	}
}
