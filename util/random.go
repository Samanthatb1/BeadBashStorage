// Utility functions to generate random parameters for the unit tests
package util

import (
	"math/rand"
	"strings"
	"time"
)

func init(){
	rand.Seed(time.Now().UnixNano())
}

// Generates random integers
func randomInt(min, max int64) int64{
	return min + rand.Int63n(max - min + 1)
}

// Generates random strings
func randomString(number int64) string {
	const alphabet = "abcdefghijlmnopqrstuvwxyz"
	var sb strings.Builder
	k := len(alphabet)
	var i int64 = 0

	for ; i<number; i++ {
		sb.WriteByte(alphabet[rand.Intn(k)])

	}
	return sb.String()
}

// Generates random floats
func randomFloat(min float64) float64{
	return min + rand.Float64()
}

/**********************************/

func RandomCost() float64 {
	return randomFloat(54.3)
}

func RandomLongString() string {
	return randomString(6)
}

func RandomOrders() int64 {
	return randomInt(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{"USD", "CAD", "EUR"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

func RandomID() int64 {
	return randomInt(0, 1000)
}