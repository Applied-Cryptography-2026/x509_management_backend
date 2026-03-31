package repository

// DBRepository abstracts transaction semantics at the use-case layer.
// This allows interactors to run multiple repository calls inside a single DB transaction.
// Returns interface{} to stay generic — callers type-assert the result.
type DBRepository interface {
	Transaction(func(interface{}) (interface{}, error)) (interface{}, error)
}
