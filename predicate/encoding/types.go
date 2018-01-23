package encoding

import . "../types"

type Serializable interface {
	Serialize(Predicate) (string, error)
	Deserialize(string) (Predicate, error)
}
