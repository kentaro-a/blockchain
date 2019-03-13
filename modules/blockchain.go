package modules

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func GenerateGenesisBlock() Block {
	ts := []Transaction{}
	return Block{0, time.Now().String(), ts, "", "", "", ""}
}

func GenerateBlock(oldBlock Block, transactions []Transaction, miner string) Block {
	var newBlock Block
	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = time.Now().String()
	newBlock.Transactions = transactions
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Miner = miner
	Mining(&newBlock)
	return newBlock
}

func Mining(newBlock *Block) {
	for {
		time.Sleep(time.Second)
		newBlock.Nonce = hash(time.Now().String())
		if validateHash(newBlock.toHash()) {
			fmt.Println(newBlock.toHash(), " -> Successfully Found")
			newBlock.Hash = newBlock.toHash()
			break
		} else {
			fmt.Println(newBlock.toHash())
			continue
		}
	}
}

func validateHash(hash string) bool {
	// Difficulty of finding Nonce
	difficulty := 1
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}

func (block Block) toHash() string {
	var buf bytes.Buffer
	b, _ := json.Marshal(block.Transactions)
	buf.Write(b)
	ts := buf.String()
	record := strconv.Itoa(block.Index) + block.Timestamp + ts + block.PrevHash + block.Nonce
	return hash(record)
}

func hash(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func ValidateBlock(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}
	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}
	if newBlock.toHash() != newBlock.Hash {
		return false
	}
	return true
}
