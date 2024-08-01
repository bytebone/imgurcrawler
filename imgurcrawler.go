package imgurcrawler

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "path"
    "strings"

    "github.com/enzo-santos/imgurcrawler/internal/iterating"
)

var magicNumbers = map[[4]byte]string{
    {0xff, 0xd8, 0xff, 0xe0}: "jpg",
    {0xff, 0xd8, 0xff, 0xe1}: "jpg",
    {0xff, 0xd8, 0xff, 0xe2}: "jpg",
    {0xff, 0xd8, 0xff, 0xe8}: "jpg",
    {0xff, 0xd8, 0xff, 0xee}: "jpg",
    {0xff, 0xd8, 0xff, 0xfe}: "jpg",
    {0xff, 0xd8, 0xff, 0xdb}: "jpg",
    {0x47, 0x49, 0x46, 0x38}: "gif",
    {0x89, 0x50, 0x4e, 0x47}: "png",
}

// ImgurImage contains information about an Imgur image.
type ImgurImage struct {
    // The ID of this image in the Imgur website.
    //
    // It's composed by 7 characters that satisfy the regular expression `[A-Za-z0-9]`.
    Id string

    // If this image exists in the Imgur website.
    //
    // A image exists in the Imgur website if trying to
    // access *https://i.imgur.com/<Id>.png* does not redirect
    // to *https://i.imgur.com/removed.png*.
    Exists bool

    // The filename of this image in the Imgur website.
    //
    // Its extension may not reflect the image contents. This field will probably
    // contain a *.png* extension, even if the image is a JPG file, for example. Users
    // should rely on [Extension] to a more precise value.
    Filename string

    // The extension of this file.
    //
    // This field contains a non-dotted extension (e.g. "jpg", "png", "gif") that
    // describe this file contents, calculated by its first four bytes.
    Extension string

    // The contents of this image.
    Content []byte
}

// Name contains a valid basename for this image.
//
// Since `ImgurImage.Filename` does not offer a precise file extension, this method
// replaces the imprecise extension with a more precise one:
//
//  img, err := imgurcrawler.GetImage("d5hU9pb")
//  if err != nil {
//      panic(err)
//  }
//  fmt.Println(img.Filename)   // d5hU9pb.png
//  fmt.Println(img.Name())     // d5hU9pb.jpeg (actually a JPEG image)
func (img ImgurImage) Name() string {
    fname := img.Filename
    fname = strings.TrimSuffix(fname, path.Ext(fname))
    fname = fmt.Sprintf("%s.%s", fname, img.Extension)
    return fname
}

// RandomId creates a random Imgur ID.
//
// This may or may not represent a valid Imgur image, since it's generated locally
// without accessing the Imgur website:
//
//  id := imgurcrawler.RandomId()    // "d5hU9pb", "8yt3zhM"...
//  img, err := imgurcrawler.GetImage(id)
func RandomId() string {
    iterator := &iterating.RandomStringIterator{}
    return iterator.Next()
}

// GetImage loads a image from the Imgur website.
//
// `id` should be a 7-character string that matches `[A-Za-z0-9]`.
//
//  img, err := imgurcrawler.GetImage("d5hU9pb")
//  if err != nil {
//      panic(err)
//  }
//  fmt.Printf("Loaded %s\n", img.Name())
func GetImage(id string) (img ImgurImage, rerr error) {
    client := &http.Client{}

    url := fmt.Sprintf("https://i.imgur.com/%s.png", id)
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        rerr = fmt.Errorf("Error while trying to create GET %s: %v", url, err)
        return
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
        rerr = fmt.Errorf("Error while trying to execute GET %s, %v", url, err)
        return
    }
    defer res.Body.Close()

    var image ImgurImage
    if res.Request.URL.Path == "/removed.png" {
        image = ImgurImage{
            Id:     id,
            Exists: false,
        }

    } else {
        bytes, err := io.ReadAll(res.Body)
        if err != nil {
            rerr = fmt.Errorf("Error while trying to read the request: %v", err)
            return
        }

        image = ImgurImage{
            Id:        id,
            Exists:    true,
            Filename:  res.Request.URL.Path[1:],
            Extension: magicNumbers[([4]byte)(bytes[:4])],
            Content:   bytes,
        }
    }
    return image, nil
}

// SaveImage stores the image in the current file system.
//
// If `dpath` does not exist, it will be created recursively:
//
//	id := imgurcrawler.RandomId()
//	img, err := imgurcrawler.GetImage(id)
//	if err != nil {
//	    panic(err)
//	}
//	imgurcrawler.SaveImage(img, "/path/to/save/directory")
func SaveImage(img ImgurImage, dpath *string) error {
    os.MkdirAll(*dpath, 0755)
    fpath := path.Join(*dpath, img.Name())
    return os.WriteFile(fpath, img.Content, 0644)
}
