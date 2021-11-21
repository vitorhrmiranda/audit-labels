package progress

import (
	"fmt"
	"log"
)

type percent struct {
	current, total int
}

func New(total int) *percent {
	return &percent{0, total}
}

func (p *percent) Percent() float64 {
	return float64(p.current) / float64(p.total) * 100
}

func (p *percent) Update(current int) {
	c := *p
	p.current = current

	if int(c.Percent()) >= int(p.Percent()) {
		return
	}
	log.Printf("\r%s", p)
}

func (p percent) String() string {
	return fmt.Sprintf("Progress: %.1f%%", p.Percent())
}
