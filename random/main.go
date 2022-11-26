package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

func main() {
	//010318263233  04
	b := "01031826323304"
	count := 0
	fmt.Println("start time ->", time.Now().Format("2006-01-04-15-04-05"))
	go timePrint(&count)
	for {
		var a string
		for i := 0; i < 7; i++ {
			a += getRandomNum(int64(32))
		}
		if a == b {
			break
		}
		count++
	}
	fmt.Println("count:", count)
	fmt.Println("end time ->", time.Now().Format("2006-01-04-15-04-05"))
}

func timePrint(count *int) {
	tick := time.Tick(1 * time.Hour)
	for {
		select {
		case <-tick:
			fmt.Printf("now ==>%s count => %d \n", time.Now().Format("2006-01-04-15-04-05"), *count)
		}
	}
}

func getRandomNum(max int64) string {
	b := new(big.Int).SetInt64(max)
	n, err := rand.Int(rand.Reader, b)
	if err != nil {
		fmt.Println("big int ", b, "=====>", err)
		return "01"
	}
	return fmt.Sprintf("%02d", n.Int64()+1)
}
