package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	svc := New(&NewInput{
		RedisURL: os.Getenv("REDIS_URL"),
	})

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Send message to the listeners...")
	fmt.Println("---------------------")

	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)
		err := svc.Publish("input/console", text)
		if err != nil {
			log.Fatal(err)
		}
	}
}
