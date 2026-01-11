package db

import (
	"github.com/oschwald/geoip2-golang"
	"go-novel/global"
	"log"
)

func InitGeoReadre() {
	cityFile := "public/resource/GeoLite2-City_20221007/GeoLite2-City.mmdb"
	cityDb, err := geoip2.Open(cityFile)
	if err != nil {
		log.Fatal("err", err.Error())
		return
	}
	global.GeoCityReader = cityDb

	asnFile := "public/resource/GeoLite2-City_20221007/GeoLite2-ASN.mmdb"
	asnDb, err := geoip2.Open(asnFile)
	if err != nil {
		log.Fatal("err", err.Error())
		return
	}
	global.GeoAsnReader = asnDb
	return
}
