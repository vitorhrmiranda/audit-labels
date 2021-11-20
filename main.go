package main

import (
	_ "embed"
	"encoding/json"
	"log"

	"github.com/vitorhrmiranda/audit/entities"
	"github.com/vitorhrmiranda/audit/internal"
	"github.com/vitorhrmiranda/audit/persistence"
)

//go:embed 2021-11-20T21_12_15.json
var file []byte

func main() {
	sync()
}

func sync() {
	// 2021/11/20 18:22:37 START
	// 2021/11/20 18:28:06 FINISH
	var items []entities.Input
	if err := json.Unmarshal(file, &items); err != nil {
		log.Fatal(err)
	}

	var i []internal.Input
	for j := 0; j < 1000; j++ {
		i = append(i, items[j])
	}

	pdf := internal.Perform(i)
	db := persistence.New()

	if err := db.Table("pdf").Create(pdf).Error; err != nil {
		log.Fatal(err)
	}
}

func async() {

}
