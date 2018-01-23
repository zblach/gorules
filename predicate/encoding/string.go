package encoding

import (
	"errors"
	"fmt"
	"strings"

	. "../errors"
	. "../types"
)

type strSerialize struct{}

var String strSerialize

func (strSerialize) Serialize(p Predicate) (st string, er error) {
	switch ps := p.(type) {
	case StringPredicate:
		switch ps.Type() {
		case EQUALS:
			return fmt.Sprintf("Equals(%q)", ps.Value()), nil
		case BEGINS_WITH:
			return fmt.Sprintf("BeginsWith(%q)", ps.Value()), nil
		case ENDS_WITH:
			return fmt.Sprintf("EndsWith(%q)", ps.Value()), nil
		case CONTAINS:
			return fmt.Sprintf("Contains(%q)", ps.Value()), nil
		default:
			return "", errors.New("unhandled stringPredicate type")
		}
	case NumericPredicate:
		switch ps.Type() {
		case LESS_THAN:
			return fmt.Sprintf("LessThan(%f)", ps.Value()), nil
		case EQUALS:
			return fmt.Sprintf("Equals(%f)", ps.Value()), nil
		default:
			return "", errors.New("unhandled numericPredicate type")
		}
	case CompositePredicate:
		children := make([]string, len(ps.Values()))

		errs := CompositeError{}
		for i, r := range ps.Values() {
			c, err := String.Serialize(r)
			if err != nil {
				errs.Append(err)
			}
			children[i] = c
		}
		er = errs.NilZero()

		serialChild := strings.Join(children, ", ")

		switch ps.Type() {
		case NOT:
			st = fmt.Sprintf("Not(%s)", serialChild)
		case AND:
			st = fmt.Sprintf("And(%s)", serialChild)
		case OR:
			st = fmt.Sprintf("Or(%s)", serialChild)
		default:
			return "", errors.New("unhandled composite type")
		}
		return
	case BasePredicate:
		{
			switch ps.Type() {
			case EXISTS:
				return "Exists()", nil
			default:
				return "", errors.New("unhandled predicate type")
			}
		}
	default:
		panic("no serialization known for this rule type")
	}
}

func IndentedSerial(p Predicate) (string, error) {
	return indentedSerial(p, 0, "\t")
}

func indentedSerial(p Predicate, amount int, str string) (st string, er error) {
	indent := strings.Repeat(str, amount)

	switch ps := p.(type) {
	case StringPredicate, BasePredicate, NumericPredicate:
		s, err := String.Serialize(ps)
		return indent + s, err

	case CompositePredicate:
		children := make([]string, len(ps.Values()))
		errs := CompositeError{}
		for i, r := range ps.Values() {
			s, err := indentedSerial(r, amount+1, str)
			if err != nil {
				errs.Append(err)
			}
			children[i] = s
		}
		er = errs.NilZero()

		serialChild := strings.Join(children, ",\r\n")

		switch ps.Type() {
		case NOT:
			st = fmt.Sprintf("%sNot(\n%s\n%s)", indent, serialChild, indent)
		case AND:
			st = fmt.Sprintf("%sAnd(\n%s\n%s)", indent, serialChild, indent)
		case OR:
			st = fmt.Sprintf("%sOr(\n%s\n%s)", indent, serialChild, indent)
		default:
			return "", errors.New("unhandled compositePredicate type")
		}
		return
	default:
		panic("no serialization known for this rule type")
	}
}

func (strSerialize) Deserialize(string) (Predicate, error) {
	// TODO
	return nil, nil
}
