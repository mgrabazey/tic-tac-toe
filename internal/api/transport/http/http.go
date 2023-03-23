package httpx

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/mgrabazey/tic-tac-toe/internal/app/module/game"
	"github.com/mgrabazey/tic-tac-toe/internal/domain/error"
)

func Run(publicUrl string, gameService *game.Service) error {
	r := mux.NewRouter()

	g := newGameController(publicUrl, gameService)

	r.Methods(http.MethodGet).Path("/api/v1/games").HandlerFunc(g.all)
	r.Methods(http.MethodGet).Path("/api/v1/games/{id}").HandlerFunc(g.get)
	r.Methods(http.MethodPost).Path("/api/v1/games").HandlerFunc(g.create)
	r.Methods(http.MethodPut).Path("/api/v1/games/{id}").HandlerFunc(g.update)
	r.Methods(http.MethodDelete).Path("/api/v1/games/{id}").HandlerFunc(g.remove)

	s := &http.Server{
		Handler: handlers.CORS(
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"}),
			handlers.AllowedHeaders([]string{"Content-Type"}),
		)(r),
		Addr: ":80",
	}
	log.Printf("Run HTTP server...\n")
	return s.ListenAndServe()
}

func writeResponse(writer http.ResponseWriter, code int, data any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	if data == nil {
		return
	}
	b, err := json.Marshal(data)
	if err != nil {
		log.Printf("unable to encode response data: %v\n", err)
		return
	}
	_, err = writer.Write(b)
	if err != nil {
		log.Printf("unable to write response data: %v\n", err)
		return
	}
}

func writeError(writer http.ResponseWriter, err error, data any) {
	var code int
	switch true {
	case errorx.IsBadRequest(err):
		code = http.StatusBadRequest
	case errorx.IsNotFound(err):
		code = http.StatusNotFound
	default:
		code = http.StatusInternalServerError
	}
	writeResponse(writer, code, data)
}
