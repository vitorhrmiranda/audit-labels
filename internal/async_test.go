package internal_test

import (
	"testing"

	audit "github.com/vitorhrmiranda/audit/internal"
)

func TestPerformAsync(t *testing.T) {
	type args struct {
		item
	}

	tests := map[string]struct {
		args     args
		wantCode string
	}{
		"should return a producer of PDFs": {
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

			producer := audit.PerformAsync(items)

			got := []audit.PDF{}
			want := []audit.PDF{{Code: tt.wantCode, ID: tt.args.id}}

			got = append(got, <-producer)
			assertPDF(t, got, want)
		})
	}
}
