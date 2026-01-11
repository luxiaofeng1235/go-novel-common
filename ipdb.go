package main

import (
	"fmt"
	"log"

	"github.com/ipipdotnet/ipdb-go"
)

func main() {
	db, err := ipdb.NewCity("public/resource/ipv4_china.ipdb")
	if err != nil {
		log.Fatal(err)
	}

	ip := "61.165.0.5"
	info, err := db.FindMap(ip, "CN")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(info)

	fmt.Printf("Country: %s\n", info["country_name"])
	fmt.Printf("Region: %s\n", info["region_name"])
	fmt.Printf("City: %s\n", info["city_name"])
}
