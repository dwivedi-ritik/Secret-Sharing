package main

// import "github.com/redis/go/v9"

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func main() {

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	err := client.Set(context.Background(), "Hi", "World", 0).Err()
	if err != nil {
		panic(err)
	}
	// t := time.Duration(60000000000)
	cmd := client.Get(context.Background(), "Hi")
	if errors.Is(cmd.Err(), redis.Nil) {
		fmt.Println("Key isn't present")
	} else {
		// fmt.Printf("%v", cmd.Val())

	}

	fmt.Printf("%v", client.Exists(context.Background(), "Hi").Val())
}
