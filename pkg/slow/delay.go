package slow

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	rnd = rand.New(rand.NewSource(time.Now().Unix()))
)

type DelayRange struct {
	Min, Max int
}

func (r DelayRange) Wait() {
	d := time.Second

	if r.Min == r.Max {
		if r.Min < 0 {
			for {
				time.Sleep(time.Second)
			}
		} else {
			d *= time.Duration(r.Min)
		}
	} else {
		d *= time.Duration(rnd.Intn(r.Max-r.Min) + r.Min)
	}

	time.Sleep(d)
}

func (r DelayRange) ToString() string {
	if r.Min == r.Max {
		if r.Min < 0 {
			return "∞"
		}

		return fmt.Sprintf("%ds", r.Min)
	} else {
		return fmt.Sprintf("%d..%ds", r.Min, r.Max)
	}
}

type DelayOptions struct {
	Provision, Delete DelayRange
}

func (o *DelayOptions) ToString() string {
	return fmt.Sprintf("Provision: %s, Delete: %s", o.Provision.ToString(), o.Delete.ToString())
}
