package main

import (
	"go-novel/utils"
	"log"
)

//const (
//	appAESKey = "5742306e4d5a48586c784e6e644f5265" // 约定的 AES key，使用十六进制表示
//)
//
//func encrypt(data string) (string, error) {
//	key, err := hex.DecodeString(appAESKey)
//	if err != nil {
//		return "", err
//	}
//
//	plaintext := []byte(data)
//
//	block, err := aes.NewCipher(key)
//	if err != nil {
//		return "", err
//	}
//
//	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
//	iv := ciphertext[:aes.BlockSize]
//	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
//		return "", err
//	}
//
//	stream := cipher.NewCFBEncrypter(block, iv)
//	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
//
//	return hex.EncodeToString(ciphertext), nil
//}
//
//func decrypt(encryptedData string) (string, error) {
//	key, err := hex.DecodeString(appAESKey)
//	if err != nil {
//		return "", err
//	}
//
//	ciphertext, err := hex.DecodeString(encryptedData)
//	if err != nil {
//		return "", err
//	}
//
//	block, err := aes.NewCipher(key)
//	if err != nil {
//		return "", err
//	}
//
//	if len(ciphertext) < aes.BlockSize {
//		return "", fmt.Errorf("ciphertext too short")
//	}
//	iv := ciphertext[:aes.BlockSize]
//	ciphertext = ciphertext[aes.BlockSize:]
//
//	stream := cipher.NewCFBDecrypter(block, iv)
//	stream.XORKeyStream(ciphertext, ciphertext)
//
//	return string(ciphertext), nil
//}

func main() {
	//data := "Hello, Golang AES Encryption and Decryption!"
	//aa, _ := encrypt(data)
	//log.Println(aa)
	//aa = "631835bb182cf46e30534b39444f45674352357069653565bae80a6f8f069fb44a24703f3b044b51db0ffa7bd10394d54a68304da5fa127947e67bb1cea5e207a86dcaf4464c2becd680a4f6897a5b242b5919"
	//log.Println(decrypt(aa))
	aa := "{\n    //页数\n    \"page\": 1,\n    //每页条数\n    \"size\": 10\n}"
	aa = "hello world"
	aa = `{
			"account_id": 1
		}`
	bb, err := utils.AesEncryptByCFB(utils.ApiAesKey, aa)
	log.Println(bb, err)
	bb = "610f6e208f7c68fe587758704b624662563472644c6e3766d60ed9ad6a5546c8d6425508b3eae661864d584c5154e9e746"
	cc, _ := utils.AesDecryptByCFB(utils.ApiAesKey, bb)
	log.Println(cc)
}
