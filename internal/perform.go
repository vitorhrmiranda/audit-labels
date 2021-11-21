package internal

import (
	"bytes"
	"io"
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

func download(url string, w io.Writer) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(w, resp.Body)
	return err
}

func perform(item Input) (*PDF, error) {
	pdf := new(PDF)

	b := bytes.Buffer{}
	if err := download(item.URL(), &b); err != nil {
		return pdf, err
	}

	res, err := docconv.Convert(&b, "application/pdf", true)
	if err != nil {
		return pdf, err
	}

	var orderCode string
	if code := regexp.MustCompile(`#([A-Z]|[0-9])+`).Find([]byte(res.Body)); len(code) > 1 {
		orderCode = string(code[1:])
	}

	equal := 0
	if orderCode != item.Code() {
		equal = 1
	}

	return &PDF{
		ID:        item.ID(),
		Code:      orderCode,
		PlainText: res.Body,
		Error:     equal,
	}, nil
}
