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

const GenesisTimestamp = int64(0)
const GenesisPreviousHash = "0000000000000000000000000000000000000000000000000000000000000000"

func NewBlockChain() *BlockChain {
	blockChain := &BlockChain{}
	nonce := blockChain.ProofOfWork(GenesisTimestamp)
	blockChain.AddBlock(GenesisTimestamp, nonce)
	return blockChain
}

func (blockChain *BlockChain) AddBlock(timestamp int64, nonce int) *Block {
	block := NewBlock(
		timestamp,
		nonce,
		blockChain.previousHash(),
		blockChain.CurrentTransactions,
	)
	if block == nil {
		return nil
	}

	blockChain.CurrentTransactions = nil
	blockChain.appendBlock(block)

	return block
}

func (blockChain *BlockChain) appendBlock(block *Block) {
	blockChain.Chain = append(blockChain.Chain, *block)
}

func (blockChain *BlockChain) AddTransaction(transaction *Transaction) {
	blockChain.CurrentTransactions = append(blockChain.CurrentTransactions, *transaction)
}

func (blockChain *BlockChain) lastBlock() *Block {
	if blockChain.Chain == nil {
		return nil
	}
	return &blockChain.Chain[len(blockChain.Chain)-1]
}

func (blockChain *BlockChain) previousHash() string {
	lastBlock := blockChain.lastBlock()
	if lastBlock != nil {
		return lastBlock.Hash
	}
	return GenesisPreviousHash
}

func (blockChain *BlockChain) ProofOfWork(timestamp int64) int {
	nonce := 0
	for !blockChain.validNonce(timestamp, nonce) {
		nonce++
	}
	return nonce
}

func (blockChain *BlockChain) validNonce(timestamp int64, nonce int) bool {
	block := NewBlock(
		timestamp,
		nonce,
		blockChain.previousHash(),
		blockChain.CurrentTransactions,
	)
	return block != nil
}

func (blockChain *BlockChain) AddNode(node string) {
	blockChain.Nodes = append(blockChain.Nodes, node)
}

func (blockChain *BlockChain) isValidChain(chain []Block) bool {
	lastBlock := chain[0]
	for i := 1; i < len(chain); i++ {
		block := chain[i]
		if block.PreviousHash != lastBlock.Hash || !block.IsValid() {
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

		if len(chain) > maxLength && blockChain.isValidChain(chain) {
			newChain = chain
		}
	}

	if newChain == nil {
		return false
	}

	blockChain.Chain = newChain
	return true
}

func (blockChain *BlockChain) PrintDump() {
	m, err := json.MarshalIndent(blockChain, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(m))
}
