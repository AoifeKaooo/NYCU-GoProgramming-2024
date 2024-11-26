package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	doorStatus string
	handStatus string
	mutex      sync.Mutex
	wg         sync.WaitGroup
)

func hand() {
	mutex.Lock() // 確保只有 hand() 或 door() 其中一個能進入臨界區
	defer mutex.Unlock()
	
	handStatus = "in"
	time.Sleep(time.Millisecond * 200)
	handStatus = "out"
	wg.Done()
}

func door() {
	mutex.Lock() // 確保只有 hand() 或 door() 其中一個能進入臨界區
	defer mutex.Unlock()
	
	doorStatus = "close"
	time.Sleep(time.Millisecond * 200)
	if handStatus == "in" {
		fmt.Println("夾到手了啦！")
	} else {
		fmt.Println("沒夾到喔！")
	}
	doorStatus = "open"
	wg.Done()
}

func main() {
	for i := 0; i < 50; i++ {
		wg.Add(2)
		go door()
		go hand()
		wg.Wait()
		time.Sleep(time.Millisecond * 200)
	}
}
