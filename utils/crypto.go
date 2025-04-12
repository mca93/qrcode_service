package utils

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "io"
)

var secretKey = []byte("a_very_secret_key_123!") // 24 bytes (AES-192)

func Encrypt(text string) (string, error) {
    block, err := aes.NewCipher(secretKey)
    if err != nil {
        return "", err
    }

    b := []byte(text)
    ciphertext := make([]byte, aes.BlockSize+len(b))
    iv := ciphertext[:aes.BlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return "", err
    }

    stream := cipher.NewCFBEncrypter(block, iv)
    stream.XORKeyStream(ciphertext[aes.BlockSize:], b)

    return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(cryptoText string) (string, error) {
    ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)
    block, err := aes.NewCipher(secretKey)
    if err != nil {
        return "", err
    }

    iv := ciphertext[:aes.BlockSize]
    ciphertext = ciphertext[aes.BlockSize:]
    stream := cipher.NewCFBDecrypter(block, iv)

    stream.XORKeyStream(ciphertext, ciphertext)
    return string(ciphertext), nil
}
