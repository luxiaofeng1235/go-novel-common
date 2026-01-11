package main

import (
	"log"
	"math/rand"
	"time"
)

func main() {

	percent := 0.85   //出现的百分比概率
	totalNum := 10000 //总次数
	ret := float64(totalNum)
	endNum := ret * percent
	var startNum float64
	startNum = ret - endNum
	log.Printf("随机的判断范围：%v", percent)
	log.Printf("total的位置开始处 %v", totalNum)
	log.Printf("20%%的位置开始处 %v", startNum)
	log.Printf("80%%的位置开始处 %v", endNum)
	// 生成多个随机数以展示结果
	n := generateWeightedRandomNumber(startNum, endNum, percent)
	log.Printf("本次随机生成的随机种子数为%v", n)
}

// 随机生成的种子数：按照一定概率来生成
func generateWeightedRandomNumber(startNum float64, endNum float64, percent float64) int {
	rand.Seed(time.Now().UnixNano())
	startNum1 := int(startNum)
	endNum1 := int(endNum)
	// 生成一个随机数，范围在 1 到 N
	// 根据权重来决定生成哪个范围的随机数
	if rand.Float64() < percent { // 85% 概率选择 8001-10000
		return rand.Intn(startNum1) + (endNum1 + 1) // 81 到 100
	} else { // 15% 概率选择 1-8000
		return rand.Intn(endNum1) + 1 // 1 到 80
	}
}
