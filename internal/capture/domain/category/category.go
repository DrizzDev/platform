package category

// Category is the fixed set of sensitive product-data classes; anything unclassified fails closed (SEC-011).
type Category string

const (
	Prompt   Category = "PROMPT"
	Response Category = "RESPONSE"

	Tool      Category = "TOOL"
	Screen    Category = "SCREEN"
	Hierarchy Category = "HIERARCHY"

	Log  Category = "LOG"
	File Category = "FILE"
)

func (category Category) Valid() bool {
	switch category {
	case Prompt, Response, Tool, Screen, Hierarchy, File, Log:
		return true
	default:
		return false
	}
}

func (category Category) String() string {
	return string(category)
}
