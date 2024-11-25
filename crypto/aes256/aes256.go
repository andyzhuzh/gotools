package aes256

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"io"
)

func AES256Encrypt(text string, passkey string) string {

	saltFlg := "SalteD_ZzH__"
	passphrase := passkey
	if len(passkey) > 48 {
		passphrase = passkey[:48]
	}
	saltlen := 48 - len(passphrase)
	if saltlen <= 0 {
		saltlen = 0
	}
	salt := make([]byte, saltlen)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		panic(err.Error())
	}

	key, iv := __DeriveKeyAndIv(passphrase, string(salt))

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}

	pad := __PKCS7Padding([]byte(text), block.BlockSize())
	ecb := cipher.NewCBCEncrypter(block, []byte(iv))
	encrypted := make([]byte, len(pad))
	ecb.CryptBlocks(encrypted, pad)

	return base64.StdEncoding.EncodeToString([]byte(saltFlg + string(salt) + string(encrypted)))
}

// Decrypts encrypted text with the passphrase
func AES256Decrypt(encrypted string, passkey string) string {
	saltFlg := "SalteD_ZzH__"
	passphrase := passkey
	if len(passkey) > 48 {
		passphrase = passkey[:48]
	}
	saltlen := 48 - len(passphrase)
	if saltlen <= 0 {
		saltlen = 0
	}
	ct, _ := base64.StdEncoding.DecodeString(encrypted)
	if len(ct) < len(saltFlg)+saltlen || string(ct[:len(saltFlg)]) != saltFlg {
		return ""
	}

	salt := ct[len(saltFlg) : len(saltFlg)+saltlen]
	ct = ct[len(saltFlg)+saltlen:]
	key, iv := __DeriveKeyAndIv(passphrase, string(salt))
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}

	cbc := cipher.NewCBCDecrypter(block, []byte(iv))
	dst := make([]byte, len(ct))
	cbc.CryptBlocks(dst, ct)
	return string(__PKCS7Trimming(dst))
}

func __PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func __PKCS7Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}

func __DeriveKeyAndIv(passphrase string, salt string) (string, string) {
	salted := ""
	dI := ""

	for len(salted) < 48 {
		md := md5.New()
		md.Write([]byte(dI + passphrase + salt))
		dM := md.Sum(nil)
		dI = string(dM[:16])
		salted = salted + dI
	}

	key := salted[0:32] // 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.
	iv := salted[32:48] // block.BlockSize

	return key, iv
}
