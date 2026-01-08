package blockchain

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

type Block struct {
	Data     string `json:"data"`
	Hash     string `json:"hash"`
	PrevHash string `json:"prevHash"`
}

type blockchain struct {
	blocks []*Block
}

var b *blockchain
var once sync.Once

func (b *Block) calculateHash() {
	hash := sha256.Sum256([]byte(b.Data + b.PrevHash))
	b.Hash = fmt.Sprintf("%x", hash)
}

func (b *blockchain) AddBlock(data string) *Block {
	newBlock := createBlock(data)
	b.blocks = append(b.blocks, newBlock)
	return newBlock
}

func createBlock(data string) *Block {
	newBlock := Block{Data: data, PrevHash: getLastHash()}
	newBlock.calculateHash()
	return &newBlock
}

func (b *blockchain) AllBlocks() []*Block {
	return b.blocks
}

func getLastHash() string {
	totalBlocks := len(GetBlockchain().blocks)
	if totalBlocks == 0 {
		return ""
	}
	return GetBlockchain().blocks[totalBlocks-1].Hash
}

func GetBlockchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{}
			b.AddBlock("Genesis Block")
		})
	}
	return b
}
