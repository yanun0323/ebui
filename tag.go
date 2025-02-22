package ebui

type tag [20]byte

func makeTag(s string) tag {
	var v tag
	copy(v[:], s)
	return v
}

var (
	tagNone      tag = makeTag("")
	tagVStack    tag = makeTag("vstack")
	tagHStack    tag = makeTag("hstack")
	tagZStack    tag = makeTag("zstack")
	tagText      tag = makeTag("text")
	tagImage     tag = makeTag("image")
	tagButton    tag = makeTag("button")
	tagSpacer    tag = makeTag("spacer")
	tagRectangle tag = makeTag("rectangle")
)

func (v tag) String() string {
	return string(v[:])
}
