package domain

import (
	"fmt"

	"github.com/google/uuid"
)

// Games is a list of Game.
type Games []*Game

// Game represents a game.
type Game struct {
	Id     GameId
	Board  GameBoard
	Status GameStatus
	Char   GameBoardChar
}

// GameId represents Game identifier.
type GameId string

// NewGameId generates a new GameId.
func NewGameId() GameId {
	return GameId(uuid.NewString())
}

// GameIdFromString creates a new GameId from string. It returns error
// if string is not a valid UUID.
func GameIdFromString(s string) (GameId, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return "", fmt.Errorf("invalid domain.GameId: %v", err)
	}
	return GameId(id.String()), nil
}

// MustGameIdFromString wraps GameIdFromString. It panics if
// GameIdFromString returns an error.
func MustGameIdFromString(s string) GameId {
	id, err := GameIdFromString(s)
	if err != nil {
		panic(fmt.Sprintf("MustGameIdFromString: %v", err))
	}
	return id
}

// GameBoard represents Game board.
type GameBoard [3][3]GameBoardChar

// NewGameBoard creates a new blank GameBoard.
func NewGameBoard() GameBoard {
	return GameBoard{
		{GameBoardCharNone, GameBoardCharNone, GameBoardCharNone},
		{GameBoardCharNone, GameBoardCharNone, GameBoardCharNone},
		{GameBoardCharNone, GameBoardCharNone, GameBoardCharNone},
	}
}

// GameBoardFromString creates a new GameBoard from string. It returns an error if
// the string doesn't contain 9 chars or contains chars other than GameBoardCharCross,
// GameBoardCharNought and GameBoardCharNone.
func GameBoardFromString(s string) (GameBoard, error) {
	var b GameBoard
	if len(s) != 9 {
		return b, fmt.Errorf("invalid domain.GameBoard")
	}
	n := 0
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			c := GameBoardChar(s[n])
			n++
			switch c {
			case GameBoardCharNone, GameBoardCharCross, GameBoardCharNought:
				// OK
			default:
				return b, fmt.Errorf("invalid domain.GameBoard")
			}
			b[i][j] = c
		}
	}
	return b, nil
}

// MustGameBoardFromString wraps GameBoardFromString. It panics if
// GameBoardFromString returns an error.
func MustGameBoardFromString(s string) GameBoard {
	b, err := GameBoardFromString(s)
	if err != nil {
		panic(fmt.Sprintf("MustGameBoardFromString: %v", err))
	}
	return b
}

func (a *GameBoard) String() string {
	var s string
	for _, i := range a {
		for _, j := range i {
			s += string(j)
		}
	}
	return s
}

// IsFull check if the GameBoard is full.
func (a *GameBoard) IsFull() bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if a[i][j] == GameBoardCharNone {
				return false
			}
		}
	}
	return true
}

// Diff provides different positions.
func (a *GameBoard) Diff(board GameBoard) [][2]int {
	var s [][2]int
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if a[i][j] != board[i][j] {
				s = append(s, [2]int{i, j})
			}
		}
	}
	return s
}

// Winner checks if there is a winner and returns winner's GameBoardChar.
// It returns GameBoardCharNone if there is no winner yet.
func (a *GameBoard) Winner() GameBoardChar {
	// Check rows
	for i := 0; i < 3; i++ {
		c := a[i][0]
		if c != GameBoardCharNone && c == a[i][1] && c == a[i][2] {
			return c
		}
	}

	// Check cols
	for i := 0; i < 3; i++ {
		c := a[0][i]
		if c != GameBoardCharNone && c == a[1][i] && c == a[2][i] {
			return c
		}
	}

	// Check left-right diagonal
	c00 := a[0][0]
	if c00 != GameBoardCharNone && c00 == a[1][1] && c00 == a[2][2] {
		return c00
	}

	// Check right-left diagonal
	c02 := a[0][2]
	if c02 != GameBoardCharNone && c02 == a[1][1] && c02 == a[2][0] {
		return c02
	}

	return GameBoardCharNone
}

type GameBoardChar string

const (
	GameBoardCharNone   GameBoardChar = "-"
	GameBoardCharCross  GameBoardChar = "X"
	GameBoardCharNought GameBoardChar = "0"
)

// GameStatus represents Game status.
type GameStatus string

const (
	// GameStatusRunning means that the Game is in progress.
	GameStatusRunning GameStatus = "RUNNING"
	// GameStatusCrossWon means that the Game is over and GameBoardCharCross is winner.
	GameStatusCrossWon GameStatus = "X_WON"
	// GameStatusNoughtWon means that the Game is over and GameBoardCharNought is winner.
	GameStatusNoughtWon GameStatus = "O_WON"
	// GameStatusDraw means that the Game is over because the GameBoard is full.
	GameStatusDraw GameStatus = "DRAW"
)
