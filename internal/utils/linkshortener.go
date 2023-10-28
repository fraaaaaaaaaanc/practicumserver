package utils

import (
	"math/rand"
	"strings"
)

// Constants defining the set of characters from which a random link is generated.
// The 'alphabet' includes lowercase and uppercase letters, as well as digits.
const (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// LinkShortening is a method that generates a random shortened link.
// It utilizes a pseudo-random number and the 'alphabet' to create a random string.
// The generated string is typically 10 characters long.
// This method is used to create unique short links for original URLs.
func LinkShortening() string {
	number := rand.Uint64()
	length := len(alphabet)
	var encodedBuilder strings.Builder
	encodedBuilder.Grow(10)
	for ; number > 0; number = number / uint64(length) {
		encodedBuilder.WriteByte(alphabet[(number % uint64(length))])
	}

	return encodedBuilder.String()
}
