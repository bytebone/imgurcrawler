package iterating

import (
	"fmt"
	"math/rand"
)

const asciiLowercase = "abcdefghijklmnopqrstuvwxyz"
const asciiUppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const digits = "0123456789"

type RandomStringIterator struct{}

func (self *RandomStringIterator) HasNext() bool {
	return true
}

func (self *RandomStringIterator) Next() string {
	characters := fmt.Sprintf("%s%s%s", asciiLowercase, asciiUppercase, digits)
	path := make([]byte, 7)
	for i := 0; i < 7; i++ {
		path[i] = characters[rand.Intn(len(characters))]
	}
	return string(path)
}

func (self *RandomStringIterator) Close() {}
