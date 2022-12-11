package resource

import (
	"flag"
	"github.com/SnDragon/go-treasure-chest/internal/app/shorturl/config"
	"github.com/SnDragon/go-treasure-chest/internal/app/shorturl/storage"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

var (
	ConfFile = flag.String("conf", "configs/shorturl/app.yaml", "配置文件路径")
	Storage  storage.Storage
)

type InitFunc func() error

func InitResource() error {
	fns := []InitFunc{
		InitAppConfig,
		InitRedisStorage,
	}
	for _, fn := range fns {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}

func InitAppConfig() error {
	flag.Parse()
	log.Printf("confFile: %v", *ConfFile)
	data, err := ioutil.ReadFile(*ConfFile)
	if err != nil {
		return errors.Wrapf(err, "[InitAppConfig] read config file:%v err", ConfFile)
	}
	// 加载配置文件
	if err := yaml.Unmarshal(data, config.AppConfig); err != nil {
		return errors.Wrap(err, "[InitAppConfig] yaml.Unmarshal err")
	}
	log.Printf("app config init succeed, conf: %+v", config.AppConfig)
	return nil
}

func InitRedisStorage() error {
	var err error
	Storage, err = storage.NewRedisStorage()
	if err != nil {
		return errors.Wrap(err, "InitRedisStorage err")
	}
	return nil
}
