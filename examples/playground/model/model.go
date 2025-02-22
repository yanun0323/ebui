package model

type IPerson interface {
	SetInfo(ii *Info)
}

type Person struct {
	*Info
}

type Info struct {
	Name string
}

func (i *Info) SetInfo(ii *Info) {
	if i == nil {
		i = &Info{}
	}

	i.Name = ii.Name
}
