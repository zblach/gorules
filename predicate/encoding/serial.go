package encoding

import (
	"fmt"
	"strings"

	. "../types"
)

func StringSerial(p Predicate) string {
	switch ps := p.(type) {
	case *StringPredicate:
		switch ps.Type() {
		case EQUALS:
			return fmt.Sprintf("Equals(%q)", ps.Value())
		case BEGINS_WITH:
			return fmt.Sprintf("BeginsWith(%q)", ps.Value())
		case ENDS_WITH:
			return fmt.Sprintf("EndsWith(%q)", ps.Value())
		case CONTAINS:
			return fmt.Sprintf("Contains(%q)", ps.Value())
		default:
			panic("unhandled stringPredicate type")
		}
	case *CompositePredicate:
		children := make([]string, len(ps.Values()))
		for i, r := range ps.Values() {
			children[i] = StringSerial(r)
		}

		serialChild := strings.Join(children, ", ")

		switch ps.Type() {
		case NOT:
			return fmt.Sprintf("Not(%s)", serialChild)
		case AND:
			return fmt.Sprintf("And(%s)", serialChild)
		case OR:
			return fmt.Sprintf("Or(%s)", serialChild)
		default:
			panic("unhandled compositePredicate type")
		}
	case *BasePredicate:
		{
			switch ps.Type() {
			case EXISTS:
				return "Exists()"
			default:
				panic("unhandled predicate type")
			}
		}
	default:
		panic("no serialization known for this rule")
	}
}

func IndentedSerial(p Predicate) string {
	return indentedSerial(p, 0, "\t") + "\r\n"
}

func indentedSerial(p Predicate, amount int, str string) string {
	indent := strings.Repeat(str, amount)

	switch ps := p.(type) {
	case *StringPredicate, *BasePredicate:
		return indent + StringSerial(ps)

	case *CompositePredicate:
		children := make([]string, len(ps.Values()))
		for i, r := range ps.Values() {
			children[i] = indentedSerial(r, amount+1, str)
		}

		serialChild := strings.Join(children, ",\r\n")

		switch ps.Type() {
		case NOT:
			return fmt.Sprintf("%sNot(\n%s\n%s)", indent, serialChild, indent)
		case AND:
			return fmt.Sprintf("%sAnd(\n%s\n%s)", indent, serialChild, indent)
		case OR:
			return fmt.Sprintf("%sOr(\n%s\n%s)", indent, serialChild, indent)
		default:
			panic("unhandled compositePredicate type")
		}
	default:
		panic("no serialization known for this rule")
	}
}
