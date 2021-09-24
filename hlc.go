package hlc

import (
	"sync"
	"time"
)

const PhysicalClockMask = int64(^0x00FFFF)

type Clock struct {
	ts int64
	mu sync.Mutex
}

type tickOption struct {
	step int64
	sync int64
}

type TickOption func(opt *tickOption)

func Step(n int64) TickOption {
	return func(opt *tickOption) {
		opt.step = n
	}
}
func Sync(ts int64) TickOption {
	return func(opt *tickOption) {
		opt.sync = ts
	}
}

func max(vals ...int64) int64 {
	if len(vals) == 0 {
		return 0
	}
	max := vals[0]
	for _, v := range vals[1:] {
		if max < v {
			max = v
		}
	}
	return max
}

// now returns the physical time and logical clock
func (c *Clock) now() (pt int64, lc int64) {
	ts := c.ts
	return ts & PhysicalClockMask, ts & (^PhysicalClockMask)
}

func (c *Clock) Tick(topts ...TickOption) int64 {
	c.mu.Lock()
	defer c.mu.Unlock()

	pt, lc := c.now()
	wall := time.Now().UnixNano() & PhysicalClockMask

	opt := &tickOption{step: 1}
	for _, o := range topts {
		o(opt)
	}

	sync := Clock{ts: opt.sync}
	spt, slc := sync.now()
	// new pt and new lc
	npt, nlc := max(wall, pt, spt), lc
	if npt == pt && npt == spt {
		nlc = max(lc, slc) + opt.step
	} else if npt == pt {
		nlc = lc + opt.step
	} else if npt == spt {
		nlc = slc + opt.step
	} else {
		nlc = 0
	}

	c.ts = npt | nlc
	return c.ts
}

// Now returns the nanosecond with microsecond precise in fact
func (c *Clock) Now() int64 {
	return c.Tick(Step(0))
}

var globalC Clock

func Now() int64 {
	return globalC.Now()
}

func Tick(topts ...TickOption) int64 {
	return globalC.Tick(topts...)
}
