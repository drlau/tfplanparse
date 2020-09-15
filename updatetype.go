package tfplanparse

type UpdateType string

const (
	NoOpResource          UpdateType = "no-op"
	NewResource           UpdateType = "created"
	UpdateInPlaceResource UpdateType = "updateInPlace"
	ForceReplaceResource  UpdateType = "forceReplace"
	DestroyResource       UpdateType = "destroyed"
	ReadResource          UpdateType = "read"
)
