package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"time"
)

func consumer(workChan chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	matched := 0
	expectedMatched := 3

	for u := range workChan {
		<-time.After(1 * time.Second)

		if u%2 == 0 {
			fmt.Println("accept job ", u)
			matched++
			if matched == expectedMatched {
				break
			}
		} else {
			fmt.Println("refuse job ", u)
		}

	}

	if matched == expectedMatched {
		fmt.Println("finish work")
		close(workChan)
	} else {
		fmt.Println("cannot get work")
	}
}

func main() {
	defer RecoverWhenWorkChanIsClosed()

	workChan := make(chan int)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go consumer(workChan, wg)

	for i := 0; i < 5; i++ {

		value, _ := rand.Int(rand.Reader, big.NewInt(1000))

		workChan <- int(value.Int64())
	}

	close(workChan)
	wg.Wait()

	fmt.Println("producer stops.")
}

// in case of panic when send on closed channel when accept > expectedMatched
func RecoverWhenWorkChanIsClosed() {
	if r := recover(); r != nil {
		if r.(error).Error() == "send on closed channel" {
			fmt.Println("producer stops.")
		} else {
			panic(r)
		}
	}
}
