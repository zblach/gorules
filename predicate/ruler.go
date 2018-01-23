package predicate

import (
	"strings"

	. "./errors"
	. "./types"
	"reflect"
)

// Basic functions
var exists = MakeBasePredicate(EXISTS, func(obj Object) (bool, error) { return obj == nil, nil })

func Exists() Predicate {
	return exists
}

func Equals(value interface{}) Predicate {
	switch val := value.(type) {
	case string:
		return MakeStringPredicate(
			EQUALS,
			func(obj Object) (bool, error) {
				if v, ok := obj.(string); ok {
					return v == val, nil
				} else {
					return false, TypeError
				}
			},
			val)
	case float64:
		return MakeNumericPredicate(
			EQUALS,
			func(obj Object) (bool, error) {
				if v, ok := obj.(float64); ok {
					return v == val, nil
				} else {
					return false, TypeError
				}
			},
			val)
	default:
		typ := reflect.TypeOf(val)
		return MakeBasePredicate(
			EQUALS,
			func(obj Object) (bool, error) {
				if reflect.TypeOf(obj) == typ {
					return value == obj, nil
				} else {
					return false, TypeError
				}
			},
		)
	}
}

// String functions

func BeginsWith(prefix string) Predicate {
	return MakeStringPredicate(
		BEGINS_WITH,
		func(obj Object) (bool, error) {
			if v, ok := obj.(string); ok {
				return strings.HasPrefix(v, prefix), nil
			} else {
				return false, TypeError
			}
		},
		prefix)
}

func EndsWith(suffix string) Predicate {
	return MakeStringPredicate(
		ENDS_WITH,
		func(obj Object) (bool, error) {
			if v, ok := obj.(string); ok {
				return strings.HasSuffix(v, suffix), nil
			} else {
				return false, TypeError
			}
		},
		suffix)
}

func Contains(substring string) Predicate {
	return MakeStringPredicate(
		CONTAINS,
		func(obj Object) (bool, error) {
			if v, ok := obj.(string); ok {
				return strings.Contains(v, substring), nil
			} else {
				return false, TypeError
			}
		},
		substring)
}

// Numeric

func LessThan(value float64) Predicate {
	return MakeNumericPredicate(
		LESS_THAN,
		func(obj Object) (bool, error) {
			if v, ok := obj.(float64); ok {
				return v < value, nil
			} else {
				return false, TypeError
			}
		},
		value)
}

func GreaterThan(value float64) Predicate {
	return Not(Or(LessThan(value), Equals(value)))
}

func NearEqual(value, epsilon float64) Predicate {
	return Not(Or(GreaterThan(value+epsilon), LessThan(value-epsilon)))
}

// Predicate Composition

func Not(p Predicate) Predicate {
	if p.Type() == NOT {
		np := p.(CompositePredicate)
		return np.Values()[0] // :)
	}

	return MakeCompositePredicate(
		NOT,
		func(obj Object) (bool, error) {
			if b, err := p.Match(obj); err == nil {
				return !b, nil
			} else {
				return false, err
			}
		},
		[]Predicate{p})
}

func And(p1 Predicate, p2 Predicate, ps ...Predicate) Predicate {
	allp := make([]Predicate, len(ps)+2)
	copy(allp, append([]Predicate{p1, p2}, ps...))
	ps = allp

	return MakeCompositePredicate(
		AND,
		func(obj Object) (bool, error) {
			errs := CompositeError{}
			for _, r := range ps {
				b, err := r.Match(obj)
				if err != nil {
					errs.Append(err)
				}
				if !b {
					return false, errs.NilZero()
				}
			}
			return true, errs.NilZero()
		},
		ps)
}

func Or(p1 Predicate, p2 Predicate, ps ...Predicate) Predicate {
	allp := make([]Predicate, len(ps)+2)
	copy(allp, append([]Predicate{p1, p2}, ps...))
	ps = allp

	return MakeCompositePredicate(
		OR,
		func(obj Object) (bool, error) {
			errs := CompositeError{}
			for _, r := range ps {
				b, err := r.Match(obj)
				if err != nil {
					errs.Append(err)
				}
				if b {
					return true, errs.NilZero()
				}
			}
			return false, errs.NilZero()
		},
		ps)
}
