package main

import (
	"fmt"
	"strings"

	"github.com/namlulu/learngo/hello"
)

func multiply(a int, b int) int {
	return a * b
}

func lenAndUpper(name string) (length int, uppercase string) {
	defer fmt.Println("lenAndUpper is done") // defer structure is stacked and executed at the end of the function
	length = len(name)
	uppercase = strings.ToUpper(name)
	return
}

func reaptMe(word string, count int) string {
	return strings.Repeat(word, count)
}

func superAdd(numbers ...int) int {
	total := 0
	for _, number := range numbers {
		total += number
	}
	return total
}

func canIDrink(age int) bool {
	if koreanAge := age + 2; koreanAge >= 20 {
		return true
	} else {
		fmt.Println("Korean age:", koreanAge)
		return false
	}
}

func canIDrinkSwitch(age int) bool {
	switch koreanAge := age + 2; koreanAge {
	case 10:
		fmt.Println("Too young to drink")
		return false
	case 15:
		fmt.Println("Too young to drink but you can drink soju")
		return false
	case 18:
		fmt.Println("You can drink beer")
		return true
	case 20:
		fmt.Println("You can drink anything")
		return true
	default:
		if koreanAge > 20 {
			fmt.Println("You can drink anything")
			return true
		} else {
			fmt.Println("Too young to drink")
			return false
		}
	}
}

type person struct {
	name         string
	age          int
	favoriteFood []string
}

func main() {
	name := "namlulu"
	println(name)

	hello.SayHello()

	result := multiply(2, 3)
	println("2 * 3 =", result)

	length, upperName := lenAndUpper(name)
	println("length:", length, "upperName:", upperName)

	repeated := reaptMe("na", 5)
	println(repeated)

	total := superAdd(1, 2, 3, 4, 5)
	println("total:", total)

	canDrink := canIDrink(18)
	println("Can I drink?", canDrink)

	canDrinkSwitch := canIDrinkSwitch(18)
	println("Can I drink?", canDrinkSwitch)

	a := 2
	b := &a
	a = 5
	println("a:", a, "b:", *b)

	names := [3]string{"namlulu", "nico", "lynn"}
	names[2] = "dal"
	for _, name := range names {
		fmt.Println(name)
	}

	namesSlice := []string{"namlulu", "nico", "lynn"}
	namesSlice = append(namesSlice, "dal")
	for _, name := range namesSlice {
		fmt.Println(name)
	}

	namlulu := map[string]string{"name": "namlulu", "age": "20"}
	for key, _ := range namlulu {
		fmt.Println(key)
	}

	favoriteFood := []string{"kimchi", "ramen"}
	namluluChar := person{name: "namlulu", age: 20, favoriteFood: favoriteFood}
	fmt.Println(namluluChar)

}
