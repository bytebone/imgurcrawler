package iterating

type CombinerStringIterator struct {
	Iterators []StringIterator
	index     int
}

func (self *CombinerStringIterator) HasNext() bool {
	if len(self.Iterators) == 0 {
		return false
	}
	iterator := self.Iterators[self.index]
	if iterator.HasNext() {
		return true
	}
	if self.index < len(self.Iterators)-1 {
		self.index += 1
		return self.HasNext()
	}
	return false
}

func (self *CombinerStringIterator) Next() string {
	return self.Iterators[self.index].Next()
}

func (self *CombinerStringIterator) Close() {
	for _, iterator := range self.Iterators {
		iterator.Close()
	}
}
