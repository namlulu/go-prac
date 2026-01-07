package main

import (
	"github.com/namlulu/namlulucoin/blockchain"
)

func main() {
	chain := blockchain.GetBlockchain()
	chain.AddBlock("First Block")
	chain.AddBlock("Second Block")
	chain.AddBlock("Third Block")
	chain.ListBlocks()
}
