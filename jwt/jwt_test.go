package jwt

import (
	"testing"
	"time"
)

func TestSigned(t *testing.T) {
	j := NewJwt()
	claims := IMClaims{}
	claims.Mid = 123456
	claims.Accepts = []int32{1}
	claims.ExpiresAt = time.Now().Add(3600 * time.Second).Unix()
	t.Log(j.Signed(claims))
}

func TestParse(t *testing.T) {
	var token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTE3MDA2OTYsIm1pZCI6MTIzNDU2LCJhY2NlcHRzIjpbMV19.7TkQeza_jWtny2BvK_PLfGMebLI7Czf91kNQIULwc-c"

	j := NewJwt()
	c, _ := j.Parse(token, &IMClaims{})

	if v, ok := c.(IMClaims); ok {
		t.Log(v.Mid)
	}
}
