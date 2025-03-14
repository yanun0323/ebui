package ebui

func getTypes(types ...formulaType) formulaType {
	if len(types) == 0 {
		return formulaZStack
	}
	return types[0]
}
