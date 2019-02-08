package main

import (
	"fmt"
	"time"
)

type message struct {
	id int64
}

func (m message) GetID() int64 {
	return m.id
}

func main() {
	q := NewQ()
	go func() {
		var i int64 = 0
		for ; i < 5; i++ {
			msg := message{id: i + 1}
			q.Append(msg)
			time.Sleep(10 * time.Millisecond)
		}

		time.Sleep(5 * time.Second)
		q.Close()
	}()

	go func() {
		for {
			msg, closed := q.Pop()
			if !closed {
				fmt.Println("closed")
				break
			}

			fmt.Printf("got message id: %d\n", msg.GetID())
		}
	}()

	time.Sleep(10 * time.Second)
}
