package goutils

import (
	"fmt"
	"regexp"
	"testing"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestHello(t *testing.T) {

	fmt.Printf("%v \n", "Begin TestHello")

	want := regexp.MustCompile(`\b` + "hello from utils" + `\b`)
	msg := Hello()
	if want.MatchString(msg) == false {
		t.Fatalf(`%v does not match %v`, msg, want)
	}

	fmt.Printf("%v \n", "End TestHello")

}
