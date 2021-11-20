package entities

import "fmt"

type Input struct {
	Timezone     string `json:"timezone"`
	ExternalCode string `json:"external_order_code"`
	Order        string `json:"order_uuid"`
	Identifier   int    `json:"id"`
	Tracking     string `json:"tracking"`
	Link         string `json:"url"`
}

func (i Input) ID() string { return fmt.Sprint(i.Identifier) }

func (i Input) URL() string { return i.Link }

func (i Input) Code() string { return i.ExternalCode }
