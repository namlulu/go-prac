package main

import (
	"github.com/namlulu/namlulucoin/explorer"
	"github.com/namlulu/namlulucoin/rest"
)

func main() {
	go explorer.Start(3000)
	rest.Start(4000)
}
