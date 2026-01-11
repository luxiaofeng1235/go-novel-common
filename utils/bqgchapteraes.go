package utils

//笔趣阁解密函数的主要流程
import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
)

const (
	sKey        = "Pxga!h*e4@T8xfOm"
	ivParameter = "E&z!EHGLd$fli*8R"
)

// 加密
func PswEncrypt(src string) (encodeString string, err error) {
	key := []byte(sKey)
	iv := []byte(ivParameter)

	result, err := Aes128Encrypt([]byte(src), key, iv)
	if err != nil {
		return
	}
	encodeString = base64.RawStdEncoding.EncodeToString(result)
	return
}

// 加密
func PswEncryptForKeyIV(src string, sKey string, ivParameter string) (encodeString string, err error) {
	key := []byte(sKey)
	iv := []byte(ivParameter)

	result, err := Aes128Encrypt([]byte(src), key, iv)
	if err != nil {
		return
	}
	encodeString = base64.RawStdEncoding.EncodeToString(result)
	return
}

// 解密章节内容的接口数据
type ReqChapterDecrypt struct {
	Path      string `form:"path" json:"path" binding:"required"`             //传入相关的路径信息
	DomainUrl string `form:"domain_url" json:"domain_url" binding:"required"` //输入对应替换的域名url进行拼接不要在那边再拼接拉
}

// 解析笔趣阁的章节内容解锁页面信息
type ReqBqgContentDecrypt struct {
	Content string `form:"content" json:"content" binding:"required"`
}

type BiqugeChapterDataItem struct {
	Name      string `json:"name"`
	URL       string `json:"url"`
	IsContent bool   `json:"is_content"`
	Path      string `json:"path"`
	UpdatedAt int64  `json:"updated_at"`
}

type BiqugeChapterListItem struct {
	Code      int                     `json:"code"`
	Data      []BiqugeChapterDataItem `json:"data"`
	UpdatedAt int64                   `json:"updated_at"`
}

// 解码操作
func ParseB64String(b64String string) ([]byte, error) {
	missingPadding := len(b64String) % 4
	if missingPadding != 0 {
		b64String = b64String + strings.Repeat("=", missingPadding)
	}
	decodedBytes, err := base64.RawURLEncoding.DecodeString(b64String)
	if err != nil {
		decodedBytes, err = base64.URLEncoding.DecodeString(b64String)
		if err != nil {
			decodedBytes, err = base64.StdEncoding.DecodeString(b64String)
			if err != nil {
				log.Println("decode base64 fail:", err.Error())
				return []byte{}, err
			}
		}
	}
	return decodedBytes, nil
}

// 解密
func PswDecrypt(src string) (origString string, err error) {
	key := []byte(sKey)
	iv := []byte(ivParameter)

	var result []byte
	result, err = base64.RawStdEncoding.DecodeString(src)
	if err != nil {
		fmt.Println("base64 error11111111111111 ----")
		return "", err
	}
	//fmt.Println(cc)
	//转换解析解码对应，不足用对应的去补齐
	//result, err := ParseB64String(src)
	//if err != nil {
	//	fmt.Println("base64 error!!!!")
	//	return
	//}
	//fmt.Println(result)

	origData, err := Aes128Decrypt(result, key, iv)

	if err != nil {
		return
	}
	origString = string(origData)
	return
}

func Aes128Encrypt(origData, key []byte, IV []byte) ([]byte, error) {
	if key == nil || len(key) != 16 {
		return nil, nil
	}
	if IV != nil && len(IV) != 16 {
		return nil, nil
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, IV[:blockSize])
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func Aes128Decrypt(crypted, key []byte, IV []byte) ([]byte, error) {
	if key == nil || len(key) != 16 {
		return nil, nil
	}
	if IV != nil && len(IV) != 16 {
		return nil, nil
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, IV[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
