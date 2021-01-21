package jwt

import (
	"crypto/rsa"

	"github.com/dgrijalva/jwt-go"
)

const defaultSignKey = "test"

// Jwt .
type Jwt struct {
	// Realm name to display to the user. Required.
	Realm string

	// signing algorithm - possible values are HS256, HS384, HS512
	// Optional, default is HS256.
	SigningAlgorithm string

	// Secret key used for signing. Required.
	Key []byte

	// Private key file for asymmetric algorithms
	PrivKeyFile string

	// Public key file for asymmetric algorithms
	PubKeyFile string

	// Private key
	privKey *rsa.PrivateKey

	// Public key
	pubKey *rsa.PublicKey
}

func (j *Jwt) init() {
	if j.SigningAlgorithm == "" {
		j.SigningAlgorithm = "HS256"
	}

	if len(j.Key) == 0 {
		j.Key = []byte(defaultSignKey)
	}
	return
}

// IMClaims .
type IMClaims struct {
	jwt.StandardClaims
	Mid     int64   `json:"mid,omitempty" form:"mid"`
	Accepts []int32 `json:"accepts,omitempty" form:"accepts"`
	AppType int32   `json:"app_type,omitempty" form:"app_type"`
}

// Signed 签名数据得到完整 jwt token
func (j *Jwt) Signed(claims jwt.Claims) (token string, err error) {
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod(j.SigningAlgorithm), claims)
	return jwtToken.SignedString(j.Key)
}

// Parse 解析 jwt token
func (j *Jwt) Parse(token string, claims jwt.Claims) (c jwt.Claims, err error) {
	var (
		jwtToken *jwt.Token
	)
	jwtToken, err = jwt.ParseWithClaims(token, claims, func(*jwt.Token) (interface{}, error) {
		return j.Key, nil
	})
	if err != nil {
		return
	}
	c = jwtToken.Claims
	return
}

// NewJwt return a jwt
func NewJwt() *Jwt {
	jwt := &Jwt{}
	jwt.init()
	return jwt
}
