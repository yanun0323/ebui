package ebui

import "github.com/hajimehoshi/ebiten/v2"

type zstack struct {
	*view

	views []SomeView
}

func ZStack(views ...SomeView) SomeView {
	v := &zstack{
		views: views,
	}
	v.view = newView(idZStack, v)
	return v
}

func (v *zstack) update(container Size) {
	v.view.update(container)

	for _, child := range v.views {
		child.update(container)
	}
}

func (v *zstack) draw(screen *ebiten.Image) {
	for _, child := range v.views {
		child.draw(screen)
	}
}
