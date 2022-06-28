package initialize

import (
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/olivere/elastic/v7"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"github.com/xiaobeicn/go-es-location/global"
	"net"
	"net/http"
	"time"
)

func Init() {
	InitViperConfig()
	InitES()
}

func InitViperConfig() {
	var configFile string
	flag.StringVar(&configFile, "c", "./config.yaml", "指定配置文件")
	if len(configFile) == 0 {
		panic("配置文件不存在！")
	}
	v := viper.New()
	v.SetConfigFile(configFile)
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("配置解析失败:%s\n", err))
	}
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件发生改变")
		if err := v.Unmarshal(&global.GConfig); err != nil {
			panic(fmt.Errorf("配置重载失败:%s\n", err))
		}
	})
	if err := v.Unmarshal(&global.GConfig); err != nil {
		panic(fmt.Errorf("配置重载失败:%s\n", err))
	}
}

func InitES() {
	elasticConfig := global.GConfig.Elastic
	client, err := elastic.NewClient(
		elastic.SetHttpClient(&http.Client{
			Transport: &http.Transport{
				DialContext:         (&net.Dialer{Timeout: cast.ToDuration("1s")}).DialContext,
				MaxIdleConns:        50,
				MaxIdleConnsPerHost: 10,
			},
			Timeout: cast.ToDuration("1s"),
		}),
		elastic.SetURL(elasticConfig.Url),
		elastic.SetSniff(elasticConfig.Sniff),
		elastic.SetGzip(true),
		elastic.SetDecoder(&elastic.NumberDecoder{}), // critical to ensure decode of int64 won't lose precise
		elastic.SetRetrier(elastic.NewBackoffRetrier(elastic.NewExponentialBackoff(128*time.Millisecond, 513*time.Millisecond))),
	)
	if err != nil {
		panic("创建ES客户端错误:" + err.Error())
	}
	global.GElastic = client
}
