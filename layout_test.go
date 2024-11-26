package ebui

import "testing"

func TestLayoutVStackFill(t *testing.T) {
	view1 := VStack().Frame(80, 40)
	view2 := VStack().Frame(500, -1)
	view3 := VStack().Frame(20, 60)
	view0 := VStack(
		view1,
		view2,
		view3,
	)

	assert(t, view0.getSize().w, 500, "view0.w")
	assert(t, view0.getSize().h, 100, "view0.h")
	assert(t, view0.getPosition().x, 0, "view0.x")
	assert(t, view0.getPosition().y, 0, "view0.y")

	assert(t, view1.getSize().w, 80, "view1.w")
	assert(t, view1.getSize().h, 40, "view1.h")
	assert(t, view1.getPosition().x, 0, "view1.x")
	assert(t, view1.getPosition().y, 0, "view1.y")

	assert(t, view2.getSize().w, 500, "view2.w")
	assert(t, view2.getSize().h, -1, "view2.h")
	assert(t, view2.getPosition().x, 0, "view2.x")
	assert(t, view2.getPosition().y, 0, "view2.y")

	assert(t, view3.getSize().w, 20, "view3.w")
	assert(t, view3.getSize().h, 60, "view3.h")
	assert(t, view3.getPosition().x, 0, "view3.x")
	assert(t, view3.getPosition().y, 0, "view3.y")

	layout(view0, point{}, size{500, 300})

	assert(t, view0.getSize().w, 500, "view0.w")
	assert(t, view0.getSize().h, 100, "view0.h")
	assert(t, view0.getPosition().x, 0, "view0.x")
	assert(t, view0.getPosition().y, 0, "view0.y")

	assert(t, view1.getSize().w, 80, "view1.w")
	assert(t, view1.getSize().h, 40, "view1.h")
	assert(t, view1.getPosition().x, (500-80)/2, "view1.x")
	assert(t, view1.getPosition().y, 0, "view1.y")

	assert(t, view2.getSize().w, 500, "view2.w")
	assert(t, view2.getSize().h, 200, "view2.h")
	assert(t, view2.getPosition().x, 0, "view2.x")
	assert(t, view2.getPosition().y, 40, "view2.y")

	assert(t, view3.getSize().w, 20, "view3.w")
	assert(t, view3.getSize().h, 60, "view3.h")
	assert(t, view3.getPosition().x, (500-20)/2, "view3.x")
	assert(t, view3.getPosition().y, 240, "view3.y")
}

func TestLayoutVStackNotFill(t *testing.T) {
	view1 := VStack().Frame(80, 40)
	view2 := VStack().Frame(200, 100)
	view3 := VStack().Frame(20, 60)
	view0 := VStack(
		view1,
		view2,
		view3,
	)

	assert(t, view0.getSize().w, 200, "view0.w")
	assert(t, view0.getSize().h, 200, "view0.h")
	assert(t, view0.getPosition().x, 0, "view0.x")
	assert(t, view0.getPosition().y, 0, "view0.y")

	assert(t, view1.getSize().w, 80, "view1.w")
	assert(t, view1.getSize().h, 40, "view1.h")
	assert(t, view1.getPosition().x, 0, "view1.x")
	assert(t, view1.getPosition().y, 0, "view1.y")

	assert(t, view2.getSize().w, 200, "view2.w")
	assert(t, view2.getSize().h, 100, "view2.h")
	assert(t, view2.getPosition().x, 0, "view2.x")
	assert(t, view2.getPosition().y, 0, "view2.y")

	assert(t, view3.getSize().w, 20, "view3.w")
	assert(t, view3.getSize().h, 60, "view3.h")
	assert(t, view3.getPosition().x, 0, "view3.x")
	assert(t, view3.getPosition().y, 0, "view3.y")

	layout(view0, point{}, size{500, 300})

	assert(t, view0.getSize().w, 200, "view0.w")
	assert(t, view0.getSize().h, 200, "view0.h")
	assert(t, view0.getPosition().x, 0, "view0.x")
	assert(t, view0.getPosition().y, 0, "view0.y")

	assert(t, view1.getSize().w, 80, "view1.w")
	assert(t, view1.getSize().h, 40, "view1.h")
	assert(t, view1.getPosition().x, (500-80)/2, "view1.x")
	assert(t, view1.getPosition().y, 50, "view1.y")

	assert(t, view2.getSize().w, 200, "view2.w")
	assert(t, view2.getSize().h, 100, "view2.h")
	assert(t, view2.getPosition().x, (500-200)/2, "view2.x")
	assert(t, view2.getPosition().y, 90, "view2.y")

	assert(t, view3.getSize().w, 20, "view3.w")
	assert(t, view3.getSize().h, 60, "view3.h")
	assert(t, view3.getPosition().x, (500-20)/2, "view3.x")
	assert(t, view3.getPosition().y, 190, "view3.y")
}
