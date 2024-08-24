package test

import (
	"fmt"
	"testing"
)

func Test_for(t *testing.T) {

	ch := make(chan int, 1)
	ch <- 9
	fmt.Println(<-ch)
}

func f(ch chan<- int) {}
