package example

import (
	"encoding/json"
	"fmt"
)

//go:generate go run .. example.go Person

type Person struct {
	Name string
	Age  int
}

func Example() {
	person := NewPersonBuilder().Age(123).Name("Sherlock Homes").Build()

	b, err := json.Marshal(person)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}
