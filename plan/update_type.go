package plan

type UpdateType string

const (
	NewResource           UpdateType = "created"
	UpdateInPlaceResource UpdateType = "updateInPlace"
	ForceReplaceResource  UpdateType = "forceReplace"
	DestroyResource       UpdateType = "destroyed"
	ReadResource          UpdateType = "read"
)
