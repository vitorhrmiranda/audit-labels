package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"log"
	"os"
	"regexp"

	"github.com/vitorhrmiranda/audit/entities"
	"github.com/vitorhrmiranda/audit/internal"
	"github.com/vitorhrmiranda/audit/persistence"
)

//go:embed input.json
var file []byte

var methods = map[string]func(){
	"sync":  sync,
	"async": async,
	"meta":  meta,
}

func main() {
	var method string
	flag.StringVar(&method, "method", "async", "sync, async or seller")
	flag.Parse()

	if f, ok := methods[method]; ok {
		f()
		return
	}

	log.Fatal("invalid method")
}

func sync() {
	log.Println("Starting sync...")

	var items []entities.Input
	if err := json.Unmarshal(file, &items); err != nil {
		log.Fatal(err)
	}

	var i []internal.Input
	for j := 0; j < len(items); j++ {
		i = append(i, items[j])
	}

	log.Printf("COUNT: %d", len(i))
	pdf := internal.Perform(i)
	db := persistence.New()

	if err := db.Table("pdf").Create(pdf).Error; err != nil {
		log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile).Println(err)
	}
}

func async() {
	log.Println("Starting async...")
	var items []entities.Input
	if err := json.Unmarshal(file, &items); err != nil {
		log.Fatal(err)
	}

	var i []internal.Input
	for j := 0; j < len(items); j++ {
		i = append(i, items[j])
	}

	log.Printf("COUNT: %d", len(i))
	pdfs := internal.PerformAsync(i)
	db := persistence.New()

	for pdf := range pdfs {
		if err := db.Table("pdf").Create(pdf).Error; err != nil {
			log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile).Printf("ID: %s, Order: %s, %v", pdf.ID, pdf.Code, err)
		}
	}
}

func meta() {
	db := persistence.New()

	pdfs := []internal.PDF{}
	db.Table("pdf").Where("error > 0").Find(&pdfs)

	for _, pdf := range pdfs {
		text := []byte(pdf.PlainText)

		code := regexp.MustCompile(`(Phone)(.*)+`).Find(text)
		pdf.Phone = string(code[5:])

		pedido := regexp.MustCompile(`(Pedido)(.*)+`).Find(text)
		pdf.Order = string(pedido[7:])

		b := regexp.MustCompile(`((.*)(\n)){4}(DECLARACAO)`).Find(text)
		buyer := regexp.MustCompile(`\r?\n`).ReplaceAll(b, []byte(" "))
		pdf.Buyer = string(buyer[:len(buyer)-10])

		s := regexp.MustCompile(`(Pedido )(\d+)((.*)(\n)){5}`).Find(text)
		seller := regexp.MustCompile(`\r?\n`).ReplaceAll(s, []byte(" "))
		seller = regexp.MustCompile(`(Pedido )(\d+)`).ReplaceAll(seller, []byte(""))
		pdf.Seller = string(seller)

		db.Table("pdf").Save(pdf)
	}
}
