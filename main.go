package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	. "./entity"
	. "./util"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func decodeBody(r *http.Request, out interface{}) error {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(out)
	if err != nil {
		log.Printf("Error occured when decoding RequestBody")
	}

	return err
}

func getBoardHandler(w http.ResponseWriter, r *http.Request) {
	board := Board{}
	board.GetBoardFromDb()

	fmt.Fprintf(w, board.ToJsonString())
}

func postBoardHandler(w http.ResponseWriter, r *http.Request) {
	var board Board

	if err := decodeBody(r, &board); err != nil {
		http.Error(w, http.StatusText(422), 422)
		log.Panic(err)
	}

	if err := board.WriteBoardToDb(); err != nil {
		http.Error(w, http.StatusText(422), 422)
		log.Panic(err)
	}

	fmt.Fprintf(w, "success")
}

func originTokenHandler(w http.ResponseWriter, r *http.Request) {
	originToken, err := GenerateOriginToken(OriginTokenLength)
	if err != nil {
		log.Printf("Error occuerred during generating origin-token.")
		log.Panic(err)
	}

	log.Printf("Generated Token: " + originToken)

	fmt.Fprintf(w, originToken)
}

func staticHundler(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/", http.FileServer(http.Dir(HtmlResourceDir))).ServeHTTP(w, r)
}

func main() {

	InitDB()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/api", func(r chi.Router) {
		// TODO token should be handed with header(?)
		r.Get("/origin-token", originTokenHandler)

		r.Get("/board", getBoardHandler)
		r.Post("/board", postBoardHandler)
	})

	r.Get("/*", staticHundler)

	log.Printf("Started to listen: " + ListenPort)
	http.ListenAndServe(":"+ListenPort, r)
}
