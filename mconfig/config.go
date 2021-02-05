package mconfig

import (
	"github.com/spf13/viper"
)

// sourceType 配置来源类型
type sourceType int

// OnChange 配置变化处理方法
type OnChange func(v *viper.Viper) error

// Provider .
type Provider interface {
	Watch(OnChange) error
	Viper() *viper.Viper
}

var mgr Provider

// providerOptions .
type providerOptions struct {
	appID string
}

// ProviderOption .
type ProviderOption interface {
	apply(*providerOptions)
}

type optionFunc func(*providerOptions)

func (fn optionFunc) apply(o *providerOptions) {
	fn(o)
}

// WithAppID 设置 APP ID
func WithAppID(appID string) ProviderOption {
	return optionFunc(func(options *providerOptions) {
		options.appID = appID
	})
}

// Init .
func Init(opts ...ProviderOption) (err error) {
	options := &providerOptions{}

	for _, opt := range opts {
		opt.apply(options)
	}

	// 有指定配置文件,从本地文件读取配置
	if configFile != "" {
		mgr, err = NewFileProvider()
	} else { // 从阿里云 acm 读取配置
		mgr, err = NewAcmProvider(options)
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
