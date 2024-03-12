package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type BlockData struct {
	Message string `json:"message"`
}

type Block struct {
	Data         BlockData `json:"data"`
	Hash         string    `json:"hash"`
	PreviousHash string    `json:"previous_hash"`
	Timestamp    int64     `json:"timestamp"`
	Height       int       `json:"height"`
	Pow          int       `json:"pow"`
}

func (b *Block) calculateHash() string {
	byt, _ := json.Marshal(b.Data)
	return fmt.Sprintf("%x", sha256.Sum256([]byte(b.PreviousHash+string(byt)+strconv.FormatInt(b.Timestamp, 10)+strconv.Itoa(b.Pow))))
}

func (b *Block) mine(difficulty int) {
	for !strings.HasPrefix(b.Hash, strings.Repeat("0", difficulty)) {
		b.Pow++
		b.Hash = b.calculateHash()
	}
}
