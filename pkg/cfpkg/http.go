package cfpkg

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/kramerul/shalm/pkg/shalm"
)

var httpClient = &http.Client{Timeout: time.Second * 60}
var get func(url string) (io.ReadCloser, error)

func init() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	httpClient = &http.Client{
		Timeout: time.Second * 60,
	}
	dirCache := shalm.NewDirCache(path.Join(homedir, ".shalm", "cf"))
	get = dirCache.WrapReader(getWithEtag)
}

func getWithEtag(url string, etagOld string) (io.ReadCloser, string, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Add("If-None-Match", etagOld)
	res, err := httpClient.Do(request)
	if err != nil {
		return nil, "", fmt.Errorf("Error fetching %s: %v", url, err)
	}
	if res.StatusCode == 304 {
		return nil, etagOld, nil
	}
	if res.StatusCode != 200 {
		return nil, "", fmt.Errorf("Error fetching %s: status=%d", url, res.StatusCode)
	}
	etag := res.Header.Get("Etag")
	if len(etag) == 0 {
		etag = fmt.Sprintf("%x", time.Now().Unix())
	}
	return res.Body, etag, nil
}
