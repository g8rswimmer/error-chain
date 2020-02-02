package main

import (
	"errors"
	"fmt"

	chain "github.com/g8rswimmer/error-chain"
)

type myError struct {
	code int
}

func (e *myError) Error() string {
	return fmt.Sprintf("%d", e.code)
}

func (e *myError) Is(target error) bool {
	te, ok := target.(*myError)
	if ok == false {
		return false
	}
	return e.code == te.code
}

func someFunction() error {
	ec := chain.New()
	ec.Add(errors.New("some error"))
	ec.Add(fmt.Errorf("wrap it up %w", &myError{code: 12}))
	return ec
}

func otherFunction() {
	err := someFunction()
	if errors.Is(err, &myError{code: 12}) {
		fmt.Println("got an error")
	}
}

func main() {
	fmt.Println("runing the error chain example")
	otherFunction()
}
