package gaana

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"strconv"
)

// https://gaana.com/2440.0d0b41c2.chunk.js - function decAESCBCPKCS
// https://jsfiddle.net/zL5ojhbr/
var key = [4]int{1735995764, 593641578, 1814585892, 2004118885}
var k = make([]byte, 16)

func init() {
	// Convert to 16-byte AES key
	binary.BigEndian.PutUint32(k[0:4], uint32(key[0]))
	binary.BigEndian.PutUint32(k[4:8], uint32(key[1]))
	binary.BigEndian.PutUint32(k[8:12], uint32(key[2]))
	binary.BigEndian.PutUint32(k[12:16], uint32(key[3]))
}

func decAESCBCPKCS(message string) (string, error) {
	if len(message) == 0 {
		return "", fmt.Errorf("empty message")
	}

	// the first byte is the offset
	offset, err := strconv.ParseInt(string(message[0]), 10, 8)
	if err != nil {
		return "", fmt.Errorf("failed to parse offset: %w", err)
	}

	// iv
	iv := []byte(message[offset : offset+aes.BlockSize])

	// cipher text
	cipherText := message[offset+aes.BlockSize:]

	c, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(k)
	if err != nil {
		return "", err
	}

	if len(c)%aes.BlockSize != 0 {
		return "", fmt.Errorf("cipher text length is not a multiple of block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	decrypted := make([]byte, len(c))
	mode.CryptBlocks(decrypted, c)

	// Remove PKCS7 padding
	decrypted, err = pkcs7Unpad(decrypted)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}

func pkcs7Unpad(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty data")
	}

	padding := int(data[len(data)-1])

	if padding > len(data) {
		return nil, fmt.Errorf("invalid padding")
	}

	return data[:len(data)-padding], nil
}
