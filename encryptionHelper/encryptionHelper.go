package encryptionHelper

import(
	"fmt"
	"crypto/sha256"
	"yukon_go/torHelper"
	"net/url"
	//"encoding/hex"
)

func VerifySignature(remoteAddress string, hashType string, signature string, body []byte) bool {
	if hashType == "sha256" {
		return verifySignature256(remoteAddress, signature, body)
	}
	return false
}

func verifySignature256(remoteAddress string, signature string, body []byte) bool {
	h := sha256.New()
	h.Write(body)
	hash := h.Sum(nil)

	remoteAddress = signaturePath(remoteAddress)
	fmt.Println(remoteAddress)

	remoteSignature, err := torHelper.PostResponse(remoteAddress, hash)
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

func SignHash(body []byte) string {
	return "hi"
}