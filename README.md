# httpfileserver

## httpfileserver主要用于http的方式上传文件，下载文件，已经列出一个目录下面的文件列表

### 使用方法

```
获取一个文件
$ curl "xxxx:12345/file/v1/xxxx"

推送一个文件
$ curl "xxxx:12345/file/v1/xxxx" -d @xxxx.txt

列举所有文件
$ curl "xxxx:12345/file/v1"
```

### 二进制构建

```
go build cmd/main.go
```

### docker构建

```
docker build -t httpfileserver -f make/Dockerfile .
```

### 使用方法

```
$ go run cmd/main.go -h
start http file server
it is expected to download file in compress mode via base auth

Usage:
  httpfileserver [flags]

Flags:
  -h, --help              help for httpfileserver
      --password string   user password
      --port uint16       the port exported outside (default 80)
      --root string       the root dir of file server (default ".")
      --user string       username (default "root")
```