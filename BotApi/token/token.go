package token

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

// 定义常量
const (
	TypeBot string = "Bot"
)

// 用于存储调用API所需的Token信息
type Token struct {
	AppID       uint64
	AccessToken string
	Type        string
}

// yaml 配置文件中的信息格式
type Config struct {
	AppID uint64 `yaml:"appid"`
	Token string `yaml:"token"`
}

// 创建一个默认的 token，默认是 bot 身份
func CreateDefaultToken() *Token {
	return &Token{
		Type: TypeBot,
	}
}

// 创建一个机器人身份的 token
func CreateBotToken(appID uint64, accessToken string) *Token {
	return &Token{
		AppID:       appID,
		AccessToken: accessToken,
		Type:        TypeBot,
	}
}

// 通过 token 获取 Authorization 字符串
func (t *Token) GetAuthorization() (res string) {
	switch t.Type {
	case TypeBot:
		res = fmt.Sprintf("%d.%s", t.AppID, t.AccessToken)
		return
	default:
		return
	}
}

// 从 yaml 配置文件中读取 token
func (t *Token) ReadFromConfig(filePath string) (_t *Token, err error) {
	_t = t
	var config Config
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("read token from file failed, err: %v", err)
	}
	if err = yaml.Unmarshal(content, &config); err != nil {
		log.Printf("parse config failed, err: %v", err)
	}
	t.AppID = config.AppID
	t.AccessToken = config.Token
	return
}
