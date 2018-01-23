package types

type PredicateType int

type Object interface{}

const (
	UNDEFINED PredicateType = iota

	EXISTS
	EQUALS

	LESS_THAN

	BEGINS_WITH
	ENDS_WITH
	CONTAINS
	REGEX_MATCH

	AND
	OR
	NOT
)

type Predicate interface {
	Match(Object) (bool, error)
	Type() PredicateType
}

type Matcher func(Object) (bool, error)

type BasePredicate struct {
	ptype   PredicateType
	matcher Matcher
}

func MakeBasePredicate(predicateType PredicateType, matcher Matcher) BasePredicate {
	return BasePredicate{
		ptype:   predicateType,
		matcher: matcher,
	}
}

func (p BasePredicate) Match(obj Object) (bool, error) { return p.matcher(obj) }
func (p BasePredicate) Type() PredicateType            { return p.ptype }

type StringPredicate struct {
	BasePredicate
	param string
}

func (p StringPredicate) Value() string { return p.param }

func MakeStringPredicate(predicateType PredicateType, matcher Matcher, value string) StringPredicate {
	return StringPredicate{
		BasePredicate: MakeBasePredicate(predicateType, matcher),
		param:         value,
	}
}

type NumericPredicate struct {
	BasePredicate
	value float64
}

func MakeNumericPredicate(predicateType PredicateType, matcher Matcher, value float64) NumericPredicate {
	return NumericPredicate{
		BasePredicate: MakeBasePredicate(predicateType, matcher),
		value:         value,
	}
}

func (p NumericPredicate) Value() float64 { return p.value }

type CompositePredicate struct {
	BasePredicate
	components []Predicate
}

func (p CompositePredicate) Values() []Predicate { return p.components }

func MakeCompositePredicate(predicateType PredicateType, matcher Matcher, components []Predicate) CompositePredicate {
	return CompositePredicate{
		BasePredicate: MakeBasePredicate(predicateType, matcher),
		components:    components,
	}
}
