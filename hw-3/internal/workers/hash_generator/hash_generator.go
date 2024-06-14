package hashgen

import (
	"context"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-3/pkg/hash"
)

type HashGenerator struct {
	hashes chan string
}

func NewHashGenerator(hashesN int) *HashGenerator {
	return &HashGenerator{make(chan string, hashesN)}
}

func (g *HashGenerator) Run(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(g.hashes)

				return
			default:
				g.hashes <- hash.GenerateHash()
			}
		}
	}()
}

func (g *HashGenerator) GetHash() string {
	return <-g.hashes
}
