package srv

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/ability-sh/abi-lib/dynamic"
	"github.com/ability-sh/abi-micro/micro"
	"github.com/google/uuid"
)

const (
	SERVICE_CONFIG = "abi-app-user"
)

type ConfigService struct {
	name   string
	config interface{}

	Secret     string `json:"secret"`
	Db         string `json:"db"`
	Collection string `json:"collection"`
}

func newConfigService(name string, config interface{}) *ConfigService {
	return &ConfigService{name: name, config: config}
}

/**
* 服务名称
**/
func (s *ConfigService) Name() string {
	return s.name
}

/**
* 服务配置
**/
func (s *ConfigService) Config() interface{} {
	return s.config
}

/**
* 初始化服务
**/
func (s *ConfigService) OnInit(ctx micro.Context) error {

	dynamic.SetValue(s, s.config)

	return nil
}

/**
* 校验服务是否可用
**/
func (s *ConfigService) OnValid(ctx micro.Context) error {
	return nil
}

func (s *ConfigService) Recycle() {

}

func (s *ConfigService) NewID(ctx micro.Context) string {
	return strconv.FormatInt(ctx.Runtime().NewID(), 36)
}

func (s *ConfigService) SecPassword(p string) string {
	m := md5.New()
	m.Write([]byte(p))
	m.Write([]byte(s.Secret))
	return hex.EncodeToString(m.Sum(nil))
}

func (s *ConfigService) NewPassword() string {
	return uuid.New().String()
}

func GetConfigService(ctx micro.Context, name string) (*ConfigService, error) {
	s, err := ctx.GetService(name)
	if err != nil {
		return nil, err
	}
	ss, ok := s.(*ConfigService)
	if ok {
		return ss, nil
	}
	return nil, fmt.Errorf("service %s not instanceof *ConfigService", name)
}

func init() {
	micro.Reg(SERVICE_CONFIG, func(name string, config interface{}) (micro.Service, error) {
		return newConfigService(name, config), nil
	})
}
