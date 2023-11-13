# imgurcrawler

A image crawler that collects random images from Imgur.

This script generates a fixed string of 7 characters composed of `A-Za-z0-9` and tries 
to download a Imgur image whose ID is equal to that string. If it finds something, it 
downloads to the *build/images* folder.

Since this process is "CPU bound", the script also sends a notification whenever it 
finds an image. It also detects if the downloaded image is a .gif, a .png or a .jpg, 
renaming the file accordingly.

Future features are command-line related:

- `-i`: inputs via command-line one or more 7-character string for the script to test
- `-f`: inputs a file where each line is a 7-character string for the script to test
- `--(no-)notify`: whether should the script send the notification when finds something


## Installing

Clone this repository:

```shell
git clone https://github.com/enzo-santos/imgurcrawler.git
cd imgurcrawler
```


## Usage

```shell
go run main.go
```
