package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"github.com/ncw/swift"
	"io"
	"os"
	"path/filepath"
)

var rs *swift.Connection

var rsUsername = flag.String("rsuser", "", "Rackspace username")
var rsApiKey = flag.String("rskey", "", "Rackspace API key")
var rsBucket = flag.String("bucket", "", "Rackspace Cloud Files bucket")
var rsAssetUrl = flag.String("asseturl", "", "Asset URL")

func PutCloudFile(path string, r io.Reader) error {
	_, err := rs.ObjectPut(*rsBucket, path, r, false, "", "text/html", swift.Headers{"Cache-Control": "public, max-age=300"})
	return err
}

func AssetPath(a string) string {
	if *upload {
		return *rsAssetUrl + "/" + assets[a]
	}

	return "/assets/" + assets[a]
}

var assets = map[string]string{"bootstrap.min.css": "", "style.css": ""}
var contentTypes = map[string]string{".css": "text/css", ".js": "text/javascript"}

func writeAssets() {
	for a := range assets {
		f, err := os.Open("assets/" + a)
		MaybePanic(err)
		defer f.Close()

		h := sha256.New()
		_, err = io.Copy(h, f)
		MaybePanic(err)
		hash := h.Sum(nil)
		f.Seek(0, 0)

		ext := filepath.Ext(a)
		name := a[:len(a)-len(ext)] + "-" + hex.EncodeToString(hash[:4]) + ext

		assets[a] = name
		if *upload {
			_, err = rs.ObjectPut(*rsBucket+"-assets", name, f, false, "", contentTypes[ext], swift.Headers{"Cache-Control": "public, max-age=31556925"})
		} else {
			err = os.MkdirAll("output/assets", os.ModeDir|os.ModePerm)
			MaybePanic(err)
			err = os.MkdirAll("output/firms", os.ModeDir|os.ModePerm)
			MaybePanic(err)

			w, err := os.Create("output/assets/" + name)
			MaybePanic(err)
			_, err = io.Copy(w, f)
			w.Close()
		}
		MaybePanic(err)
	}
}

func init() {
	flag.Parse()

	if *upload {
		rs = &swift.Connection{UserName: *rsUsername, ApiKey: *rsApiKey, AuthUrl: "https://identity.api.rackspacecloud.com/v1.0"}
		err := rs.Authenticate()
		MaybePanic(err)
	}
}