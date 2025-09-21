package main

import "github.com/namlulu/dict/dictionary"

func main() {
	myDict := dictionary.Dictionary{}

	err := myDict.Add("hello", "A greeting2")
	if err != nil {
		println(err.Error())
	}

	err = myDict.Update("hello2", "A greeting2")
	if err != nil {
		println(err.Error())
	}

	definition, err := myDict.Search("hello")
	if err != nil {
		println(err.Error())
	} else {
		println(definition)
	}

	err = myDict.Delete("hello2")
	if err != nil {
		println(err.Error())
	}
}
