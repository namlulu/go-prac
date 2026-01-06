package main

import "fmt"

func main() {
	fmt.Print("Hello, World!")

	food := []string{"kimchi", "ramen", "pizza"}
	for _, food := range food {
		fmt.Println(food)
	}
}
