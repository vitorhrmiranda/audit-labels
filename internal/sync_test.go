package internal_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	audit "github.com/vitorhrmiranda/audit/internal"
)

type item struct {
	id, url, code string
}

func (i item) ID() string   { return i.id }
func (i item) URL() string  { return i.url }
func (i item) Code() string { return i.code }

func TestPerform(t *testing.T) {
	type args struct {
		item
	}

	tests := map[string]struct {
		args     args
		wantCode string
	}{
		"should return a list of PDFs": {
			wantCode: "JA1BB8C4F8FB4F0F",
			args: args{
				item{
					id:   "4630683",
					url:  "https://envio.enjoei.com.br/etiqueta/jadlog/d7ff0ff8876f45e1e69d38159a789065fa1fd061c18b4c3602aed2074c481ba5/etiqueta_4630683.pdf",
					code: "JA1BB8C4F8FB4F0F",
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var items []audit.Input
			items = append(items, tt.args.item)

			pdfs := audit.Perform(items)

			want := []audit.PDF{{Code: tt.wantCode, ID: tt.args.id}}
			assertPDF(t, pdfs, want)
		})
	}
}

func assertPDF(t *testing.T, got []audit.PDF, want []audit.PDF) {
	assertion := assert.New(t)

	for i, g := range got {
		assertion.Equal(g.ID, want[i].ID)
		assertion.Equal(g.Code, want[i].Code)
	}
}
