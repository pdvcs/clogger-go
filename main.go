package main

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func read(li *list.List) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		t := scanner.Text()
		// fmt.Printf("%d bytes: %v\n", len(t), t)
		li.PushBack(t)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}

func drain(li *list.List) {
	for li.Len() > 0 {
		el := li.Front()
		// fmt.Printf("%d bytes: %v\n", len(el.Value.(string)), el.Value)
		fmt.Printf("%v\n", el.Value)
		li.Remove(el)
	}
}

func write(li *list.List) {
	for {
		drain(li)
		time.Sleep(3 * time.Second)
		PrintMemUsage()
	}
}

func trap(sigs chan os.Signal, done chan bool, li *list.List) {
	<-sigs
	fmt.Println("# interrupted, draining...")
	drain(li)
	done <- true
}

func main() {
	li := list.New()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)
	go trap(sigs, done, li)

	go read(li)
	go write(li)

	<-done
	fmt.Println("clogger-go: exiting")

}
