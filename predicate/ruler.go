package predicate

import (
	"strings"

	. "./types"
	. "./errors"
)

// Basic functions
var exists = MakeBasePredicate(EXISTS, func(obj Object) (bool, error) { return obj == nil, nil })

func Exists() Predicate {
	return exists
}

// String functions
func Equals(str string) Predicate {
	return MakeStringPredicate(
		EQUALS,
		func(obj Object) (bool, error) {
			if v, ok := obj.(string); ok {
				return v == str, nil
			} else {
				return false, TypeError
			}
		},
		str)
}

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

// Predicate Composition

func Not(p Predicate) Predicate {
	if p.Type() == NOT {
		np := p.(*CompositePredicate)
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
	ps = append([]Predicate{p1, p2}, ps...) // forces minimum size of 2

	return MakeCompositePredicate(
		AND,
		func(obj Object) (bool, error) {
			var errs *CompositeError
			for _, r := range ps {
				b, err := r.Match(obj)
				if err != nil {
					if errs == nil {
						errs = &CompositeError{}
					}
					errs.Append(err)
				}
				if !b {
					return false, errs
				}
			}
			return true, errs
		},
		ps)
}

func Or(p1 Predicate, p2 Predicate, ps ...Predicate) Predicate {
	ps = append([]Predicate{p1, p2}, ps...)

	return MakeCompositePredicate(
		OR,
		func(obj Object) (bool, error) {
			var errs *CompositeError
			for _, r := range ps {
				b, err := r.Match(obj)
				if err != nil {
					if errs == nil {
						errs = &CompositeError{}
					}
					errs.Append(err)
				}
				if b {
					return true, errs
				}
			}
			return false, errs
		},
		ps)
}
