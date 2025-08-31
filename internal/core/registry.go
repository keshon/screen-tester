package core

import (
	"sort"
)

var registry = map[string]ScreenTest{}

func Register(t ScreenTest) {
	registry[t.Name()] = t
}

func Get(name string) (ScreenTest, bool) {
	t, ok := registry[name]
	return t, ok
}

func All() []ScreenTest {
	list := make([]ScreenTest, 0, len(registry))
	for _, t := range registry {
		list = append(list, t)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Order() < list[j].Order()
	})

	return list
}

func SafeRun(test ScreenTest, ctx *WindowContext) {
	bounds := ctx.Win.Bounds()
	width := int(bounds.W())
	height := int(bounds.H())
	if width == 0 || height == 0 {
		return
	}

	test.Run(ctx)
}
