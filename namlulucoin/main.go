package main

import (
	"crypto/sha256"
	"fmt"
)

type block struct {
	data     string
	hash     string
	prevHash string
}

type blockchain struct {
	blocks []block
}

func (b *blockchain) getLastHash() string {
	if len(b.blocks) == 0 {
		return ""
	}
	return b.blocks[len(b.blocks)-1].hash
}

func (b *blockchain) listBlocks() []block {
	for _, block := range b.blocks {
		fmt.Printf("Data: %s\n", block.data)
		fmt.Printf("Hash: %s\n", block.hash)
		fmt.Printf("Previous Hash: %s\n", block.prevHash)
		fmt.Println()
	}
	return b.blocks
}

func (b *blockchain) addBlock(data string) {
	newBlock := block{data, "", b.getLastHash()}
	hash := sha256.Sum256([]byte(newBlock.data + newBlock.prevHash))
	hexHash := fmt.Sprintf("%x", hash)
	newBlock.hash = hexHash
	b.blocks = append(b.blocks, newBlock)
}

func main() {
	chain := blockchain{}
	chain.addBlock("First Block")
	chain.addBlock("Second Block")
	chain.addBlock("Third Block")

	chain.listBlocks()
}
