package main

import (
	"log"
	"time"
)

func main() {

	counter := DefaultCounter()

	go func() {
		for {
			log.Println(counter.GetTotal())
			time.Sleep(time.Second)
		}
	}()

	for i := 0; i < 10; i++ {
		//go request1ms(counter)
		//go request10ms(counter)
		//go request100ms(counter)
		go request1000timeEvery1s(counter)
		log.Println("Start finish")
	}
	log.Println("Start finish all")

	select {}
}

func request1000timeEvery1s(c *Counter) {
	for {
		for i := 0; i < 1000; i++ {
			c.CountOne()
		}
		time.Sleep(time.Second)
	}
}

func request1ms(c *Counter) {
	for {
		c.CountOne()
		time.Sleep(time.Microsecond)
	}
}

func request10ms(c *Counter) {
	for {
		c.CountOne()
		time.Sleep(time.Microsecond * 10)
	}
}

func request100ms(c *Counter) {
	for {
		c.CountOne()
		time.Sleep(time.Microsecond * 100)
	}
}
