package global

import (
	"github.com/olivere/elastic/v7"
	"github.com/xiaobeicn/go-es-location/config"
)

var (
	GConfig  config.ServerConfig // 全局配置
	GElastic *elastic.Client
)
