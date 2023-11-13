package iterating

import (
	"bufio"
	"log"
	"os"
)

type FileStringIterator struct {
	Path    string
	file    *os.File
	scanner *bufio.Scanner
}

func (self *FileStringIterator) HasNext() bool {
	scanner := self.scanner
	if scanner == nil {
		file, err := os.Open(self.Path)
		if err != nil {
			log.Println(err)
			return false
		}
		self.file = file
		scanner = bufio.NewScanner(file)
		self.scanner = scanner
	}
	return scanner.Scan()
}

func (self *FileStringIterator) Next() string {
	return self.scanner.Text()
}

func (self *FileStringIterator) Close() {
	file := self.file
	if file != nil {
		file.Close()
	}
}
