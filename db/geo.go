package db

import (
	"github.com/oschwald/geoip2-golang"
	"go-novel/global"
	"log"
)

func InitGeoReadre() {
	// 兼容历史目录：GeoLite2-City_20221007 与当前仓库常见目录 GeoLite2-City
	cityFileCandidates := []string{
		"public/resource/GeoLite2-City/GeoLite2-City.mmdb",
		"public/resource/GeoLite2-City_20221007/GeoLite2-City.mmdb",
	}
	for _, cityFile := range cityFileCandidates {
		cityDb, err := geoip2.Open(cityFile)
		if err == nil {
			global.GeoCityReader = cityDb
			break
		}
	}
	if global.GeoCityReader == nil {
		log.Printf("GeoLite2 City 数据库加载失败（将跳过 IP 解析）：candidates=%v", cityFileCandidates)
		return
	}

	asnFileCandidates := []string{
		"public/resource/GeoLite2-City/GeoLite2-ASN.mmdb",
		"public/resource/GeoLite2-City_20221007/GeoLite2-ASN.mmdb",
	}
	for _, asnFile := range asnFileCandidates {
		asnDb, err := geoip2.Open(asnFile)
		if err == nil {
			global.GeoAsnReader = asnDb
			break
		}
	}
	if global.GeoAsnReader == nil {
		log.Printf("GeoLite2 ASN 数据库加载失败（将跳过 ASN 解析）：candidates=%v", asnFileCandidates)
		return
	}
}
