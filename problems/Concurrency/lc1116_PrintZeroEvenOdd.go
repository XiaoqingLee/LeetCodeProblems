package Concurrency

import (
	"fmt"
	"sync"
)

type ZeroEvenOdd struct {
	n           int
	nextToPrint int  // need protection 保存下一个要打印的非零数字
	zerosTurn   bool // need protection
	cv          *sync.Cond
}

func NewZeroEvenOdd(n int) *ZeroEvenOdd {
	if n <= 0 {
		panic("Invalid Input")
	}
	return &ZeroEvenOdd{
		n:           n,
		nextToPrint: 1,
		zerosTurn:   true,
		cv:          sync.NewCond(&sync.Mutex{}),
	}
}

func (z *ZeroEvenOdd) Zero(printNumberFunc func(a ...any) (n int, err error)) {
	for {
		z.cv.L.Lock()
		for !z.zerosTurn {
			z.cv.Wait()
		}
		if z.nextToPrint <= z.n { // Zero 被动停下来
			printNumberFunc(0)
			z.zerosTurn = !z.zerosTurn
			z.cv.Broadcast()
			z.cv.L.Unlock()
		} else {
			z.cv.L.Unlock()
			return
		}
	}
}

func (z *ZeroEvenOdd) Odd(printNumberFunc func(a ...any) (n int, err error)) {
	needTerminate := false
	for {
		z.cv.L.Lock()
		for !(!z.zerosTurn && z.nextToPrint%2 == 1) {
			z.cv.Wait()
		}
		printNumberFunc(z.nextToPrint)
		if z.nextToPrint == z.n || z.nextToPrint == z.n-1 { // Odd 和 Even 主动停下来
			needTerminate = true
		}
		z.nextToPrint += 1
		z.zerosTurn = !z.zerosTurn
		z.cv.Broadcast()
		z.cv.L.Unlock()
		if needTerminate {
			return
		}
	}
}

func (z *ZeroEvenOdd) Even(printNumberFunc func(a ...any) (n int, err error)) {
	needTerminate := false
	for {
		z.cv.L.Lock()
		for !(!z.zerosTurn && z.nextToPrint%2 == 0) {
			z.cv.Wait()
		}
		printNumberFunc(z.nextToPrint)
		if z.nextToPrint == z.n || z.nextToPrint == z.n-1 { // Odd 和 Even 主动停下来
			needTerminate = true
		}
		z.nextToPrint += 1
		z.zerosTurn = !z.zerosTurn
		z.cv.Broadcast()
		z.cv.L.Unlock()
		if needTerminate {
			return
		}
	}
}

func TestZeroEvenOdd() {
	zeo := NewZeroEvenOdd(15)
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		zeo.Even(fmt.Println)
	}()
	go func() {
		defer wg.Done()
		zeo.Odd(fmt.Println)
	}()
	go func() {
		defer wg.Done()
		zeo.Zero(fmt.Println)
	}()
	wg.Wait()
}