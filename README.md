# imgurcrawler

An image crawler that collects random images from Imgur.

This script generates a fixed string of 7 characters composed of `A-Za-z0-9`
and tries to download an Imgur image whose ID is equal to that string. If it
finds something, it downloads to the *build/images* folder.

Since this process is "CPU bound", the script can also send a OS-notification
whenever it finds an image. It also detects if the downloaded image is
a *.gif*, a *.png* or a *.jpg*, renaming the file accordingly.


## Usage

Run the following in your command prompt:

```shell
go get -u github.com/enzo-santos/imgurcrawler
```

Import it into your code as

```go
import (
	"github.com/enzo-santos/imgurcrawler"
)
```

This package imports a single function named `DownloadImage`. Its receives an
Imgur ID as its first parameter and the directory to where the file will be
downloaded as its second parameter.

The function will then try to download the Imgur image with the given ID. If
there is no image, it returns *false*. Otherwise, it downloads the image with
the appropriate extension to the given directory and returns *true*:

```go
ok := imgurcrawler.DownloadImage("L1PQAPa", "build/images")
if ok {
	// Downloads https://i.imgur.com/L1PQAPa.jpeg to the build/images directory
} else {
	// Does nothing
}
```


## Command-line usage

To access the executable of this package, run the following:

```shell
go install github.com/enzo-santos/imgurcrawler/cmd/imgurcrawler@latest
```

It'll download it to the *bin* folder of the path shown by `go env GOPATH`.

The basic usage is

```shell
imgurcrawler
```

It tries to find random Imgur images *ad infinitum*.


### Creational flags

All creational flags can be used with each other.

```shell
imgurcrawler -i "FQHZlGg" -i "olTru93"
```

It tries to find Imgur images for each `-i` argument given. Can be used one or
more times.

```shell
imgurcrawler -f args.txt
```

It tries to find Imgur images for each line of the `-f` argument given, as a
file. Can be used one or more times.


### Behavioral flags

```shell
imgurcrawler --no-notify
```

If it finds any Imgur image, it'll not display the OS-notification.

```shell
imgurcrawler --no-stdout
```

It'll not print anything to the standard output (*stdout*).

```shell
imgurcrawler --delay 2
```

It waits for the given amount of seconds before trying to find another Imgur
image. The default delay is 1 second.
