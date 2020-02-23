package event

import (
	"decent-ft/src/JSlike"
	"sort"
)

type Event struct {
	name string
	work func(...JSlike.Any)
}

type Bus map[string]func(...JSlike.Any)

var nameQueue []string

func (bus Bus) On(name string, work func(...JSlike.Any)) {
	bus[name] = work
	nameQueue = append(nameQueue, name)
}

func (bus Bus) Emit(name string, args ...JSlike.Any) {
	if sort.SearchStrings(nameQueue, name) == len(nameQueue) {
		return
	}
	bus[name](args)
}

func main() {
	bus := Bus{}
	bus.On("aaa", func(any ...JSlike.Any) {
		println("Aaa")
	})
	bus.Emit("aaa")
}
