package utils

import (
	"fmt"
	"github.com/mojocn/base64Captcha"
	"go-novel/global"
	"log"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// 返回随机字符串
func RandomString(typeStr string, number int) string {

	var letters = []byte("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	switch typeStr {
	case "salt":
		letters = []byte("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	case "code":
		letters = []byte("1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	case "rand":
		letters = []byte("1234567890abcdefghijklmnopqrstuvwxyz")
	case "upper":
		letters = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	case "lower":
		letters = []byte("abcdefghijklmnopqrstuvwxyz")
	}

	result := make([]byte, number)
	rand.Seed(time.Now().UnixNano())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

// 获取随机的书籍的位置信息
func GetBookRandPosition(total int64) int {
	if total == 0 {
		return 0
	}
	randNum := RangeNum(0, int(total)) //随机0-N的最大的一个数作为他的一个值
	sRandInt := int64(randNum)
	maxRandInt := total - sRandInt //用最大的数-当前的随机数
	if total == 1 || total == 0 {
		maxRandInt = 0 //如果当前的总数为1，就默认从位置0开始就一条数据
	}
	//设置偏移量
	offsetNum := int(maxRandInt)
	global.Requestlog.Infof("totalNum = %v,此次随机的createRand = %v,maxRandInt = %v", total, sRandInt, maxRandInt)
	return offsetNum
}

// 生成唯一订单号
func BuildOrderNo(prefix string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	number := r.Intn(99999)
	if prefix == "" {
		prefix = "C"
	}
	return prefix + fmt.Sprintf("%05d", number) + time.Now().Format("20060102150405")
}

// 生成指定范围的随机数
func RangeNum(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(max-min) + min
	return randNum
}

func RangeNumFloat(min, max int) (randomFloat float64) {
	rand.Seed(time.Now().UnixNano())
	randomFloat = rand.Float64()*(float64(max-min)) + float64(min)
	randomFloat = float64(int(randomFloat*10)) / 10
	return
}

func GetRangeAll(min, max int, filter []int) (arr []int) {
	m := make(map[int]bool)
	if len(filter) > 0 {
		for _, val := range filter {
			m[val] = true
		}
	}
	for i := min; i <= max; i++ {
		if _, isOk := m[i]; !isOk {
			arr = append(arr, i)
		}
	}
	return arr
}

// 随机验证码
func Code() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(899999) + 100000
	res := strconv.Itoa(code) //转字符串返回
	return res
}

// GetRandomCode 随机六位验证码
func GetRandomCode() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vcode := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	return vcode
}

// 获取数字验证码
func GetVerifyImgDigit() (idKeyC string, base64stringC string) {
	driver := &base64Captcha.DriverDigit{80, 240, 5, 0.7, 5}
	store := base64Captcha.DefaultMemStore
	c := base64Captcha.NewCaptcha(driver, store)
	idKeyC, base64stringC, err := c.Generate()
	if err != nil {
		log.Println(err)
	}
	return
}

// 获取字母数字混合验证码
func GetVerifyImgString() (idKeyC string, base64stringC string) {
	driver := &base64Captcha.DriverString{
		Height:          80,
		Width:           240,
		NoiseCount:      50,
		ShowLineOptions: 20,
		Length:          4,
		Source:          "abcdefghjkmnpqrstuvwxyz23456789",
		Fonts:           []string{"chromohv.ttf"},
	}
	driver = driver.ConvertFonts()
	store := base64Captcha.DefaultMemStore
	c := base64Captcha.NewCaptcha(driver, store)
	idKeyC, base64stringC, err := c.Generate()
	if err != nil {
		log.Println(err)
	}
	return
}

// 验证输入的验证码是否正确
func VerifyString(id, answer string) bool {
	driver := new(base64Captcha.DriverString)
	store := base64Captcha.DefaultMemStore
	c := base64Captcha.NewCaptcha(driver, store)
	answer = strings.ToLower(answer)
	return c.Verify(id, answer, true)
}

func GenValidateCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

func GetCion(cion int64, perc int) int {
	result := int(math.Round(float64(cion) * float64(perc) / 100))
	return result
}

func GetRandNumBookHits() (hits, hitsDay, hitsWeek, hitsMonth, shits int64, score float64, readCount, searchCount int64) {
	hits = int64(RangeNum(5000, 9999))
	hitsDay = int64(RangeNum(1, 300))
	hitsWeek = int64(RangeNum(300, 1000))
	hitsMonth = int64(RangeNum(1000, 5000))
	shits = int64(RangeNum(1, 999))
	score = RangeNumFloat(6, 9)
	readCount = int64(RangeNum(1, 10000))
	searchCount = int64(RangeNum(1, 1000))
	return
}

func GetGuestName() (guestName string) {
	uuidStr := GetUUID().String()
	str := uuidStr[:8]
	guestName = fmt.Sprintf("G-%v", str)
	return
}

func GetRandUserName() (userName string) {
	uuidStr := GetUUID().String()
	uuidStr = strings.ReplaceAll(uuidStr, "-", "")
	str := uuidStr[:16]
	userName = fmt.Sprintf("K-%v", str)
	return
}
