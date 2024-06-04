package bloom

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateInstanceID generates a pseudo-random service
// instance identifier, using a service name
// suffixed by dash and a random number.
// The base64 function prefers a length that is a multiple of 4,
// otherwise it will append the = character to the end of the string to
// pad it to the correct length.
func generateInstanceID(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	t := base64.StdEncoding.EncodeToString(b)
	return t
}
