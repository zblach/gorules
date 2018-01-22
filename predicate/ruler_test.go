package predicate_test

import (
	"fmt"
	"testing"

	. "."
	. "./encoding"
)

func TestExample(t *testing.T) {
	rule := Not(Or(Equals("foo"), Equals("bar"), Equals("hello")))

	b, err := rule.Match(5)

	neg := Not(Equals("foo"))
	doubleneg := Not(neg)

	fmt.Printf("%s: %s, %s\n", IndentedSerial(rule), b, err)

	fmt.Printf("%s\n%s\n", StringSerial(neg), StringSerial(doubleneg))
}
