package storage

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/imroc/cos-go-sdk-v5/cos"
)

type TencentCOSBackend struct {
	Bucket *cos.Bucket
}

func NewTencentCOSBackend(url string) *TencentCOSBackend {
	secretId := os.Getenv("TENCENT_CLOUD_SECRET_ID")
	secretKey := os.Getenv("TENCENT_CLOUD_SECRET_KEY")
	if len(secretId) == 0 {
		panic("TENCENT_CLOUD_SECRET_ID environment variable is not set")
	}

	if len(secretKey) == 0 {
		panic("TENCENT_CLOUD_SECRET_KEY environment variable is not set")
	}

	client := cos.NewClient(secretId, secretKey)
	bucket := cos.NewBucketFromURL(url, client)
	return &TencentCOSBackend{
		Bucket: bucket,
	}
}

func (b TencentCOSBackend) ListObjects(prefix string) ([]Object, error) {
	var objects []Object
	cosPrefix := cos.Prefix(prefix)
	marker := cos.Marker("")
	for {
		result, err := b.Bucket.ListObjects(cosPrefix, marker)
		if err != nil {
			return nil, err
		}

		for _, obj := range result.Objects {
			path := removePrefixFromObjectPath(prefix, obj.Key)
			if objectPathIsInvalid(path) {
				continue
			}
			object := Object{
				Path:         path,
				Content:      []byte{},
				LastModified: obj.LastModified,
			}
			objects = append(objects, object)
		}
		if !result.IsTruncated {
			break
		}
		cosPrefix = cos.Prefix(result.Prefix)
		marker = cos.Marker(result.NextMarker)
	}
	return objects, nil
}

// PutObject uploads an object to Tencent Cloud COS bucket
func (b TencentCOSBackend) PutObject(path string, content []byte) error {
	return b.Bucket.PutObject(path, content)
}

// DeleteObject removes an object from Tencent Cloud COS bucket
func (b TencentCOSBackend) DeleteObject(path string) error {
	return b.Bucket.DeleteObject(path)
}

// GetObject retrieves an object from Tencent Cloud COS bucket
func (b TencentCOSBackend) GetObject(path string) (Object, error) {
	var object Object
	object.Path = path
	var content []byte
	body, err := b.Bucket.GetObject(path)
	if err != nil {
		return object, err
	}
	content, err = ioutil.ReadAll(body)
	body.Close()
	if err != nil {
		return object, err
	}
	object.Content = content

	headers, err := b.Bucket.GetObjectMeta(path)
	if err != nil {
		return object, err
	}
	lastModified, _ := http.ParseTime(headers.Get(cos.HTTPHeaderLastModified))
	object.LastModified = lastModified
	return object, nil
}
