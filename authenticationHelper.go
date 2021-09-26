package main

import(
	"fmt"
	//"crypto/sha256"
	//"yukon_go/torHelper"
	"net/url"
	"encoding/hex"
	"crypto/ed25519"
	//"crypto/x509"
	"crypto/rand"
)

var privateKey ed25519.PrivateKey
var publicKey ed25519.PublicKey

func VerifySignature(remoteAddress string, myAddress string, cryptoStandard string, signature string, body []byte) bool {
	if cryptoStandard == "ed25519" {
		return verifySignatureEd25519(remoteAddress, myAddress, signature, body)
	}
	return false
}

func verifySignatureEd25519(remoteAddress string, myAddress string, signature string, body []byte) bool {
	fmt.Println("got to verification")
	remoteKey := getRemoteKey(remoteAddress)
	signatureBytes, _ := hex.DecodeString(signature)
	return ed25519.Verify(remoteKey, body, signatureBytes)
}

func getRemoteKey(remoteAddress string) ed25519.PublicKey {
	//todo add standard type to request
	keyAddress := remoteKeyPath(remoteAddress)
	//fmt.Println(keyAddress)
	keyBytes, _ := PostThroughProxy(keyAddress, []byte(""),nil)
	//fmt.Println(keyBytes)
	key := ed25519.PublicKey(keyBytes)
	return key
}

func remoteKeyPath(address string) string {
	a, _ := url.Parse(address)
	c, _ := url.Parse("./publickey")
	a = a.ResolveReference(c)
	return a.String()
}

// func Hash(bodies ... []byte) []byte {
// 	h := sha256.New()
// 	for _, body := range bodies{
// 		h.Write(body)
// 	}
// 	hash := h.Sum(nil)
// 	return hash
// }

func MyPublicKey() ed25519.PublicKey {
	return publicKey
}

func InitializeKeys(){
	publicKey, privateKey, _ = ed25519.GenerateKey(rand.Reader)
}

func SignBody(body []byte, address string) string {
	signatureBytes := ed25519.Sign(privateKey, body)
	return hex.EncodeToString(signatureBytes)
}

// func SignHash(hash []byte) string {
// 	key := Hash([]byte("secret_key"))
// 	signed := Hash(hash,key)
// 	return hex.EncodeToString(signed)
// }

// func Xor(a []byte, b []byte) ([]byte){
// 	c := make([]byte, len(a))
// 	for i := 0; i<len(a); i++ {
// 		c[i] = a[i] ^ b[i]
// 	}
// 	return c
//  }