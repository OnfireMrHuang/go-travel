package algorithm

import (
	"fmt"
	"sync"
)

var num = 1
var mutex sync.Mutex


func main() {

	//创建无缓冲的通道 c
	//ch1 := make(chan int, 0)
	//ch2 := make(chan int, 0)

	wg := sync.WaitGroup{}
	wg.Add(100)

	dataChan := make(chan int,0)

	go func() {
		for {
			<- dataChan
			if num % 2 == 0 {
				fmt.Printf("偶数: %v\n",num)
			}
			mutex.Lock()
			num=num+1
			mutex.Unlock()
			dataChan <- 1
			wg.Done()
		}
	}()

	go func() {
		for {
			<- dataChan
			if num % 2 != 0 {
				fmt.Printf("奇数: %v\n",num)
			}
			mutex.Lock()
			num = num + 1
			mutex.Unlock()
			dataChan <- 1
			wg.Done()
		}
	}()

	dataChan <- 1

	wg.Wait()
}

