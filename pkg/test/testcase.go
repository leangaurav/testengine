package test

type IOType int

const (
	CONSOLE_INPUT IOType = iota 
	FILE_INPUT    IOType
)

type TestInput struct {
	Type IOType,
	Data string,
}

type TestOutput struct {
	Type IOType,
	Data string,
}

type TestCase struct {
	Input TestInput,
	Output TestOutput,
}