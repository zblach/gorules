package predicate_test

import (
	"fmt"
	"testing"

	. "."
	. "./encoding"
)

func TestExample(t *testing.T) {
	rule := Not(Or(Equals("foo"), Equals("bar"), Equals("hello")))

	b, err := rule.Match("hello")

	fmt.Println(b, err)

	b, err = rule.Match(5)

	fmt.Println(b, err)

	neg := Not(Equals("foo"))
	doubleneg := Not(neg)

	is, _ := IndentedSerial(rule)

	fmt.Printf("%s: %s, %s\n", is, b, err)

	sn, _ := String.Serialize(neg)
	snn, _ := String.Serialize(doubleneg)

	fmt.Printf("%s\n%s\n", sn, snn)
}

func TestNumeric(t *testing.T) {
	rule := NearEqual(1.0, 0.1)

	fmt.Println(Not(LessThan(10.0)).Match(10.1))
	fmt.Println(Not(Equals(10.0)).Match(10.1))

	fmt.Println(GreaterThan(10.0).Match(10.1))

	fmt.Println(String.Serialize(rule))

	for _, v := range []float64{0.89, 0.90, 0.95, 1.0, 1.05, 1.1, 1.11} {
		b, _ := rule.Match(v)
		fmt.Println(b, v)
	}

	fmt.Println(IndentedSerial(rule))
}
