package internal

import (
	"fmt"
	"log"
	"math"
	"sync"
	"time"

	"github.com/vitorhrmiranda/audit/progress"
)

func PerformAsync(items []Input) <-chan PDF {
	pdfs := make(chan PDF)

	go func() {
		defer close(pdfs)

		guard := make(chan struct{}, 50)
		var wg sync.WaitGroup

		log.Print("START")
		total := len(items)
		p := progress.New(total)

		for row, item := range items {
			wg.Add(1)
			guard <- struct{}{}

			p.Update(row + 1)

			go func(item Input) {
				defer wg.Done()
				defer func() { <-guard }()

				pdf, err := perform(item)
				for i := 1; err != nil && i < 5; i++ {
					log.Printf("ERROR: %s, retry=%d, id=%s", err, i, item.ID())
					time.Sleep(time.Duration(math.Pow(2, float64(i))) * time.Second)
					pdf, err = perform(item)
				}

				if err != nil {
					return
				}
				pdfs <- *pdf
			}(item)
		}

		wg.Wait()

		fmt.Println("")
		log.Print("FINISH")
	}()

	return pdfs
}
