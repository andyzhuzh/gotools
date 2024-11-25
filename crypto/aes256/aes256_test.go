package aes256

import (
	"testing"
)

func Test_encrypt(t *testing.T) {
	orgString := "This is a testing string"
	aesString := AES256Encrypt(orgString, "password")
	daesString := AES256Decrypt(aesString, "password")
	if daesString != orgString {
		t.Errorf("Error")
	}

}
