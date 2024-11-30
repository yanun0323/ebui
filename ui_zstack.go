package ebui

import "github.com/hajimehoshi/ebiten/v2"

type zstack struct {
	view

	views []SomeView
}

func ZStack(views ...SomeView) SomeView {
	v := &zstack{
		views: views,
	}
	v.view = newView(v)
	return v
}

func (v *zstack) draw(screen *ebiten.Image) {
	v.modify()

	for _, child := range v.views {
		child.draw(screen)
	}
}
