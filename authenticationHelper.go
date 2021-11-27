package main

import(
	"fmt"
	"net/url"
	"encoding/hex"
	"crypto/ed25519"
	"crypto/rand"
)

var privateKey ed25519.PrivateKey
var publicKey ed25519.PublicKey

// verifies that a message is from the hidden service that it says it is
// gets the public key directly from the hidden service, then checks the signature of the message matches the public key
// signatures are hex utf-8 strings, (ie. binary signature -> hex string -> utf-8 binary encoding of string) so that they can be sent in headers
func VerifySignature(remoteAddress string, myAddress string, cryptoStandard string, signature string, body []byte) bool {
	// currently only ed25519 signatures are supported
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

// request public key of hidden service that sent the message
func getRemoteKey(remoteAddress string) ed25519.PublicKey {
	//todo add standard type to request
	keyAddress := remoteKeyPath(remoteAddress)
	_, keyBytes, _ := PostThroughProxy(keyAddress, []byte(""),nil)
	key := ed25519.PublicKey(keyBytes)
	return key
}

func remoteKeyPath(address string) string {
	a, _ := url.Parse(address)
	c, _ := url.Parse("./publickey")
	a = a.ResolveReference(c)
	return a.String()
}

func MyPublicKey() ed25519.PublicKey {
	return publicKey
}

// called at startup, generates random keys and stores as global variable, maybe change this???
// TODO: have keys updated more often
func InitializeKeys(){
	publicKey, privateKey, _ = ed25519.GenerateKey(rand.Reader)
}

func SignBody(body []byte, address string) string {
	signatureBytes := ed25519.Sign(privateKey, body)
	return hex.EncodeToString(signatureBytes)
}