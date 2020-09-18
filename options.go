package tfplanparse

type GetBeforeAfterOptions func(a attributeChange) bool

func IgnoreComputed(a attributeChange) bool {
	return a.IsComputed()
}

func IgnoreSensitive(a attributeChange) bool {
	return a.IsSensitive()
}

func IgnoreNoOp(a attributeChange) bool {
	return a.IsNoOp()
}

func ComputedOnly(a attributeChange) bool {
	return !a.IsComputed()
}
