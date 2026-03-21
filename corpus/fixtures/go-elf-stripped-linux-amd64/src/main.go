package main

import "fmt"

var fixtureSentinel = []byte("goreveal fixture\x00")
var fixtureModuleSentinel = []byte("example.com/gorevealfixture\x00")

type fixtureCounter struct {
	Value int
}

type fixtureGreeter interface {
	Greet() string
}

type fixtureGreeterImpl struct {
	Label string
}

func (g fixtureGreeterImpl) Greet() string {
	return g.Label
}

//go:noinline
func helperAdd(a, b int) int {
	return a + b
}

//go:noinline
func helperBanner() string {
	return "goreveal fixture"
}

func main() {
	counter := fixtureCounter{Value: helperAdd(20, 22)}
	var greeter fixtureGreeter = fixtureGreeterImpl{Label: helperBanner()}
	fmt.Println(greeter.Greet(), counter.Value, len(fixtureSentinel), len(fixtureModuleSentinel))
}
