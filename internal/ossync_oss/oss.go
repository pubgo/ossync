package ossync_oss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/pubgo/golugin/client/oss_cli"
)

func GetBucket(names ...string) *oss.Bucket { return oss_cli.GetClient(names...) }
