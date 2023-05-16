package token

import (
	"math/rand"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var alphabet = "abcdefghijklmnopqrstuvwxyz"

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}
func TestPasetoMaker(t *testing.T) {

	text := RandomString(10)
	p := make(map[string]interface{})
	p["name"] = text

	maker, err := NewPasetoMaker(RandomString(32))
	require.NoError(t, err)

	token, payload1, err := maker.CreateToken(p)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload1)

	payload2, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.Equal(t, reflect.TypeOf(payload2).String(), reflect.TypeOf(payload1).String())

	require.Equal(t, payload2["name"], text)
}
