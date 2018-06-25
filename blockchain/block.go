package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type Block struct {
	Timestamp    int64         `json:"timestamp"`
	Proof        int           `json:"proof"`
	Hash         string        `json:"hash"`
	PreviousHash string        `json:"previous_hash"`
	Transactions []Transaction `json:"transactions"`
}

const BlockDifficulty = 5

func NewBlock(timestamp int64, proof int, previousHash string, transactions []Transaction) *Block {
	block := &Block{
		Timestamp:    timestamp,
		Proof:        proof,
		PreviousHash: previousHash,
		Transactions: transactions,
	}
	if !block.BlockValidProof() {
		return nil
	}

	return block
}

func (block *Block) blockHash() string {
	hashSeed := struct {
		Timestamp    int64         `json:"timestamp"`
		Proof        int           `json:"proof"`
		PreviousHash string        `json:"previous_hash"`
		Transactions []Transaction `json:"transactions"`
	}{
		Timestamp:    block.Timestamp,
		Proof:        block.Proof,
		PreviousHash: block.PreviousHash,
		Transactions: block.Transactions,
	}

	marshal, err := json.Marshal(hashSeed)
	if err != nil {
		return ""
	}

	bytes := sha256.Sum256([]byte(marshal))
	block.Hash = hex.EncodeToString(bytes[:])

	return block.Hash
}

func (block *Block) BlockValidProof() bool {
	hash := block.blockHash()
	for i := 0; i < BlockDifficulty; i++ {
		if hash[i:i+1] != "0" {
			return false
		}
	}
	return true
}
