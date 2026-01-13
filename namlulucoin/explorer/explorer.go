package explorer

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/namlulu/namlulucoin/blockchain"
)

const (
	port            = ":4000"
	contentTypeJSON = "application/json"
	templatesDir    = "explorer/templates/*.html"
	partialsDir     = "explorer/templates/partials/*.html"
)

var templates *template.Template

type addBlockBody struct {
	Data string `json:"data"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type pageData struct {
	PageTitle  string
	BlockCount int
	Blocks     []*blockchain.Block
}

func jsonResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", contentTypeJSON)
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func getAllBlocks(w http.ResponseWriter, _ *http.Request) {
	jsonResponse(w, http.StatusOK, blockchain.GetBlockchain().AllBlocks())
}

func addBlock(w http.ResponseWriter, r *http.Request) {
	var body addBlockBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonResponse(w, http.StatusBadRequest, errorResponse{Error: "Invalid request body"})
		return
	}
	if body.Data == "" {
		jsonResponse(w, http.StatusBadRequest, errorResponse{Error: "Data is required"})
		return
	}
	newBlock := blockchain.GetBlockchain().AddBlock(body.Data)
	jsonResponse(w, http.StatusCreated, newBlock)
}

func home(w http.ResponseWriter, _ *http.Request) {
	blocks := blockchain.GetBlockchain().AllBlocks()
	data := pageData{
		PageTitle:  "Home",
		BlockCount: len(blocks),
		Blocks:     blocks,
	}
	templates.ExecuteTemplate(w, "home.html", data)
}

func blocksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getAllBlocks(w, r)
	case http.MethodPost:
		addBlock(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func addBlockForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	data := r.FormValue("data")
	if data != "" {
		blockchain.GetBlockchain().AddBlock(data)
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func Start() {
	templates = template.Must(template.ParseGlob(templatesDir))
	templates = template.Must(templates.ParseGlob(partialsDir))

	http.HandleFunc("/", home)
	http.HandleFunc("/add", addBlockForm)
	http.HandleFunc("/blocks", blocksHandler)

	fmt.Printf("🚀 Server running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
