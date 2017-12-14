# QUPLOAD

一个用 Golang 写的七牛云上传工具。使用七牛云官方的服务端 API 完成。
---
## 使用

1. `go install github.com/moorper/qupload`
2. 第一次运行会自动在 `~/qupload/qupload.ini` 创建配置文件：
```
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
```
3. 使用命令 `qupload ~/upload.jpg` 上传，上传成功返回以下内容：
```
key : upload.jpg
url :  http://oph6h5t6t.bkt.clouddn.com/upload.jpg
markdown : ![upload.jpg](http://oph6h5t6t.bkt.clouddn.com/upload.jpg)
```
注：
MacOS 可直接拖拽图片到控制台获取图片路径

## TODO
* 上传文件重命名
* 本地记录上传文件
* 免配置直接上传文件到 [https://sm.ms/](https://sm.ms/)

---

## 版本
### v0.0.1
* 简单配置即可使用 `qupload` 命令上传文件
### v0.0.2
* 修改配置文件路径为 `~/qupload/qupload.ini`
* 自动创建配置文件
* 支持多帐号配置
