package boot

import (
	"flag"
	"github.com/ydx1011/gopher-core"
	"os"
	"sync"
)

const EnvNameConfigFile = "GOPHER_CONFIG_FILE"

var (
	// 默认的配置路径
	ConfigPath = "application.yaml"

	creator func() gopher.Application = defaultCreator
	gApp    gopher.Application
	once    sync.Once
)

func defaultCreator() gopher.Application {
	if conf, ok := os.LookupEnv(EnvNameConfigFile); ok {
		ConfigPath = conf
	}
	flag.StringVar(&ConfigPath, "f", ConfigPath, "Application configuration file path.")
	flag.Parse()
	return gopher.NewFileConfigApplication(ConfigPath)
}
