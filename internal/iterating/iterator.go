package iterating

type StringIterator interface {
	HasNext() bool
	Next() string
	Close()
}
