package mconfig

import (
	"bytes"
	"errors"
	"flag"
	"os"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
)

var (
	namespaceID string
	accessKey   string
	secretKey   string
	endpoint    = "acm.aliyun.com"
	group       string
)

func init() {
	addAcmFlag(flag.CommandLine)
}

func addAcmFlag(fs *flag.FlagSet) {
	var (
		defNamespaceID = os.Getenv("ACM_NAMESPACE")
		defAccessKey   = os.Getenv("ACM_ACCESSKEY")
		defSecretKey   = os.Getenv("ACM_SECRETKEY")
		defGroup       = os.Getenv("ACM_GROUP")
	)
	fs.StringVar(&namespaceID, "acm.namespace", defNamespaceID, "acm namespace id.")
	fs.StringVar(&accessKey, "acm.accessKey", defAccessKey, "acm access key.")
	fs.StringVar(&secretKey, "acm.secretKey", defSecretKey, "acm secret key.")
	fs.StringVar(&group, "acm.group", defGroup, "acm group.")
}

type acm struct {
	options *providerOptions
	iClient config_client.IConfigClient
	v       *viper.Viper
}

type acmConfig struct {
	Namespace string `json:"namespace" form:"namespace"`
	AccessKey string `json:"access_key" form:"access_key"`
	SecretKey string `json:"secret_key" form:"secret_key"`
	Endpoint  string `json:"endpoint" form:"endpoint"`
	Group     string `json:"group" form:"group"`
}

func buildConfig() (c *acmConfig, err error) {
	if namespaceID == "" {
		err = errors.New("invalid acm namespace")
		return
	}

	if accessKey == "" {
		err = errors.New("invalid acm accessKey")
		return
	}

	if secretKey == "" {
		err = errors.New("invalid acm secretKey")
		return
	}

	if endpoint == "" {
		err = errors.New("invalid acm endpoint")
		return
	}

	if group == "" {
		err = errors.New("invalid acm group")
		return
	}

	return &acmConfig{
		Namespace: namespaceID,
		AccessKey: accessKey,
		SecretKey: secretKey,
		Endpoint:  endpoint,
		Group:     group,
	}, nil
}

// NewAcmProvider .
func NewAcmProvider(options *providerOptions) (c Provider, err error) {
	if options.appID == "" {
		err = errors.New("invalid app id")
		return
	}
	var acmConfig *acmConfig
	acmConfig, err = buildConfig()
	if err != nil {
		return
	}

	v := viper.New()
	v.SetConfigType("toml")
	a := &acm{
		options: options,
		v:       v,
	}
	clientConfig := constant.ClientConfig{
		Endpoint:       acmConfig.Endpoint + ":8080",
		NamespaceId:    acmConfig.Namespace,
		AccessKey:      acmConfig.AccessKey,
		SecretKey:      acmConfig.SecretKey,
		TimeoutMs:      5 * 1000,
		ListenInterval: 30 * 1000,
	}
	// Initialize client.
	a.iClient, err = clients.CreateConfigClient(map[string]interface{}{
		"clientConfig": clientConfig,
	})
	if err != nil {
		return
	}
	content, err := a.iClient.GetConfig(vo.ConfigParam{
		DataId: a.options.appID,
		Group:  acmConfig.Group,
	})
	a.v.ReadConfig(bytes.NewBufferString(content))
	c = a
	return
}

// Watch 监测配置变化
func (a *acm) Watch(onChange OnChange) (err error) {
	params := vo.ConfigParam{
		DataId: a.options.appID,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {
			a.v.ReadConfig(bytes.NewBufferString(data))
			onChange(a.v)
		},
	}
	a.iClient.ListenConfig(params)
	return
}

func (a *acm) Viper() *viper.Viper {
	return a.v
}
