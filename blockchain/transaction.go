package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type Transaction struct {
	Timestamp int64  `json:"timestamp"`
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Amount    int64  `json:"amount"`
}

func (transaction *Transaction) Hash() string {
	marshal, err := json.Marshal(transaction)
	if err != nil {
		return ""
	}

	bytes := sha256.Sum256([]byte(marshal))
	return hex.EncodeToString(bytes[:])
}

func roundupPowerOf2(n int) int {
	// Bit Twiddling Hacks
	x := uint32(n - 1)
	x |= x >> 1
	x |= x >> 2
	x |= x >> 4
	x |= x >> 8
	x |= x >> 16
	return int(x + 1)
}

func CalcMerkleHash(transactions []Transaction) string {
	len := len(transactions)
	if len == 0 {
		return ""
	}
	size := roundupPowerOf2(len)
	return calcMarkleRoot(transactions, size)
}

func calcMarkleRoot(transactions []Transaction, size int) string {
	if size == 1 {
		return transactions[0].Hash()
	}

	len := len(transactions)
	h1 := calcMarkleRoot(transactions[:size/2], size/2)
	h2 := h1
	if len > size/2 {
		h2 = calcMarkleRoot(transactions[size/2:], size/2)
	}
	bytes := sha256.Sum256([]byte(h1 + h2))

	return hex.EncodeToString(bytes[:])
}
