package game

import (
	"fmt"
	"math/rand"

	"github.com/mgrabazey/tic-tac-toe/internal/domain"
	"github.com/mgrabazey/tic-tac-toe/internal/domain/error"
)

type Engine struct {
	s Strategy
}

func NewEngine(strategy Strategy) *Engine {
	return &Engine{
		s: strategy,
	}
}

func (e *Engine) Start(board domain.GameBoard) (domain.GameBoardChar, domain.GameBoard, error) {
	// Create a new blank board.
	b := domain.NewGameBoard()
	// Calc board difference.
	d := b.Diff(board)

	var c domain.GameBoardChar
	switch len(d) {
	// No move has been made.
	case 0:
		// Select char randomly.
		if rand.Int31n(2) == 0 {
			c = domain.GameBoardCharCross
		} else {
			c = domain.GameBoardCharNought
		}
	// One move has been made.
	case 1:
		// Select char other than user selected.
		if board[d[0][0]][d[0][1]] == domain.GameBoardCharCross {
			c = domain.GameBoardCharNought
		} else {
			c = domain.GameBoardCharCross
		}
	// More than one move was made.
	default:
		return "", domain.GameBoard{}, errorx.WrapInBadRequest(fmt.Errorf("invalid board diff"))
	}

	u, err := e.move(board, c)
	if err != nil {
		return "", domain.GameBoard{}, err
	}
	return c, u, nil
}

func (e *Engine) Move(was, is domain.GameBoard, char domain.GameBoardChar) (domain.GameBoard, error) {
	if was.IsFull() {
		return domain.GameBoard{}, errorx.WrapInBadRequest(fmt.Errorf("board is full"))
	}

	if was.Winner() != domain.GameBoardCharNone {
		return domain.GameBoard{}, errorx.WrapInBadRequest(fmt.Errorf("board has a winner"))
	}

	// Calc board difference.
	d := was.Diff(is)

	switch len(d) {
	// One move has been made.
	case 1:
		c := is[d[0][0]][d[0][1]]
		if c == domain.GameBoardCharNone /*|| c == char*/ {
			return domain.GameBoard{}, errorx.WrapInBadRequest(fmt.Errorf("bad move"))
		}
		// The API doesn't specify what the char (X or O) the user selected.
		// So, let's allow to make a move with any char and simply fix if collides.
		if c == char {
			if c == domain.GameBoardCharCross {
				is[d[0][0]][d[0][1]] = domain.GameBoardCharNought
			} else {
				is[d[0][0]][d[0][1]] = domain.GameBoardCharCross
			}
		}
	// None or more than one move was made.
	default:
		return domain.GameBoard{}, errorx.WrapInBadRequest(fmt.Errorf("invalid board diff"))
	}

	if is.IsFull() {
		return is, nil
	}

	u, err := e.move(is, char)
	if err != nil {
		return domain.GameBoard{}, err
	}
	return u, nil
}

func (e *Engine) move(board domain.GameBoard, char domain.GameBoardChar) (domain.GameBoard, error) {
	// Detect best move.
	i, j := e.s.BestMove(board, char)

	// Ensure that selected call is not used.
	if board[i][j] != domain.GameBoardCharNone {
		return domain.GameBoard{}, fmt.Errorf("used cell chosen")
	}

	// Make move.
	board[i][j] = char

	return board, nil
}
