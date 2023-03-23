package game

import (
	"math"

	"github.com/mgrabazey/tic-tac-toe/internal/domain"
)

// Strategy is a common interface of move strategy.
type Strategy interface {
	// BestMove provides the best move on domain.GameBoard for domain.GameBoardChar.
	BestMove(board domain.GameBoard, char domain.GameBoardChar) (int, int)
}

type minimaxStrategy struct{}

// NewMinimaxStrategy creates a new MiniMax Strategy. The strategy is aimed at minimizing possible losses.
// The algorithm is described here https://en.wikipedia.org/wiki/Minimax#Example
func NewMinimaxStrategy() Strategy {
	return &minimaxStrategy{}
}

func (s *minimaxStrategy) BestMove(board domain.GameBoard, char domain.GameBoardChar) (int, int) {
	var i, j int

	// Set initial rating value as negative infinity.
	b := math.Inf(-1)

	for i1 := 0; i1 < 3; i1++ {
		for j1 := 0; j1 < 3; j1++ {
			// Consider not used sells only.
			if board[i1][j1] == domain.GameBoardCharNone {
				// Make a move.
				board[i1][j1] = char
				// Get move rating.
				r := s.minimax(board, false, char, math.Inf(-1), math.Inf(1))
				// Rollback board state.
				board[i1][j1] = domain.GameBoardCharNone

				// Update best rating if current rating is better.
				if r > b {
					b = r
					i, j = i1, j1
				}
			}
		}
	}

	return i, j
}

func (s *minimaxStrategy) minimax(board domain.GameBoard, maximizing bool, char domain.GameBoardChar, alpha float64, beta float64) float64 {
	w := board.Winner()
	if w != domain.GameBoardCharNone {
		if w == char {
			return 1 // Won :)
		} else {
			return -1 // Rival won :(
		}
	}

	if board.IsFull() {
		return 0
	}

	if maximizing {
		// Set initial rating value as negative infinity.
		b := math.Inf(-1)

		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				// Consider not used sells only.
				if board[i][j] == domain.GameBoardCharNone {
					// Make a move.
					board[i][j] = char
					// Get move rating.
					r := s.minimax(board, false, char, alpha, beta)
					// Rollback board state.
					board[i][j] = domain.GameBoardCharNone

					// Update best rating if current rating is better.
					b = math.Max(b, r)
					alpha = math.Max(alpha, r)

					if alpha >= beta {
						return b
					}
				}
			}
		}

		return b
	} else {
		// Set initial rating value as positive infinity.
		b := math.Inf(1)

		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				// Consider not used sells only.
				if board[i][j] == domain.GameBoardCharNone {
					c := s.reverseChar(char)
					// Make a move.
					board[i][j] = c
					// Get move rating.
					r := s.minimax(board, true, char, alpha, beta)
					// Rollback board state.
					board[i][j] = domain.GameBoardCharNone

					// Update best rating if current rating is worse.
					b = math.Min(b, r)
					beta = math.Min(beta, r)

					if alpha >= beta {
						return b
					}
				}
			}
		}

		return b
	}
}

func (s *minimaxStrategy) reverseChar(char domain.GameBoardChar) domain.GameBoardChar {
	if char != domain.GameBoardCharCross {
		return domain.GameBoardCharCross
	}
	return domain.GameBoardCharNought
}
