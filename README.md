# imgurcrawler

An image crawler that collects random images from Imgur.

This script generates a fixed string of 7 characters composed of `A-Za-z0-9`
and tries to download an Imgur image whose ID is equal to that string. If it
finds something, it downloads to the *build/images* folder.

Since this process is "CPU bound", the script can also send a OS-notification
whenever it finds an image. It also detects if the downloaded image is
a *.gif*, a *.png* or a *.jpg*, renaming the file accordingly.


## Installing

Clone this repository:

```shell
git clone https://github.com/enzo-santos/imgurcrawler.git
cd imgurcrawler
```


## Usage

The basic usage is

```shell
go run main.go
```

It tries to find random Imgur images *ad infinitum*.


### Creational flags

All creational flags can be used with each other.

```shell
go run main.go -i "FQHZlGg" -i "olTru93"
```

It tries to find Imgur images for each `-i` argument given. Can be used one or
more times.

```shell
go run main.go -f args.txt
```

It tries to find Imgur images for each line of the `-f` argument given, as a
file. Can be used one or more times.


### Behavioral flags

```shell
go run main.go --no-notify
```

If it finds any Imgur image, it'll not display the OS-notification.

```shell
go run main.go --no-stdout
```

It'll not print anything to the standard output (*stdout*).

```shell
go run main.go --delay 2
```

It waits for the given amount of seconds before trying to find another Imgur
image.
