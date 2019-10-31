package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
	"blockchain"

	"github.com/gorilla/mux"
)

var blockChain = blockchain.NewBlockChain()
//var nodeIdentifire = uuid.Must(uuid.NewV4()).String()

func createTransactionHandler(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var transaction blockchain.Transaction
	if err := decoder.Decode(&transaction); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	blockChain.AddTransaction(&transaction)
	w.WriteHeader(http.StatusCreated)
	blockChain.PrintDump()
}

func getMineHandler(w http.ResponseWriter, req *http.Request) {
	timestamp := time.Now().Unix()
	nonce := blockChain.ProofOfWork(timestamp)
	block := blockChain.AddBlock(timestamp, nonce)
	if block == nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(block); err != nil {
		log.Println("Error:", err)
	}
	blockChain.PrintDump()
}

func getChainsHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(blockChain.Chain); err != nil {
		log.Println("Error:", err)
	}
	blockChain.PrintDump()
}

func registerNodesHandler(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var nodes []string
	if err := decoder.Decode(&nodes); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for _, node := range nodes {
		blockChain.AddNode(node)
	}
	w.WriteHeader(http.StatusCreated)
	blockChain.PrintDump()
}

func consensusNodesHandler(w http.ResponseWriter, req *http.Request) {
	replaced := blockChain.ResolveConflicts()
	if replaced {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
	blockChain.PrintDump()
}

func listenAddress() string {
	port := os.Getenv("PORT")
	if port == "" {
		return "8080"
	}
	return ":" + port
}

func init() {
	router := mux.NewRouter()
	router.HandleFunc("/transactions", createTransactionHandler).Methods("POST")
	router.HandleFunc("/mine", getMineHandler).Methods("POST")
	router.HandleFunc("/chains", getChainsHandler).Methods("GET")
	router.HandleFunc("/nodes", registerNodesHandler).Methods("POST")
	router.HandleFunc("/nodes/resolve", consensusNodesHandler).Methods("GET")
	http.Handle("/", router)
}
