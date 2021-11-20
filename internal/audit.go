package internal

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"

	"code.sajari.com/docconv"
)

type Input interface {
	ID() string
	Code() string
	URL() string
}

type PDF struct {
	ID        string `json:"id" gorm:"text"`
	Code      string `json:"code" gorm:"text"`
	PlainText string `json:"plain_text" gorm:"text"`
	Error     int    `json:"error"`
}

type percent float64

func (p percent) String() string {
	return fmt.Sprintf("Progress: %.2f%%", p*100)
}

func download(url string, w io.Writer) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(w, resp.Body)
	return err
}

func Perform(items []Input) []PDF {
	var pdfs []PDF

	total := len(items)

	log.Print("START")

	for row, item := range items {
		log.Println(percent(float64(row+1) / float64(total)))

		b := bytes.Buffer{}
		if err := download(item.URL(), &b); err != nil {
			log.Println(err)
			continue
		}

		res, err := docconv.Convert(&b, "application/pdf", true)
		if err != nil {
			log.Println(err)
			continue
		}

		code := regexp.MustCompile(`#([A-Z]|[0-9])+`).Find([]byte(res.Body))
		orderCode := string(code[1:])

		equal := 0
		if orderCode != item.Code() {
			equal = 1
		}

		pdfs = append(pdfs, PDF{
			ID:        item.ID(),
			Code:      orderCode,
			PlainText: res.Body,
			Error:     equal,
		})
	}

	log.Print("FINISH")
	return pdfs
}
