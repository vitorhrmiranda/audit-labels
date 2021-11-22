package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/vitorhrmiranda/audit/entities"
	"github.com/vitorhrmiranda/audit/internal"
	"github.com/vitorhrmiranda/audit/persistence"
)

//go:embed input.json
var file []byte

func main() {
	runAsync := flag.Bool("async", true, "Runs async requests")
	dbName := flag.String("db", "audit.db", "Database name")
	flag.Parse()

	if *runAsync {
		async(*dbName)
		return
	}

	sync(*dbName)
}

func sync(dbName string) {
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
	db := persistence.New(dbName)

	if err := db.Table("pdf").Create(pdf).Error; err != nil {
		log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile).Println(err)
	}
}

func async(dbName string) {
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
	db := persistence.New(dbName)

	for pdf := range pdfs {
		if err := db.Table("pdf").Create(pdf).Error; err != nil {
			log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile).Printf("ID: %s, Order: %s, %v", pdf.ID, pdf.Code, err)
		}
	}
}
