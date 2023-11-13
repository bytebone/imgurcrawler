package main

import (
	"bufio"
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/gen2brain/beeep"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

const asciiLowercase = "abcdefghijklmnopqrstuvwxyz"
const asciiUppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const digits = "0123456789"

var magicNumbers = map[[4]byte]string{
	{0xff, 0xd8, 0xff, 0xe0}: "jpg",
	{0xff, 0xd8, 0xff, 0xe1}: "jpg",
	{0xff, 0xd8, 0xff, 0xe2}: "jpg",
	{0xff, 0xd8, 0xff, 0xe8}: "jpg",
	{0xff, 0xd8, 0xff, 0xdb}: "jpg",
	{0x47, 0x49, 0x46, 0x38}: "gif",
	{0x89, 0x50, 0x4e, 0x47}: "png",
}

type StringIterator interface {
	HasNext() bool
	Next() string
	Close()
}

type RandomStringIterator struct{}

func (self *RandomStringIterator) HasNext() bool {
	return true
}

func (self *RandomStringIterator) Next() string {
	characters := fmt.Sprintf("%s%s%s", asciiLowercase, asciiUppercase, digits)
	path := make([]byte, 7)
	for i := 0; i < 7; i++ {
		path[i] = characters[rand.Intn(len(characters))]
	}
	return string(path)
}

func (self *RandomStringIterator) Close() {}

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

func DoRequest(url string) bool {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("authority", "i.imgur.com")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("sec-ch-ua", `"Google Chrome";v="89", "Chromium";v="89", ";Not A Brand";v="99"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36")
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("dnt", "1")
	req.Header.Set("sec-fetch-site", "none")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("accept-language", "en-GB,en;q=0.9")

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.Request.URL.Path == "/removed.png" {
		return false
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	extension, ok := magicNumbers[([4]byte)(bytes[:4])]
	if !ok {
		panic(fmt.Errorf("Invalid magic number: %v", bytes[:4]))
	}

	dpath := path.Join("build", "images")
	os.MkdirAll(dpath, 755)

	fname := res.Request.URL.Path[1:]
	fname = strings.TrimSuffix(fname, path.Ext(fname))
	fname = fmt.Sprintf("%s.%s", fname, extension)
	fpath := path.Join(dpath, fname)

	err = os.WriteFile(fpath, bytes, 644)
	if err != nil {
		panic(err)
	}

	return true
}

func main() {
	parser := argparse.NewParser("imgurcrawler", "A image crawler that collects random images from Imgur")
	delay := parser.Int("d", "delay", &argparse.Options{Help: "Delay between tries, in seconds", Default: 1})
	stdinArgs := parser.StringList("i", "input", &argparse.Options{Help: "Input as strings"})
	inputFilePaths := parser.FileList("f", "file", os.O_RDONLY, 0444, &argparse.Options{Help: "Input as files"})
	shouldNotNotify := parser.Flag("", "no-notify", &argparse.Options{Help: "Do not launch OS-notification on hit"})
	shouldNotStdout := parser.Flag("", "no-stdout", &argparse.Options{Help: "Do not print to standart output"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}
	iterators := make([]StringIterator, 0)
	if len(*stdinArgs) > 0 {
		iterators = append(iterators, &ListStringIterator{Values: *stdinArgs})
	}
	for _, file := range *inputFilePaths {
		iterators = append(iterators, &FileStringIterator{Path: file.Name()})
	}
	if len(iterators) == 0 {
		iterators = append(iterators, &RandomStringIterator{})
	}
	iterator := &CombinerStringIterator{Iterators: iterators}
	defer iterator.Close()

	shouldNotify := !(*shouldNotNotify)
	shouldPrint := !(*shouldNotStdout)

	var count int
	for iterator.HasNext() {
		filename := iterator.Next()
		url := fmt.Sprintf("https://i.imgur.com/%s.png", filename)
		if shouldPrint {
			fmt.Printf("%32s=", url)
		}
		hit := DoRequest(url)
		if hit {
			count += 1
			if shouldPrint {
				fmt.Println("hit")
			}
			if shouldNotify {
				err := beeep.Notify("Nova imagem encontrada", url, "assets/information.png")
				if err != nil {
					panic(err)
				}
			}
		} else {
			if shouldPrint {
				fmt.Printf("miss (%d)\n", count)
			}
		}
		time.Sleep(time.Duration(*delay) * time.Second)
	}
}
