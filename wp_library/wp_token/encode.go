package wp_token

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// 加密函数
func encrypt(plaintext []byte, key []byte) (string, error) {
	// 创建一个新的 AES 块
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 创建一个加密模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 生成一个随机 nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// 使用 AES-GCM 加密
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// 解密函数
func decrypt(ciphertext []byte, key []byte) ([]byte, error) {

	// 创建一个新的 AES 块
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 创建一个解密模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 获取 nonce 的长度
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	// 提取 nonce
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// 解密数据
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
