package progress

import (
	"fmt"
	"math"
	"time"
)

type percent struct {
	current, total int
	started        time.Time
}

func New(total int) *percent {
	return &percent{0, total, time.Now()}
}

func (p *percent) Percent() float64 {
	return float64(p.current) / float64(p.total) * 100
}

func (p *percent) Update(current int) {
	const r = 10
	c := *p
	p.current = current

	if math.Round(c.Percent()*r)/r >= math.Round(p.Percent()*r)/r {
		return
	}
	fmt.Printf("\r%s", p)
}

func (p percent) String() string {
	elapsedTime := time.Since(p.started)
	timeLeft := time.Duration(p.total * int(elapsedTime) / p.current)
	return fmt.Sprintf(
		"Progress: \033[32m%.1f%%\033[0m\tElapsed: \033[34m%v\033[0m\tLeft: \033[31m%v\033[0m",
		p.Percent(),
		elapsedTime.Round(time.Second),
		timeLeft.Round(time.Second),
	)
}
