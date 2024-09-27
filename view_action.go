package ebui

type uiViewAction struct {
	onHover      func()
	onTapGesture func()
	onChange     func()
	onAppear     func()
	onDisappear  func()
	onDrag       func()
	onDrop       func()
}
