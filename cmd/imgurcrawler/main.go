package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/enzo-santos/imgurcrawler"
	"github.com/enzo-santos/imgurcrawler/internal/iterating"
	"github.com/gen2brain/beeep"
	"os"
	"time"
)

func main() {
	parser := argparse.NewParser("imgurcrawler", "An image crawler that collects random images from Imgur")
	pDelay := parser.Int("d", "delay", &argparse.Options{Help: "Delay between tries, in seconds", Default: 1})
	pStdinArgs := parser.StringList("i", "input", &argparse.Options{Help: "Input as strings"})
	pInputFilePaths := parser.FileList("f", "file", os.O_RDONLY, 0444, &argparse.Options{Help: "Input as files"})
	pShouldNotNotify := parser.Flag("", "no-notify", &argparse.Options{Help: "Do not launch OS-notification on hit"})
	pShouldNotStdout := parser.Flag("", "no-stdout", &argparse.Options{Help: "Do not print to standard output"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}
	iterators := make([]iterating.StringIterator, 0)
	stdinArgs := *pStdinArgs
	if len(stdinArgs) > 0 {
		iterators = append(iterators, &iterating.ListStringIterator{Values: stdinArgs})
	}
	for _, file := range *pInputFilePaths {
		iterators = append(iterators, &iterating.FileStringIterator{Path: file.Name()})
	}
	if len(iterators) == 0 {
		iterators = append(iterators, &iterating.RandomStringIterator{})
	}
	iterator := &iterating.CombinerStringIterator{Iterators: iterators}
	defer iterator.Close()

	shouldNotify := !(*pShouldNotNotify)
	shouldPrint := !(*pShouldNotStdout)
	delay := time.Duration(*pDelay)

	var count int
	for iterator.HasNext() {
		id := iterator.Next()
		if shouldPrint {
			fmt.Printf("%s=", id)
		}
		hit, err := imgurcrawler.DownloadImage(id, "build/images")
		if err != nil {
			panic(err)
		}

		if hit {
			count += 1
			if shouldPrint {
				fmt.Println("hit")
			}
			if shouldNotify {
				err := beeep.Notify("New image found", id, "assets/information.png")
				if err != nil {
					panic(err)
				}
			}
		} else {
			if shouldPrint {
				fmt.Printf("miss (%d)\n", count)
			}
		}
		if delay > 0 {
			time.Sleep(delay * time.Second)
		}
	}
}
