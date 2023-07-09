package ccrypto

// Deterministic crypto.Reader
// overview: half the result is used as the output
// [a|...] -> sha512(a) -> [b|output] -> sha512(b)

import (
	"crypto/sha512"
	"io"
)

const DetermRandIter = 2048

func NewDetermRand(seed []byte) io.Reader {
	var out []byte
	//strengthen seed
	var next = seed
	for i := 0; i < DetermRandIter; i++ {
		next, out = hash(next)
	}
	return &determRand{
		next: next,
		out:  out,
	}
}

type determRand struct {
	next, out []byte
}

func (d *determRand) Read(b []byte) (int, error) {
	n := 0
	l := len(b)
	for n < l {
		next, out := hash(d.next)
		n += copy(b[n:], out)
		// In Golang 1.20, ecdsa.GenerateKey() introduced a function called
		// MaybeReadRand() which reads 1 byte from the determRand reader
		// with 50% chance. As a result, GenerateKey() generates
		// nondeterministic keys.
		// The following conditional check neutralizes this effect.
		if l > 1 {
			d.next = next
		}
	}
	return n, nil
}

func hash(input []byte) (next []byte, output []byte) {
	nextout := sha512.Sum512(input)
	return nextout[:sha512.Size/2], nextout[sha512.Size/2:]
}
