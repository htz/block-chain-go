package blockchain

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type BlockChain struct {
	Chain               []Block       `json:"chain"`
	CurrentTransactions []Transaction `json:"current_transactions"`
	Nodes               []string      `json:"nodes"`
}

func NewBlockChain() *BlockChain {
	blockChain := &BlockChain{}
	// create genesis block
	blockChain.AddNewBlock(0, 100, "1")
	return blockChain
}

func (blockChain *BlockChain) AddNewBlock(timestamp int64, proof int, previousHash string) *Block {
	block := Block{
		Timestamp:    timestamp,
		Transactions: blockChain.CurrentTransactions,
		Proof:        proof,
		PreviousHash: previousHash,
	}
	if block.PreviousHash == "" {
		block.PreviousHash = blockChain.LastBlock().BlockHash()
	}
	blockChain.CurrentTransactions = nil
	blockChain.Chain = append(blockChain.Chain, block)
	return &block
}

func (blockChain *BlockChain) AddNewTransaction(transaction *Transaction) {
	blockChain.CurrentTransactions = append(blockChain.CurrentTransactions, *transaction)
}

func (blockChain *BlockChain) LastBlock() *Block {
	return &blockChain.Chain[len(blockChain.Chain)-1]
}

func (blockChain *BlockChain) ProofOfWork(timestamp int64, lastProof int) int {
	proof := 0
	for !blockChain.validProof(timestamp, lastProof, proof) {
		proof++
	}
	return proof
}

func (blockChain *BlockChain) validProof(timestamp int64, lastProof int, proof int) bool {
	block := &Block{
		Timestamp:    timestamp,
		Transactions: blockChain.CurrentTransactions,
		Proof:        proof,
		PreviousHash: blockChain.LastBlock().PreviousHash,
	}
	hash := block.BlockHash()
	return string(hash[:4]) == "0000"
}

func (blockChain *BlockChain) AddNode(node string) {
	blockChain.Nodes = append(blockChain.Nodes, node)
}

func (blockChain *BlockChain) validChain(chain []Block) bool {
	lastBlock := chain[0]
	for i := 1; i < len(chain); i++ {
		block := chain[i]
		if block.PreviousHash != lastBlock.BlockHash() {
			return false
		}
		if !blockChain.validProof(block.Timestamp, lastBlock.Proof, block.Proof) {
			return false
		}
		lastBlock = block
	}
	return true
}

func (blockChain *BlockChain) ResolveConflicts() bool {
	neighbours := blockChain.Nodes
	maxLength := len(blockChain.Nodes)
	var newChain []Block

	for _, node := range neighbours {
		res, err := http.Get(node + "/chains")
		if res.StatusCode != http.StatusOK {
			fmt.Fprintf(os.Stderr, "http status code: %d\n", res.StatusCode)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			continue
		}
		decoder := json.NewDecoder(res.Body)
		var chain []Block
		if err := decoder.Decode(&chain); err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			continue
		}

		if len(chain) > maxLength && blockChain.validChain(chain) {
			newChain = chain
		}
	}

	if newChain == nil {
		return false
	}

	blockChain.Chain = newChain
	return true
}

func (blockChain *BlockChain) DumpBlockChain() {
	m, err := json.MarshalIndent(blockChain, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(m))
}
