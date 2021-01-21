package mconfig

import (
	"github.com/spf13/viper"
)

// sourceType 配置来源类型
type sourceType int

// OnChange 配置变化处理方法
type OnChange func(v *viper.Viper) error

// Provier .
type Provier interface {
	Watch(OnChange) error
	Viper() *viper.Viper
}

var mgr Provier

// provierOptions .
type provierOptions struct {
	appID string
}

// ProvierOption .
type ProvierOption interface {
	apply(*provierOptions)
}

type optionFunc func(*provierOptions)

func (fn optionFunc) apply(o *provierOptions) {
	fn(o)
}

// WithAppID 设置 APP ID
func WithAppID(appID string) ProvierOption {
	return optionFunc(func(options *provierOptions) {
		options.appID = appID
	})
}

// Init .
func Init(opts ...ProvierOption) (err error) {
	options := &provierOptions{}

	for _, opt := range opts {
		opt.apply(options)
	}

	// 有指定配置文件,从本地文件读取配置
	if configFile != "" {
		mgr, err = NewFileProvider()
	} else { // 从阿里云 acm 读取配置
		mgr, err = NewAcmProvier(options)
	}
	if err == nil {
		viper.MergeConfigMap(mgr.Viper().AllSettings())
	}
	return
}

// Viper .
func Viper() *viper.Viper {
	return mgr.Viper()
}

// Watch .
func Watch(f OnChange) error {
	return mgr.Watch(f)
}
