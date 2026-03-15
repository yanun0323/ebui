package ebui

import "testing"

type fakeView struct{}

func (fakeView) Body(ctx ViewContext) Element {
	return Text("ok")
}

func TestRenderAndApply(t *testing.T) {
	rendered := Render(fakeView{})
	if rendered.Kind() != ElementKindView {
		t.Fatalf("expected view wrapper, got %v", rendered.Kind())
	}

	applied := Apply(Text("hello"), Foreground(White), Padding(All(8)))
	if !applied.Style().HasForeground || applied.Style().Padding.Top != 8 {
		t.Fatalf("modifier application failed: %#v", applied.Style())
	}
}
