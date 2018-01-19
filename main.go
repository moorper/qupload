package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/go-ini/ini"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
)

// CFG 配置文件
type CFG struct {
	AccessKey    string
	SecretKey    string
	Bucket       string
	BucketPublic bool
	BucketDomain string
}

// RequestVersion 查看当前版本
var RequestVersion bool

// Section 区分不同的上传帐号
var Section string

// Key 指定上传文件的Key
var Key string

// config 解析配置文件
var config CFG

func init() {
	flag.BoolVar(&RequestVersion, "v", false, "查看当前版本")
	flag.StringVar(&Section, "s", "default", "上传的 Section 空间")
	flag.StringVar(&Key, "k", "", "指定上传文件的 key")
	flag.Parse()
	initConfig()
}
func main() {
	if RequestVersion {
		fmt.Println("v0.0.3")
		return
	}
	//配置文件
	AccessKey := config.AccessKey
	SecretKey := config.SecretKey
	Bucket := config.Bucket
	BucketPublic := config.BucketPublic
	BucketDomain := config.BucketDomain
	//上传文件检测
	localFile := flag.Arg(0)
	localFileInfo, err := os.Stat(localFile)
	if err != nil || os.IsNotExist(err) {
		log.Fatalln("no such file or directory")
	}
	if Key == "" {
		Key = localFileInfo.Name()
	}
	// 获取 Bucket 对应的 zone
	BucketZone, err := storage.GetZone(AccessKey, Bucket)
	if err != nil {
		log.Fatalln("no such AccessKey or Bucket")
	}
	// mac
	mac := qbox.NewMac(AccessKey, SecretKey)
	// 获取上传token
	putPolicy := storage.PutPolicy{
		Scope: Bucket,
	}
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{
		Zone:          BucketZone, // 空间对应的机房
		UseHTTPS:      false,      // 是否使用https域名
		UseCdnDomains: false,      // 上传是否使用CDN上传加速
	}
	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	// 可选配置
	putExtra := storage.PutExtra{}
	err = formUploader.PutFile(context.Background(), &ret, upToken, Key, localFile, &putExtra)
	if err != nil {
		fmt.Println(err)
		return
	}
	if BucketPublic {
		publicURL := storage.MakePublicURL(BucketDomain, ret.Key)
		show(ret.Key, publicURL)
	} else {
		deadline := time.Now().Add(time.Second * 3600).Unix() //1小时有效期
		privateURL := storage.MakePrivateURL(mac, BucketDomain, ret.Key, deadline)
		show(ret.Key, privateURL)
	}
}
func show(key string, url string) {
	fmt.Println("key :", key)
	fmt.Println("url : ", url)
	fmt.Printf("markdown : ![%s](%s)", key, url)
	fmt.Println()
}

func initConfig() {
	var configDefault = `
[default]	
AccessKey = ***
SecretKey = ***
Bucket = bucket
BucketPublic = true # bool
BucketDomain = http://***
# 可不受限制的添加更多 Section 只需使用 qupload -s=other file.png 即可
# [other]
# AccessKey = ***
# SecretKey = ***
# Bucket = bucket
# BucketPublic = true # bool
# BucketDomain = http://***
	`
	currentUser, _ := user.Current()
	homeDir := currentUser.HomeDir
	quploadPath := filepath.Join(homeDir, ".qupload")
	_, err := os.Stat(quploadPath)
	if os.IsNotExist(err) {
		os.Mkdir(quploadPath, 0711)
	}
	configPath := filepath.Join(quploadPath, "qupload.ini")

	cfg, err := ini.Load(configPath)
	if err != nil {
		err = ioutil.WriteFile(configPath, []byte(configDefault), 0711)
		log.Fatal(`初始化配置文件成功，请编辑配置文件：`, configPath)
	}
	config.BucketPublic = true
	_, err = cfg.GetSection(Section)
	if err != nil {
		log.Fatalln("未找到指定的 Section")
	}
	cfg.Section(Section).MapTo(&config)
}
