package cfg

import "github.com/jinzhu/configor"

var _config = new(Cfg)

type Cfg struct {
	CacheCapacity int        `json:"cacheCapacity"`
	ThumbnailMode bool       `json:"thumbnailMode"`
	Logger        LoggerConf `json:"logger"`
	HttpServer    HttpServer `json:"httpServer"`
}

type HttpServer struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type LoggerConf struct {
	Level string `json:"level"`
}

func InitConfig(configFile string) error {
	err := configor.Load(_config, configFile)
	if err != nil {
		return err
	}

	return nil
}

func Config() *Cfg {
	return _config
}
