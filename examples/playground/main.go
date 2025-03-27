package main

type Person struct {
	Name  string
	Hobby Hobby
}

type Hobby struct {
	Elem Element
}

type Element struct {
	Name string
}

func main() {
	person := Person{}

	person.Hobby.Elem.Name = "test"

	println(person.Hobby.Elem.Name)
}
