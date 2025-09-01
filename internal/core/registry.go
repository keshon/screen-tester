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
