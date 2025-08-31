package core

type TestOptions struct {
	Brightness float64
	Extra      map[string]interface{}
}

type ScreenTest interface {
	Name() string
	Description() string
	Order() int
	Run(ctx *WindowContext)
	Options() TestOptions
}
