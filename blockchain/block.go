package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type Block struct {
	Height       int           `json:"height"`
	Timestamp    int64         `json:"timestamp"`
	Nonce        int           `json:"nonce"`
	Hash         string        `json:"hash"`
	PreviousHash string        `json:"previous_hash"`
	MerkleHash   string        `json:"merkle_hash"`
	Transactions []Transaction `json:"transactions"`
}

const BlockDifficulty = 5

func NewBlock(timestamp int64, nonce int, previousHash string, merkleHash string, transactions []Transaction) *Block {
	block := &Block{
		Timestamp:    timestamp,
		Nonce:        nonce,
		PreviousHash: previousHash,
		MerkleHash:   merkleHash,
		Transactions: transactions,
	}
	if !block.IsValid() {
		return nil
	}

	return block
}

func (block *Block) hash() string {
	hashSeed := struct {
		Timestamp    int64  `json:"timestamp"`
		Nonce        int    `json:"nonce"`
		PreviousHash string `json:"previous_hash"`
		MerkleHash   string `json:"merkle_hash"`
	}{
		Timestamp:    block.Timestamp,
		Nonce:        block.Nonce,
		PreviousHash: block.PreviousHash,
		MerkleHash:   block.MerkleHash,
	}

	marshal, err := json.Marshal(hashSeed)
	if err != nil {
		return ""
	}
	bytes := sha256.Sum256([]byte(marshal))

	return hex.EncodeToString(bytes[:])
}

func (block *Block) IsValid() bool {
	hash := block.hash()
	for i := 0; i < BlockDifficulty; i++ {
		if hash[i:i+1] != "0" {
			return false
		}
	}
	block.Hash = hash
	return true
}
