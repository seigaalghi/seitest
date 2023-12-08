package tested

import "github.com/seigaalghi/seitest/interfaces"

type tested struct {
	greeter interfaces.Sample
}

type Tested interface {
	Greetings() string
}

func NewTested(greeter interfaces.Sample) Tested {
	return &tested{
		greeter: greeter,
	}
}

func (t *tested) Greetings() string {
	return t.greeter.Greet()
}

func Jamban(nama string) string {
	return nama + " Boker"
}
