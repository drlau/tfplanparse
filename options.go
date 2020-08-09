package tfplanparse

type GetBeforeAfterOptions func(a *AttributeChange) bool

func IgnoreComputed(a *AttributeChange) bool {
	return a.IsComputed()
}

func IgnoreSensitive(a *AttributeChange) bool {
	return a.IsSensitive()
}
