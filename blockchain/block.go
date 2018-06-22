package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type Block struct {
	Timestamp    int64         `json:"timestamp"`
	Transactions []Transaction `json:"transactions"`
	Proof        int           `json:"proof"`
	PreviousHash string
}

func (block *Block) BlockHash() string {
	marshal, err := json.Marshal(block)
	if err != nil {
		return ""
	}

	bytes := sha256.Sum256([]byte(marshal))
	return hex.EncodeToString(bytes[:])
}
