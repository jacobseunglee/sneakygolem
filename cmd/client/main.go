package main

import (
	"sneakygolem/internal/client"
	"flag"
)

func main() {
	flag.Parse()
	client.Run(flag.Arg(0))
}
