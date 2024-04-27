package random

import "crypto/rand"

// TODO: replace this with a proper implementation
// NanoID generates a random string of given length using URL-safe base64 encoding
func NanoId(size int) string {
	const (
		alphabet = "-.0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz"
		mask     = 63 // 6 bits of randomness
	)

	bytes := make([]byte, size)
	rand.Read(bytes)

	for i, b := range bytes {
		bytes[i] = alphabet[int(b)&mask]
	}

	return string(bytes)
}
