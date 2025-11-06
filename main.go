package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"	
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

const difficulty = 4

// Block represents each 'item' in the blockchain
type Block struct {
	Index     int
	Timestamp string
	BPM       int
	Hash      string
	PrevHash  string
	Nonce     int
}

// Blockchain is a series of validated Blocks
var Blockchain []Block

// calculateHash returns the SHA256 hash of a Block
func calculateHash(block Block) string {
	record := strconv.Itoa(block.Index) + block.Timestamp + strconv.Itoa(block.BPM) + block.PrevHash + strconv.Itoa(block.Nonce)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// generateBlock creates a new block using previous block's hash
func generateBlock(oldBlock Block, BPM int) Block {
	var newBlock Block
	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = oldBlock.Hash

	for i := 0; ; i++ {
		newBlock.Nonce = i
		newBlock.Hash = calculateHash(newBlock)
		if isHashValid(newBlock.Hash, difficulty) {
			fmt.Printf("Block mined! Nonce: %d, Hash: %s\n", i, newBlock.Hash)
			break
		}
	}
	return newBlock
}

// isHashValid checks if hash satisfies the difficulty level
func isHashValid(hash string, difficulty int) bool {
	prefix := ""
	for i := 0; i < difficulty; i++ {
		prefix += "0"
	}
	return hash[:difficulty] == prefix
}

// isBlockValid verifies new block's integrity
func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

// handleGetBlockchain returns the blockchain as JSON
func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(Blockchain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

// handleMineBlock mines a new block with BPM from URL
func handleMineBlock(w http.ResponseWriter, r *http.Request) {
	bpmStr := r.URL.Query().Get("bpm")
	bpm, err := strconv.Atoi(bpmStr)
	if err != nil {
		http.Error(w, "Invalid BPM value", http.StatusBadRequest)
		return
	}

	oldBlock := Blockchain[len(Blockchain)-1]
	newBlock := generateBlock(oldBlock, bpm)

	if isBlockValid(newBlock, oldBlock) {
		Blockchain = append(Blockchain, newBlock)
		fmt.Println("New block added to blockchain!")
	}

	bytes, err := json.MarshalIndent(Blockchain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

// main initializes the blockchain and web server
func main() {
	fmt.Println("Starting blockchain web server with Proof of Work (Go)")

	t := time.Now()
	genesisBlock := Block{0, t.String(), 0, calculateHash(Block{}), "", 0}
	Blockchain = append(Blockchain, genesisBlock)

	r := mux.NewRouter()
	r.HandleFunc("/", handleGetBlockchain).Methods("GET")
	r.HandleFunc("/mine", handleMineBlock).Methods("GET")

	fmt.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

