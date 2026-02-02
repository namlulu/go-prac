package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/namlulu/namlulucoin/blockchain"
	"github.com/namlulu/namlulucoin/utils"
)

type url string

var port string

func (u url) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost:4000%s%s", port, u)
	return []byte(url), nil
}

type urlDescription struct {
	URL         url    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

type addBlockBody struct {
	Message string
}

func (u urlDescription) String() string {
	return "URL: " + string(u.URL) + "\nMethod: " + u.Method + "\nDescription: " + u.Description + "\nPayload:"
}

func documentation(w http.ResponseWriter, r *http.Request) {
	data := []urlDescription{{
		URL:         url("/"),
		Method:      "GET",
		Description: "Hello, World!",
	}, {
		URL:         url("/blocks"),
		Method:      "POST",
		Description: "Add a new block",
		Payload:     "data: string",
	}, {
		URL:         url("/blocks/{id}"),
		Method:      "GET",
		Description: "See a block",
	}}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func blocks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(blockchain.GetBlockchain().AllBlocks())
	case "POST":
		var addBlockBody addBlockBody
		utils.HandleErr(json.NewDecoder(r.Body).Decode(&addBlockBody))
		blockchain.GetBlockchain().AddBlock(addBlockBody.Message)
		w.WriteHeader(http.StatusCreated)
	}
}

func Start(aPort int) {
	handler := http.NewServeMux()
	port = fmt.Sprintf(":%d", aPort)
	http.HandleFunc("/", documentation)
	http.HandleFunc("/blocks", blocks)
	fmt.Printf("Server is running on port %d", aPort)
	http.ListenAndServe(port, handler)
}
