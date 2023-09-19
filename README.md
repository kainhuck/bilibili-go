![bilibili-go](https://socialify.git.ci/kainhuck/bilibili-go/image?description=1&descriptionEditable=%E7%AE%80%E5%8D%95%E5%A5%BD%E7%94%A8%E7%9A%84%20bilibili%20golang%20sdk&font=Inter&forks=1&issues=1&language=1&name=1&owner=1&pattern=Floating%20Cogs&pulls=1&stargazers=1&theme=Auto)

## 简介 📜

![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/kainhuck/bilibili-go) 
![GitHub](https://img.shields.io/github/license/kainhuck/bilibili-go) 
![GitHub tag (with filter)](https://img.shields.io/github/v/tag/kainhuck/bilibili-go)

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

4. options介绍

   1. 自定义http客户端
      用户可以通过自定义http客户端来使用代理或者其他需求，比如
         
      ```go
      proxyURL, err := url.Parse("http://proxy.example.com:8080")
      if err != nil {
          panic(err)
      }
      
      client := bilibili_go.NewClient(
          bilibili_go.WithHttpClient(&http.Client{
              Transport: &http.Transport{
                  Proxy: http.ProxyURL(proxyURL),
              },
          }),
      )
      ```
   
   2. 缓存cookie
      
      用户可以实现下面这个接口来定义自己的存储，如果设置了缓存，在第二次登陆时将不再需要授权，除非缓存过期，在加载缓存时会自动校验授权信息是否过期
      ```go
      type AuthStorage interface {
          // LoadAuthInfo 加载AuthInfo
          LoadAuthInfo() (*AuthInfo, error)

          // SaveAuthInfo 保存AuthInfo
          SaveAuthInfo(*AuthInfo) error
      
          // LogoutAuthInfo 账号退出登陆时会调用该方法
          LogoutAuthInfo(*AuthInfo) error
      }
      ```
      默认提供了一个文件缓存的实现`fileAuthStorage`可以如下使用
      ```go
      cient := bilibili_go.NewClient(
           bilibili_go.WithAuthStorage(bilibili_go.NewFileAuthStorage("文件路径")),
      )
      ```
         
   3. 开启调试
      
      开启debug模式后，将会向指定的文件中写入http的报文
      ```go
      client := bilibili_go.NewClient(
           bilibili_go.WithDebug(true), // 将会向 stdout 输出http报文
      )
      ```
      ```go
      f, err := os.Open("debug.txt")
      if err != nil {
          panic(err)
      }
      defer f.Close()

      client := bilibili_go.NewClient(
          bilibili_go.WithDebug(true, f), // 将会向 debug.txt 输出http报文
      )
      ```
      
   4. 自定义处理登陆二维码
      
      在使用`LoginWithQrCode`方法登陆时，默认会将登陆二维码输出到标准输出，用户可以配置自己的输出方法来自定义处理登陆二维码，比如将其发送到指定的群组或个人
      ```go
      client := bilibili_go.NewClient(
		  bilibili_go.WithShowQRCodeFunc(func(code *qrcode.QRCode) error {
			  // ....
			  return nil
          }),
      )
      ```
      
   5. 自定义User-Agent 
      
      默认UA: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36"
      ```go
      client := bilibili_go.NewClient(
		bilibili_go.WithUserAgent("abc"),
      )
      ```
   
   6. 使用自定义logger
      
      用户可以设置任何实现了`Logger`接口的日志，默认使用`logrus.StandardLogger()`
      ```go
      type Logger interface {
          Debug(args ...any)
          Info(args ...any)
          Warn(args ...any)
          Error(args ...any)
          Debugf(format string, args ...any)
          Infof(format string, args ...any)
          Warnf(format string, args ...any)
          Errorf(format string, args ...any)
      }
      ```
      ```go
      client := bilibili_go.NewClient(
          bilibili_go.WithLogger(log),
      )
      ```


## 特别鸣谢 🥰

[bilibili-API-collect](https://github.com/SocialSisterYi/bilibili-API-collect)

## 更新日志 🐥

### v0.3.5
1. 新增token刷新功能
2. 查询每日奖励状态
3. 视频投币

### v0.3.4
1. 新增登出功能

### v0.3.3
1. 封装了关系操作接口
2. 新增查询用户与自己的关系接口
3. 新增查询用户与自己的互相关系接口
4. 新增批量查询用户与自己的关系
5. 新增查询关注分组列表
6. 新增其他分组操作

### v0.3.2
1. 新增操作用户关系接口
2. 新增批量操作用户关系接口

### v0.3.1
1. 新增查询用户粉丝列表接口
2. 新增查询用户关注列表接口
3. 新增搜索用户关注列表接口
4. 新增查询共同关注列表接口
5. 新增查询悄悄关注列表接口
6. 新增查询互相关注列表接口
7. 新增查询黑名单列表接口

### v0.3.0
1. 新增关系状态数接口
2. 新增up状态数接口
3. 新增相簿投稿数接口
4. 优化封面上传方式，支持从io读取，http读取
5. auth info 缓存用户信息

### v0.2.6
1. 更新auth info的存储模块
2. 可自定义auth info存储

### v0.2.5
1. 更新分区

### v0.2.4
1. 可自定义处理登陆二维码

### v0.2.3
1. 添加视频分区

### v0.2.2
1. 更新debug调试，可输出到指定文件
2. 新增日志配置，可使用自定义日志
3. 更新多渠道视频上传，本地磁盘，io，http

### v0.2.1
1. 更新认证信息缓存
2. 新增refresh_token

### v0.2.0
1. 修改登陆接口
2. 增加debug参数
3. 新增 获取硬币数 接口
4. 新增 用户空间详细信息 接口
5. 新增 用户名片信息 接口
6. 新增 登陆用户空间详细信息 接口

### v0.1.0
1. 支持扫码登陆
2. 支持缓存cookie
3. 支持视频上传
4. 支持封面上传
5. 支持投稿视频
6. 支持查询个人信息相关接口
