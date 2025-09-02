package core

import (
	"sort"
)

var registry = map[string]ScreenTest{}

func RegisterTest(t ScreenTest) {
	registry[t.Name()] = t
}

func GetTest(name string) (ScreenTest, bool) {
	t, ok := registry[name]
	return t, ok
}

func AllTests() []ScreenTest {
	list := make([]ScreenTest, 0, len(registry))
	for _, t := range registry {
		list = append(list, t)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Order() < list[j].Order()
	})

	return list
}
