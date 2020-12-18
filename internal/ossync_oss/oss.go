package ossync_oss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/pubgo/golug/golug_consts"
	"github.com/pubgo/golug/plugins/golug_oss"
)

var client *oss.Bucket

func GetBucket() *oss.Bucket { return client }
func InitBucket(name string) {
	if name == "" {
		name = golug_consts.Default
	}
	client = golug_oss.GetClient(name)
}
