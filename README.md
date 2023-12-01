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
cd your-go-project-path
go get -u github.com/enzo-santos/imgurcrawler
```

Import it into your code as

```go
import (
    "github.com/enzo-santos/imgurcrawler"
)
```

Sample usage:

```go
img, err := imgurcrawler.GetImage("L1PQAPa")
// This executes a network request to https://i.imgur.com
// `img` is an instance of `imgurcrawler.ImgurImage`

if err != nil {
    panic(err)
}

fmt.Println(img.Id)     // = "L1PQAPa"
fmt.Println(img.Exists) // = (if it exists in the Imgur website: true or false)

fmt.Println(img.Filename)  // = "L1PQAPa.png" (always ends with `.png`)
fmt.Println(img.Extension) // = (the true extension of this image file: "jpg", "png"...)
fmt.Println(img.Content)   // = (the contents of this image file, as []byte)

fmt.Println(img.Name()) // = (the true name of this file: '<img.Id>.<img.Extension>')

err := imgurcrawler.SaveImage(img, "/home/username/Downloads")
if err != nil {
    panic(err)
} 
// At this point of execution, the image is downloaded to the provided directory

randomId := imgurcrawler.RandomId() 
// A random Imgur ID that can be provided to `imgurcrawler.ImgurImage`
// This does not necessarily refer to an existing image in the Imgur website
```

See the source code for better documentation.


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
