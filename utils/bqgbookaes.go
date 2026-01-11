package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"hash/crc32"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

func Random(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

// 解密传入对应的数据内容信息
type ReqGeneralDecrypt struct {
	Content string `form:"content" json:"content" binding:"required"` //传入解码的内容
}

func RandFloats(min, max float64, n int) []float64 {
	rand.Seed(time.Now().UnixNano())
	res := make([]float64, n)
	for i := range res {
		res[i] = min + rand.Float64()*(max-min)
	}
	return res
}

func RandString(l int, t int) string {
	var (
		result bytes.Buffer
		temp   byte
	)
	switch t {
	case 1: //数字字母混合
		var (
			iniByte []byte = []byte{
				49, 50, 51, 52, 53, 54, 55, 56, 57, 97, 98, 99, 100, 101, 102, 103, 104, 106, 107, 109, 110, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122,
			}
			iniByteLen int = len(iniByte)
		)
		for i := 0; i < l; {
			if iniByte[Random(0, iniByteLen)] != temp {
				temp = iniByte[Random(0, iniByteLen)]
				result.WriteByte(temp)
				i++
			}
		}
	case 2: //前面字母，后面数字
		var (
			letterByte  []byte = []byte{97, 98, 99, 100, 101, 102, 103, 104, 106, 107, 109, 110, 112, 113, 114, 115, 116, 117, 119, 115, 121}
			letterLen   int    = 21
			numeralByte []byte = []byte{49, 50, 51, 52, 53, 54, 55, 56, 57}
			numeralLen  int    = 9
		)
		temp = letterByte[Random(0, letterLen)]
		result.WriteByte(temp)
		temp = letterByte[Random(0, letterLen)]
		result.WriteByte(temp)
		for i := 0; i < l-2; {
			if numeralByte[Random(0, numeralLen)] != temp {
				temp = numeralByte[Random(0, numeralLen)]
				result.WriteByte(temp)
				i++
			}
		}
	}
	return strings.ToTitle(result.String())
}

func GetEncryptPath(str string) (path string) {
	hash := md5.New()
	hash.Write([]byte(str))
	encrypt := hex.EncodeToString(hash.Sum(nil))

	startString := string([]byte(encrypt)[:2])
	endString := string([]byte(encrypt)[2:4])

	path = startString + "/" + endString + "/" + str
	return
}

func HttpPut(data string, url string) (body []byte, err error) {
	payload := strings.NewReader(data)
	req, err := http.NewRequest("PUT", url, payload)
	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	body, err = ioutil.ReadAll(res.Body)
	_ = res.Body.Close()
	return
}

func HttpGet(url string) (httpContent string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	httpContent = string(body)
	_ = resp.Body.Close()
	return
}

func HttpPost(params url.Values, url string) (body []byte, err error) {
	resp, err := http.PostForm(url, params)
	if err != nil {
		return
	}
	body, err = ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	return
}
func HttpPostByte(params []byte, url, contentType string) (body []byte, err error) {
	resp, err := http.Post(url, contentType, strings.NewReader(string(params)))
	if err != nil {
		return
	}
	body, err = ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	return
}
func HttpDel(params []byte, url string) (body []byte, err error) {
	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ = ioutil.ReadAll(res.Body)
	return
}

func GetBookImagePath(image string) (imagePath string) {
	if image == "" {
		image = "default/book.jpg"
	}
	imagePath = image
	return
}

func GetBookSourceImagePath(sourceId string, image string) (imagePath string) {
	if image == "" {
		image = "default/book.jpg"
	}

	imagePath = sourceId + "/" + GetEncryptPath(image)
	return
}

func StringConvertJson(dataString string) (dataJson interface{}) {
	_ = json.Unmarshal([]byte(dataString), &dataJson)
	return
}

func GetPassword(password string) (hashPwd string, err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	hashPwd = string(hash)
	return
}

func CheckPassword(password string, hasPassword string) (err error) {
	err = bcrypt.CompareHashAndPassword([]byte(hasPassword), []byte(password))
	return
}

func StringLength(str string) (count int) {
	//count = strings.Count(str, "") - 1
	count = utf8.RuneCountInString(str)
	return
}

func SubString(source string, start int, end int) string {
	var r = []rune(source)
	length := len(r)

	if start < 0 || end > length || start > end {
		return ""
	}

	if start == 0 && end == length {
		return source
	}

	return string(r[start:end])
}

func Md5String(str string) (md5String string) {
	return Md5Str([]byte(str))
}

func Md5Str(b []byte) (md5String string) {
	has := md5.Sum(b)
	md5String = fmt.Sprintf("%x", has)
	return
}

func IdToCrc32(id int, scope int) (int, error) {
	idstr := strconv.Itoa(id)
	sprints := fmt.Sprint(crc32.ChecksumIEEE([]byte(idstr)))
	num, err := strconv.Atoi(sprints)
	if err != nil {
		return 0, err
	}
	return num % scope, nil
}

// 获取采集章节路径
func GetCollectBookCatalogPath(siteId string, crawlerBookId int) (path string) {
	idToCrc32, _ := IdToCrc32(crawlerBookId, 500)
	idToCrc32Str := strconv.Itoa(idToCrc32)
	encryptPath := GetEncryptPath(strconv.Itoa(crawlerBookId))
	if siteId != "" {
		path = siteId + "/"
	}
	path += idToCrc32Str + "/" + encryptPath
	return
}

func GetCollectBookDebugPath(bookName, auther, sourceBookId string) (path string) {
	sourceBookId = strings.ReplaceAll(sourceBookId, "/", "_")
	return bookName + "_" + auther + "_" + sourceBookId
}
func GetOldCollectBookDebugPath(bookName, auther string) (path string) {
	return bookName + "_" + auther
}

func GetListPath(BookId, MainId int, Source string) (urlPath string) {
	if MainId > 0 {
		BookId = MainId
	}
	urlPath = GetCollectBookCatalogPath(Source, BookId) + ".html"
	return
}
func GetListReloadPath(BookId, MainId int, Source string) (urlPath string) {
	if MainId > 0 {
		BookId = MainId
	}
	urlPath = GetCollectBookCatalogPath(Source, BookId) + "_reload.html"
	return
}

func GetBookListsPath(bookListsId int) (path string) {
	path = strconv.Itoa(bookListsId/1000) + "/" + strconv.Itoa(bookListsId) + ".html"
	return
}

func GetTableNumber(prefix string, id int, scope int) (tableName string) {
	tableName = prefix + strconv.Itoa(id%scope)
	return
}

func GetMiddlewarePackageDataConf(dataConfInterface interface{}) (dataConf int8) {
	dataConf, _ = dataConfInterface.(int8)
	if dataConf == 0 {
		dataConf = 1
	}
	return
}

type FilBaseEncryptData struct {
	Ver     int    `json:"ver"`
	Content string `json:"content"`
}

// 加密引用处理
func ComicApiEncrypt(ver int, content []byte) (data FilBaseEncryptData, err error) {
	switch ver {
	case EncryptVer0:
		data = FilBaseEncryptData{
			Content: base64.StdEncoding.EncodeToString(content),
			Ver:     EncryptVer0,
		}
	case EncryptVer1:
		encryptV1, err := ComicApiEncryptV1(content)
		if err == nil {
			data = FilBaseEncryptData{
				Content: encryptV1,
				Ver:     EncryptVer1,
			}
		}
	default:
		err = errors.New(fmt.Sprintln("ver ", ver, "加密版本不支持"))
	}
	return
}

// 默认用这个值来进行加密
func ComicApiEncryptV1(content []byte) (string, error) {
	key := RandString(16, 1)
	aesKey := SHA256String(key)
	ivStr := RandString(16, 1)
	ivMd5 := Md5String(ivStr)
	iv := []byte{}
	for i := 0; i < 16; i++ {
		item := ^(int(ivMd5[i]) ^ int(ivStr[i]))
		if item < 0 {
			item = item + 256
		}
		iv = append(iv, byte(item))
	}
	aescbcEncrypt, err := AESCBCPKCS5PaddingEncrypt(content, aesKey, iv)
	if err != nil {
		return "", err
	}
	encrypt := append([]byte{}, []byte(key)...)
	encrypt = append(encrypt, aescbcEncrypt...)
	encrypt = append(encrypt, []byte(ivStr)...)
	base64str := base64.StdEncoding.EncodeToString(encrypt)
	return base64str, nil
}

// 解密
func ComicApiDecrypt(ver int, base64str string) (string, error) {
	switch ver {
	case EncryptVer0:
		//不需要解密
		return base64str, nil
	case EncryptVer1:
		//去到对应的解密函数
		return ComicApiDecryptV1(base64str)
	default:
		return "", errors.New(fmt.Sprintln("ver ", ver, "暂没有解密方法"))
	}
}

func ComicApiDecryptV1(content string) (string, error) {

	//把接口返回的base64解密
	contentByte, err := base64.StdEncoding.DecodeString(content)
	//fmt.Println(contentByte)
	//return "", nil
	content = string(contentByte)
	//前16位是key
	key := content[0:16]
	//拿key去SHA256后的到的值就是aes的key，例子这里的SHA256返回值是byte
	aesKey := SHA256String(key)
	//中间部分是加密内容
	encrypt := content[16 : len(content)-16]
	//最后16位是iv
	ivStr := content[len(content)-16:]
	//iv拿去md5
	ivMd5 := Md5String(ivStr)
	//fmt.Println(ivMd5)
	//return "", nil

	iv := []byte{}
	for i := 0; i < 16; i++ {
		//这里拿了了ivMd5[i]转int类型 和 ivStr[i] 做按位异或运算 再按位取反运算
		//go的取反是^int,其他语言应该是~int  ~(int ^ int)
		//java 的byte是int8类型，需要转uint8来运算
		item := ^(int(ivMd5[i]) ^ int(ivStr[i]))
		//运算后，如果小于0 需要+ 256
		if item < 0 {
			item = item + 256
		}

		//添加到[]byte里，就是aes的iv了iv
		iv = append(iv, byte(item))
	}
	//aes/cbc/PKCS5Padding 方式解密
	decrypt, err := AESCBCPKCS5PaddingDecrypt([]byte(encrypt), aesKey, iv)
	if err != nil {
		return "", err
	}
	return string(decrypt), nil
}

func AESCBCPKCS5PaddingDecrypt(data, key, iv []byte) (content []byte, err error) {
	if len(data) == 0 {
		return content, errors.New("AESCBCPKCS5PaddingDecrypt len(data) == 0")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(data))
	blockMode.CryptBlocks(origData, data)
	content = PKCS5UnPadding(origData)
	return
}

func AESCBCPKCS5PaddingEncrypt(data, key, iv []byte) (content []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	origData := PKCS5Padding(data, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, iv)
	content = make([]byte, len(origData))
	blockMode.CryptBlocks(content, origData)
	return
}

func SHA256String(str string) []byte {
	return SHA256([]byte(str))
}

func SHA256(b []byte) []byte {
	s := sha256.New()
	s.Write(b)
	return s.Sum(nil)
}

func SHA1String(str string) []byte {
	return SHA1([]byte(str))
}
func SHA1(b []byte) []byte {
	s := sha1.New()
	s.Write(b)
	return s.Sum(nil)
}
