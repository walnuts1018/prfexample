package random

import (
	"crypto/rand"
	mathrand "math/rand/v2"

	"github.com/cockroachdb/errors"
)

const UpperLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const LowerLetters = "abcdefghijklmnopqrstuvwxyz"
const Numbers = "0123456789"
const Symbols = "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
const Alphabets = UpperLetters + LowerLetters
const Alphanumeric = Alphabets + Numbers
const AlphanumericSymbols = Alphanumeric + Symbols

type Random interface {
	InsecureString(length uint, base string) string
	SecureString(length uint, base string) (string, error)
}

type random struct{}

func New() Random {
	return random{}
}

func (r random) InsecureString(length uint, base string) string {
	runes := []rune(base)
	result := make([]rune, length)
	for i := range result {
		result[i] = runes[mathrand.IntN(len(runes))]
	}
	return string(result)
}

func (r random) SecureString(length uint, base string) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", errors.Wrap(err, "failed to read random")
	}

	var result string
	for _, v := range b {
		result += string(base[int(v)%len(base)])
	}
	return result, nil
}

type dummy struct{}

func NewDummy() Random {
	return dummy{}
}

func (d dummy) InsecureString(length uint, base string) string {
	return "dummy"
}

func (d dummy) SecureString(length uint, base string) (string, error) {
	return "dummy", nil
}
