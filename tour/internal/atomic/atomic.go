package atomic

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Config struct {
	a []int
}

func (c *Config) T()  {
	
}

func Demo1()  {
	var v atomic.Value
	v.Store(&Config{})

	go func() {
		i := 0
		for  {
			i++
			cfg := &Config{[]int{i,i+1,i+2,i+3,i+4,i+5}}
			v.Store(cfg)
		}
	}()

	var wg sync.WaitGroup
	for n := 0; n < 4;n++ {
		wg.Add(1)
		go func() {
			for n:=0; n<100;n++ {
				cfg := v.Load().(*Config)
				cfg.T()
				fmt.Println("%v\n",cfg)
			}
		}()
	}
	wg.Wait()
}

func demo2()  {

	var cfg = &Config{}
	var l sync.RWMutex
	go func() {
		i := 0
		for  {
			i++
			l.Lock()
			cfg = &Config{[]int{i,i+1,i+2,i+3,i+4,i+5}}
			l.Unlock()
		}
	}()

	var wg sync.WaitGroup
	for n := 0; n < 4;n++ {
		wg.Add(1)
		go func() {
			for n:=0; n<100;n++ {
				l.RLock()
				fmt.Println("%v\n",cfg)
				l.Unlock()
			}
		}()
	}
	wg.Wait()
}