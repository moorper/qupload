package main

import (
	"context"
	"flag"
	"fmt"
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
var config CFG

func init() {
	flag.BoolVar(&RequestVersion, "v", false, "查看当前版本")
	flag.Parse()
	initConfig()
}
func main() {
	if RequestVersion {
		fmt.Println("v0.0.1")
		return
	}
	//配置文件
	var AccessKey = config.AccessKey
	var SecretKey = config.SecretKey
	var Bucket = config.Bucket
	var BucketPublic = config.BucketPublic
	var BucketDomain = config.BucketDomain
	//上传文件检测
	localFile := flag.Arg(0)
	localFileInfo, err := os.Stat(localFile)
	if err != nil || os.IsNotExist(err) {
		log.Fatalln("no such file or directory")
	}
	var key = localFileInfo.Name()
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
	err = formUploader.PutFile(context.Background(), &ret, upToken, key, localFile, &putExtra)
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
	currentUser, _ := user.Current()
	homeDir := currentUser.HomeDir
	configPath := filepath.Join(homeDir, "qupload.ini")

	cfg, err := ini.Load(configPath)
	if err != nil {
		log.Fatalln(`
配置文件解析错误
在用户根目录下新建：qupload.ini
配置如下：
AccessKey = ***
SecretKey = ***
Bucket = bucket
BucketPublic = true # bool
BucketDomain = http://***
			`)
	}
	config.BucketPublic = true
	err = cfg.MapTo(&config)
	if err != nil {
		log.Fatalln(err)
	}
}
