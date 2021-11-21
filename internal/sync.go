package internal

import (
	"log"

	"github.com/vitorhrmiranda/audit/progress"
)

func Perform(items []Input) []PDF {
	var pdfs []PDF

	total := len(items)
	p := progress.New(total)
	log.Print("START")

	for row, item := range items {
		p.Update(row + 1)

		pdf, err := perform(item)
		if err != nil {
			log.Printf("ERROR: id=%s, %s", item.ID(), err)
			continue
		}

		pdfs = append(pdfs, *pdf)
	}

	log.Print("FINISH")
	return pdfs
}
