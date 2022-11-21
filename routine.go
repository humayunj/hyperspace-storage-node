package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

var wg sync.WaitGroup

func sayGreetings(greeting string, times int) {
	for i := 0; i < times; i++ {
		log.Println(greeting)
		d := time.Second * time.Duration(rand.Intn(5)) / 2
		time.Sleep(d)
	}
	wg.Done()
}

func Triple(n int) (r int) {
	defer func() {
		r += n
	}()
	return n + n
}

func _main() {
	// rand.Seed(time.Now().Unix())
	// log.SetFlags(0)
	// wg.Add(2)
	// go sayGreetings("hi!", 10)
	// go sayGreetings("hello!", 10)
	// time.Sleep(2 * time.Second)
	// wg.Wait()
	// runtime.GOMAXPROCS()
	// fmt.Println("R", runtime.NumCPU())
	// wg.Add(1)
	// go func() {
	// 	time.Sleep(2 * time.Second)
	// 	wg.Wait()
	// }()
	// wg.Wait()

	fmt.Println("Hi")
	defer func() {
		// v := recover()
		fmt.Println("Rec", "v")
	}()

	defer func() {
		fmt.Println("exit norm")
	}()
	panic("bye")
	fmt.Println("Unreachable")
}
