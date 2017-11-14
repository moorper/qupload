# QUPLOAD

一个用 Golang 写的七牛云上传工具。使用七牛云官方的服务端 API 完成。
---
## 使用

1. `go install github.com/moorper/qupload`
2. 在用户根目录创建 `qupload.ini` 配置如下：
```
AccessKey = ***
SecretKey = ***
Bucket = bucket
BucketPublic = true # bool
BucketDomain = http://***
```
3. 使用命令 `qupload ~/upload.jpg` 上传，上传成功返回以下内容：
```
key : upload.jpg
url :  http://oph6h5t6t.bkt.clouddn.com/upload.jpg
markdown : ![upload.jpg](http://oph6h5t6t.bkt.clouddn.com/upload.jpg)
```
注：
MacOS 可直接拖拽图片到控制台获取图片路径