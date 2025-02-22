package main

import (
	"fmt"

	"github.com/yanun0323/ebui/examples/playground/model"
)

func main() {
	i := &model.Info{"Yanun"}
	p := model.IPerson(&model.Person{i})

	p.SetInfo(&model.Info{"Yanun2"})

	fmt.Printf("%+v", p.(*model.Person).Info)
}
