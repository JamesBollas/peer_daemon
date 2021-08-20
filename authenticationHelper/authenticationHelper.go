package authenticationHelper

import(
	"fmt"
	"crypto/sha256"
	"yukon_go/torHelper"
	"net/url"
	"encoding/hex"
)

func VerifySignature(remoteAddress string, myAddress string, hashType string, signature string, body []byte) bool {
	if hashType == "sha256" {
		return verifySignature256(remoteAddress, myAddress, signature, body)
	}
	return false
}

func verifySignature256(remoteAddress string, myAddress string, signature string, body []byte) bool {
	//bodyAndAddress := append(body, []byte(myAddress)...)
	hash := Hash(body, []byte(myAddress))

	remoteHashAddress := signaturePath(remoteAddress)
	//fmt.Println(remoteAddress)

	remoteSignature, err := torHelper.PostResponse(remoteHashAddress, hash)
	if err != nil{
		fmt.Println(err)
		return false
	}
	if string(remoteSignature) == signature {
		return true
	}
	return false
}

func signaturePath(address string) string {
	a, _ := url.Parse(address)
	c, _ := url.Parse("./hashcookie")
	a = a.ResolveReference(c)
	return a.String()
}

func Hash(bodies ... []byte) []byte {
	h := sha256.New()
	for _, body := range bodies{
		h.Write(body)
	}
	hash := h.Sum(nil)
	return hash
}

func SignBody(body []byte, address string) string {
	//bodyAndAddress := append(body, []byte(address)...)
	hash := Hash(body, []byte(address))
	return SignHash(hash)
}

func SignHash(hash []byte) string {
	key := Hash([]byte("secret_key"))
	signed := Hash(hash,key)
	return hex.EncodeToString(signed)
}

// func Xor(a []byte, b []byte) ([]byte){
// 	c := make([]byte, len(a))
// 	for i := 0; i<len(a); i++ {
// 		c[i] = a[i] ^ b[i]
// 	}
// 	return c
//  }