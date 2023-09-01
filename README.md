# bilibili-go

![GitHub](https://img.shields.io/github/license/kainhuck/bilibili-go) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/kainhuck/bilibili-go) ![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/kainhuck/bilibili-go)

使用`bilibili-go`来顺畅的接入bilibili，支持视频投稿，个人信息查询...

## Warning ⚠️
！！！目前项目处于开发阶段，并非稳定版本，接口可能会变更，其他接口陆续接入中🔨...

## 使用 🥑

1. 下载包
    ```bash
    go get github.com/kainhuck/bilibili-go
    ```

2. 导入包
    ```go
    import bilibili_go "github.com/kainhuck/bilibili-go"
    ```

3. 使用包
    参考 👉[demo](test/main.go)

## 更新日志 🐥

### v0.1.0
1. 支持扫码登陆
2. 支持缓存cookie
3. 支持视频上传
4. 支持封面上传
5. 支持投稿视频
6. 支持查询个人信息相关接口