package ossync_oss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/pubgo/golugin/client/golug_oss"
)

var client *oss.Bucket

func GetBucket() *oss.Bucket { return client }
func InitBucket(name string) { client = golug_oss.GetClient(name) }
