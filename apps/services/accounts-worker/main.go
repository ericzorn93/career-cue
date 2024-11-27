package main

import (
	"fmt"
	"time"
)

func Hello(name string) string {
	result := "Hello " + name
	return result
}

func main() {
	time.Sleep(time.Millisecond * 500)
	fmt.Println(Hello("accounts-worker"))
}
