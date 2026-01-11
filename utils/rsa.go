package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"os"
)

// 从文件中读取RSA key
func RSAReadKeyFromFile(filename string) []byte {
	f, err := os.Open(filename)
	var b []byte

	if err != nil {
		return b
	}
	defer f.Close()
	fileInfo, _ := f.Stat()
	b = make([]byte, fileInfo.Size())
	f.Read(b)
	return b
}

// RSA加密
func RSAEncrypt(data, publicBytes []byte) ([]byte, error) {
	var res []byte
	// 解析公钥
	block, _ := pem.Decode(publicBytes)

	if block == nil {
		return res, fmt.Errorf("无法加密, 公钥可能不正确")
	}

	// 使用X509将解码之后的数据 解析出来
	// x509.MarshalPKCS1PublicKey(block):解析之后无法用，所以采用以下方法：ParsePKIXPublicKey
	keyInit, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return res, fmt.Errorf("无法加密, 公钥可能不正确, %v", err)
	}
	// 使用公钥加密数据
	pubKey := keyInit.(*rsa.PublicKey)
	res, err = rsa.EncryptPKCS1v15(rand.Reader, pubKey, data)
	if err != nil {
		return res, fmt.Errorf("无法加密, 公钥可能不正确, %v", err)
	}
	// 将数据加密为base64格式
	return []byte(EncodeStr2Base64(string(res))), nil
}

// 对数据进行解密操作
func RSADecrypt(base64Data, privateBytes []byte) ([]byte, error) {
	var res []byte
	// 将base64数据解析
	data := []byte(DecodeStrFromBase64(string(base64Data)))
	// 解析私钥
	block, _ := pem.Decode(privateBytes)
	if block == nil {
		return res, fmt.Errorf("无法解密, 私钥可能不正确")
	}
	// 还原数据
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return res, fmt.Errorf("无法解密, 私钥可能不正确, %v", err)
	}
	res, err = rsa.DecryptPKCS1v15(rand.Reader, privateKey, data)
	if err != nil {
		return res, fmt.Errorf("无法解密, 私钥可能不正确, %v", err)
	}
	return res, nil
}

// 加密base64字符串
func EncodeStr2Base64(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// 解密base64字符串
func DecodeStrFromBase64(str string) string {
	decodeBytes, _ := base64.StdEncoding.DecodeString(str)
	return string(decodeBytes)
}

// aes加密 user.ApiKey=hex.EncodeToString(utils.AesEncryptECB([]byte(user.ApiKey))) 使用时需要更换加密key
func AesEncryptECB(origData []byte) (encrypted []byte) {
	key := []byte("0f90023fc9b9b8ff")
	cipher, _ := aes.NewCipher(generateKey(key))
	length := (len(origData) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, origData)
	pad := byte(len(plain) - len(origData))
	for i := len(origData); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted = make([]byte, len(plain))
	// 分组分块加密
	for bs, be := 0, cipher.BlockSize(); bs <= len(origData); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}

	return encrypted
}

// ak,_:=hex.DecodeString(user.ApiKey) key := utils.AesDecryptECB(ak)
func AesDecryptECB(encrypted []byte) (decrypted []byte) {
	key := []byte("0f90023fc9b9b8ff")
	cipher, _ := aes.NewCipher(generateKey(key))
	decrypted = make([]byte, len(encrypted))
	//
	if len(encrypted)%16 != 0 {
		return nil
	}
	for bs, be := 0, cipher.BlockSize(); bs < len(encrypted); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}

	trim := 0

	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}

	return decrypted[:trim]
}

func generateKey(key []byte) (genKey []byte) {
	genKey = make([]byte, 16)
	copy(genKey, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}

func getRndStr(length int) (str string, err error) {
	bytes := make([]byte, length/2)
	if _, err = rand.Read(bytes); err != nil {
		return
	}
	str = hex.EncodeToString(bytes)
	return
}

func AesEncryptByCFB(key, data string) (finalHex string, err error) {
	iv, err := getRndStr(aes.BlockSize)
	if err != nil {
		return
	}
	//iv = "9ceeb8286dd2716b"
	dataByte := []byte(data)
	keyByte := []byte(key)
	ivByte := []byte(iv)
	block, err := aes.NewCipher(keyByte)
	if err != nil {
		err = fmt.Errorf("NewCipher error %s", err.Error())
		return
	}

	ciphertext := make([]byte, len(dataByte))
	stream := cipher.NewCFBEncrypter(block, ivByte)
	stream.XORKeyStream(ciphertext, dataByte)

	encryptedHex := hex.EncodeToString(ciphertext)
	ivHex := hex.EncodeToString([]byte(iv))
	finalHex = encryptedHex[:16] + ivHex + encryptedHex[16:]
	return
}

func AesDecryptByCFB(key, cipherText string) (decryptedData string, err error) {
	strlen := len(cipherText)

	var content, iv string
	if strlen < 48 {
		if strlen <= 32 {
			err = fmt.Errorf("%v", "密文长度错误")
			return
		}
		content = cipherText[:strlen-32]
		iv = cipherText[strlen-32:]
	} else {
		content = cipherText[:16] + cipherText[48:]
		iv = cipherText[16:48]
	}

	keyByte := []byte(key)
	// 将十六进制字符串转换为字节数组
	contentByte, err := hex.DecodeString(content)
	if err != nil {
		err = fmt.Errorf("Error decoding content: %v", err.Error())
		return
	}

	ivByte, err := hex.DecodeString(iv)
	if err != nil {
		err = fmt.Errorf("Error decoding IV: %v", err.Error())
		return
	}

	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return
	}

	stream := cipher.NewCFBDecrypter(block, ivByte)
	decryptedByte := make([]byte, len(contentByte))
	stream.XORKeyStream(decryptedByte, contentByte)
	decryptedData = string(decryptedByte)
	return
}
