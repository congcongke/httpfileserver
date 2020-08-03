# httpfileserver

## httpfileserver主要用于http的方式上传文件，下载文件，已经列出一个目录下面的文件列表

### 获取文件

```
获取一个文件
$ curl "xxxx:12345/file/v1/xxxx"

推送一个文件
$ curl "xxxx:12345/file/v1/xxxx" -d @xxxx.txt

列举所有文件
$ curl "xxxx:12345/file/v1"
```