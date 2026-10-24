package benchmark

import (
	"crypto/rand"
	"testing"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
)

// BenchmarkPeerIDFromPublicKey measures a representative identity operation
// without depending on a specific messenger application.
func BenchmarkPeerIDFromPublicKey(b *testing.B) {
	_, publicKey, err := crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	for b.Loop() {
		if _, err := peer.IDFromPublicKey(publicKey); err != nil {
			b.Fatal(err)
		}
	}
}
