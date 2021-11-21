package internal

import (
	"log"
	"sync"

	"github.com/vitorhrmiranda/audit/progress"
)

func PerformAsync(items []Input) <-chan PDF {
	pdfs := make(chan PDF)

	go func() {
		defer close(pdfs)

		guard := make(chan struct{}, 25)
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
				pdfs <- *perform(item)
				<-guard
			}(item)
		}

		wg.Wait()
		log.Print("FINISH")
	}()

	return pdfs
}
