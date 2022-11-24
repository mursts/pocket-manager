package main

import (
	"context"
	"os"

	"github.com/mursts/pocket-manager"
)

func main() {
	if err := pocket.Run(context.Background(), pocket.PubSubMessage{}); err != nil {
		os.Exit(1)
	}
}
