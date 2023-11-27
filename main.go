package main

import (
	"context"
	"fmt"

	"github.com/jonggulee/go-ssh/internal"
)

func main() {
	fmt.Printf("Hi, This is go ssh \n\n")

	cfg, err := internal.NewConfig(context.Background())
	if err != nil {
		fmt.Println(err)
	}

	ctx := context.Background()
	table, err := internal.FindeInstances(ctx, cfg)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(table)
}
