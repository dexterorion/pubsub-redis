package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	cli := New(&NewInput{
		RedisURL: os.Getenv("REDIS_URL"),
	})

	reply := make(chan []byte)
	err := cli.Subscribe("input/console", reply)
	if err != nil {
		log.Fatal(err)
	}

	for msg := range reply {
		fmt.Println("message: " + string(msg))
	}
}
