package main

import (
	"context"

	"github.com/mursts/pocket-manager"
)

func main() {
	pocket.Run(context.Background(), pocket.PubSubMessage{})
}
