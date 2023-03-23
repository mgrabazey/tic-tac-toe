package httpx

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mgrabazey/tic-tac-toe/internal/api/protocol/json"
	"github.com/mgrabazey/tic-tac-toe/internal/app/module/game"
	"github.com/mgrabazey/tic-tac-toe/internal/domain"
	"github.com/mgrabazey/tic-tac-toe/internal/domain/error"
)

type gameController struct {
	u string
	s *game.Service
}

func newGameController(publicUrl string, service *game.Service) *gameController {
	return &gameController{
		u: publicUrl,
		s: service,
	}
}

func (c *gameController) all(writer http.ResponseWriter, request *http.Request) {
	v, err := c.s.All(request.Context())
	if err != nil {
		writeError(writer, err, nil)
	}
	writeResponse(writer, http.StatusOK, jsonx.NewGames(v))
}

func (c *gameController) get(writer http.ResponseWriter, request *http.Request) {
	id, ok := c.validateId(writer, request)
	if !ok {
		return
	}
	v, err := c.s.Get(request.Context(), id)
	if err != nil {
		writeError(writer, err, nil)
		return
	}
	writeResponse(writer, http.StatusOK, jsonx.NewGame(v))
}

func (c *gameController) create(writer http.ResponseWriter, request *http.Request) {
	b, ok := c.validateBody(writer, request)
	if !ok {
		return
	}
	v, err := c.s.Create(request.Context(), &game.CreateRequest{
		Board: b,
	})
	if err != nil {
		writeError(writer, err, jsonx.NewError(err.Error()))
		return
	}
	u := fmt.Sprintf("%s/api/v1/games/%s", c.u, v.Id)
	writer.Header().Set("Location", u)
	writeResponse(writer, http.StatusCreated, jsonx.NewGameLocation(u))
}

func (c *gameController) update(writer http.ResponseWriter, request *http.Request) {
	id, ok := c.validateId(writer, request)
	if !ok {
		return
	}
	b, ok := c.validateBody(writer, request)
	if !ok {
		return
	}
	v, err := c.s.Update(request.Context(), &game.UpdateRequest{
		Id:    id,
		Board: b,
	})
	if err != nil {
		writeError(writer, err, jsonx.NewError(err.Error()))
		return
	}
	writeResponse(writer, http.StatusOK, jsonx.NewGame(v))
}

func (c *gameController) remove(writer http.ResponseWriter, request *http.Request) {
	id, ok := c.validateId(writer, request)
	if !ok {
		return
	}
	err := c.s.Delete(request.Context(), id)
	if err != nil {
		writeError(writer, err, nil)
		return
	}
	writeResponse(writer, http.StatusOK, nil)
}

func (c *gameController) validateId(writer http.ResponseWriter, request *http.Request) (domain.GameId, bool) {
	id, err := domain.GameIdFromString(mux.Vars(request)["id"])
	if err != nil {
		writeError(writer, errorx.WrapInBadRequest(err), jsonx.NewError("Invalid id"))
		return "", false
	}
	return id, true
}

func (c *gameController) validateBody(writer http.ResponseWriter, request *http.Request) (domain.GameBoard, bool) {
	g := &jsonx.Game{}
	err := json.NewDecoder(request.Body).Decode(g)
	if err != nil {
		writeError(writer, errorx.WrapInBadRequest(err), jsonx.NewError("Invalid request body"))
		return domain.GameBoard{}, false
	}
	if g.Id != "" {
		writeError(writer, errorx.NewBadRequest(), jsonx.NewError("Unexpected parameter: id"))
		return domain.GameBoard{}, false
	}
	if g.Status != "" {
		writeError(writer, errorx.NewBadRequest(), jsonx.NewError("Unexpected parameter: status"))
		return domain.GameBoard{}, false
	}
	b, err := domain.GameBoardFromString(g.Board)
	if err != nil {
		writeError(writer, errorx.WrapInBadRequest(err), jsonx.NewError("Invalid board"))
		return domain.GameBoard{}, false
	}
	return b, true
}
