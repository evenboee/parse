package main

import (
	"fmt"
	"os"
	"time"

	"github.com/evenboee/parse"
	"github.com/evenboee/parse/env"
)

func main() {
	n, err := parse.Try[int]("123")
	fmt.Printf("%T %v %v\n", n, n, err)

	boolSlice := parse.Must[[]bool]("t,t,t,f,,")
	fmt.Printf("%T %v\n", boolSlice, boolSlice)

	myString := parse.Must[MyString]("hello")
	fmt.Printf("%T %v\n", myString, myString)

	t := parse.Must[time.Time]("2021-01-01", parse.WithTimeFormat("2006-01-02"))
	fmt.Printf("%T %v\n", t, t)

	t = parse.Must[time.Time]("1234-01-23T12:34:56Z")
	fmt.Printf("%T %v\n", t, t)

	os.Setenv("EXPIRES_IN", "1h")
	expiresIn := env.Get[*[]time.Duration]("EXPIRES_IN", "")
	fmt.Printf("%T %v\n", expiresIn, expiresIn)
}

type MyString string

func (m *MyString) UnmarshalString(text string) error {
	*m = MyString(text) + "!"
	return nil
}
