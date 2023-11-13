package iterating

type ListStringIterator struct {
	Values []string
	index  int
}

func (self *ListStringIterator) HasNext() bool {
	return self.index < len(self.Values)
}

func (self *ListStringIterator) Next() string {
	value := self.Values[self.index]
	self.index += 1
	return value
}

func (self *ListStringIterator) Close() {}
