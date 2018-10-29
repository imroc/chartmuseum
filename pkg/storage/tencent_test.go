package storage

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type TencentTestSuite struct {
	suite.Suite
	TencentCOSBackend       *TencentCOSBackend
	BrokenTencentCOSBackend *TencentCOSBackend
}

func (suite *TencentTestSuite) SetupSuite() {
	backend := NewTencentCOSBackend("")
	suite.BrokenTencentCOSBackend = backend

	url := os.Getenv("TEST_STORAGE_TENCENT_BUCKET_URL")
	backend = NewTencentCOSBackend(url)
	suite.TencentCOSBackend = backend

	data := []byte("some object")
	path := "deleteme.txt"

	err := suite.TencentCOSBackend.PutObject(path, data)
	suite.Nil(err, "no error putting deleteme.txt using TencentCOS backend")
}

func (suite *TencentTestSuite) TearDownSuite() {
	err := suite.TencentCOSBackend.DeleteObject("deleteme.txt")
	suite.Nil(err, "no error deleting deleteme.txt using TencentCOS backend")
}

func (suite *TencentTestSuite) TestListObjects() {
	_, err := suite.BrokenTencentCOSBackend.ListObjects("")
	suite.NotNil(err, "cannot list objects with bad bucket")

	_, err = suite.TencentCOSBackend.ListObjects("data/www/app/latest/tgit-web-files/uploads/user/avatar/2")
	suite.Nil(err, "can list objects with good bucket")
}

func (suite *TencentTestSuite) TestGetObject() {
	_, err := suite.BrokenTencentCOSBackend.GetObject("this-file-cannot-possibly-exist.tgz")
	suite.NotNil(err, "cannot get objects with bad bucket")

	obj, err := suite.TencentCOSBackend.GetObject("deleteme.txt")
	suite.Nil(err, "can get object")
	suite.Equal([]byte("some object"), obj.Content, "able to get object")
}

func (suite *TencentTestSuite) TestPutObject() {
	err := suite.BrokenTencentCOSBackend.PutObject("this-file-will-not-upload.txt", []byte{})
	suite.NotNil(err, "cannot put objects with bad bucket")
}

func TestTencentStorageTestSuite(t *testing.T) {
	// req.Debug = true
	if os.Getenv("TEST_CLOUD_STORAGE") == "1" &&
		os.Getenv("TEST_STORAGE_TENCENT_BUCKET_URL") != "" {
		suite.Run(t, new(TencentTestSuite))
	}
}
