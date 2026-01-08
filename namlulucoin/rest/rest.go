package rest

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/namlulu/namlulucoin/blockchain"
)

const port string = ":4000"

var templates *template.Template

type addBlockBody struct {
	Data string `json:"data"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type homeData struct {
	BlockCount int
	Blocks     []*blockchain.Block
}

func getAllBlocks(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blockchain.GetBlockchain().AllBlocks())
}

func addBlock(w http.ResponseWriter, r *http.Request) {
	var body addBlockBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "Invalid request body"})
		return
	}
	if body.Data == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "Data is required"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	newBlock := blockchain.GetBlockchain().AddBlock(body.Data)
	json.NewEncoder(w).Encode(newBlock)
}

func home(w http.ResponseWriter, r *http.Request) {
	blocks := blockchain.GetBlockchain().AllBlocks()
	data := homeData{
		BlockCount: len(blocks),
		Blocks:     blocks,
	}
	templates.ExecuteTemplate(w, "home.html", data)
}

func Start() {
	templates = template.Must(template.ParseGlob("templates/*.html"))

	http.HandleFunc("/", home)
	http.HandleFunc("/blocks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getAllBlocks(w, r)
		case http.MethodPost:
			addBlock(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	fmt.Printf("🚀 Server running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
