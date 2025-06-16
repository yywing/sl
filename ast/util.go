package ast

func GetDeterministicType(ts ...ValueType) ValueType {
	for _, t := range ts {
		if t.Kind() != TypeKindAny {
			return t
		}
	}
	return AnyType
}

func TypeEquals(t1, t2 ValueType) bool {
	return t1.Equals(t2)
}
