package constants

type InputMode uint8

const (
	InputModeTraverse InputMode = iota
	InputModeSearch
	InputModeCommand
)

type ActiveMode uint8

const (
	ActiveModeDefault ActiveMode = iota
	ActiveModeAutocomplete
	ActiveModeSearchResultSelect
)
