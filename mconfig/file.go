package mconfig

import (
	"flag"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var configFile string

func init() {
	addFileFlag(flag.CommandLine)
}

func addFileFlag(fs *flag.FlagSet) {
	fs.StringVar(&configFile, "config", "", "config file.")
}

type file struct {
	v *viper.Viper
}

// NewFileProvider 从文件读取配置
func NewFileProvider() (c Provier, err error) {
	v := viper.New()
	v.SetConfigType("toml")
	v.SetConfigFile(configFile)
	if err = v.ReadInConfig(); err != nil {
		return
	}
	f := &file{
		v: v,
	}
	c = f
	return
}

// WatchConfig 监测配置变化
func (f *file) Watch(onChange OnChange) (err error) {
	f.v.OnConfigChange(func(e fsnotify.Event) {
		switch e.Op {
		case fsnotify.Write, fsnotify.Create:
			onChange(f.v)
		}
	})
	go f.v.WatchConfig()
	return
}

func (f *file) Viper() *viper.Viper {
	return f.v
}
