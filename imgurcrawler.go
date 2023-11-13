package imgurcrawler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

var magicNumbers = map[[4]byte]string{
	{0xff, 0xd8, 0xff, 0xe0}: "jpg",
	{0xff, 0xd8, 0xff, 0xe1}: "jpg",
	{0xff, 0xd8, 0xff, 0xe2}: "jpg",
	{0xff, 0xd8, 0xff, 0xe8}: "jpg",
	{0xff, 0xd8, 0xff, 0xdb}: "jpg",
	{0x47, 0x49, 0x46, 0x38}: "gif",
	{0x89, 0x50, 0x4e, 0x47}: "png",
}

func DownloadUrl(url string) bool {
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
