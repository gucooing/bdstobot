## 项目名称 ： bdstobot （随意取的名字）

### 连接bds插件[PFLP ws](https://github.com/PixelFaramita/PixelFaramitaLuminousPolymerizationRes) 然后与[QQ(cqhttp)](https://github.com/Mrs4s/go-cqhttp)/discord互通

## 支持的功能：

1.进、退服发送消息提醒（discord使用的是webhook）

2.执行指令（仅QQ）

3.服务器状态获取（仅QQ）

4.服务器状态监控（断线警告）

5.服务器聊天与QQ群互通（discord目前仅支持单向互通）

6.绑定，解绑服务器白名单

7.远程启动服务器（仅QQ："mc 启动!"）

### 开发以discord bot为主（其实是QQ被风控不能发送消息）

## 支持内、外置discord bot（用来避免discor被墙的困扰）

### 外置discord bot通信采用rsa加密，公钥加密，私钥签名

## 请自行创建两对密钥对a和b，将a公钥与b私钥存放在外置discord botdata目录中替换，将b公钥与a私钥存放在本项目data目录中替换！！！！！请妥善保管好创建的ab两对密钥对，只要密钥对不泄露，数据传输99.9999999999%安全

### 还需要啥功能提issues就行了，一般都会加

## 使用指南

### 1.前往[actions](https://github.com/gucooing/bdstobot/actions)下载运行最新版

### 2.手动编译

1. 克隆仓库到本地：

   ```bash
   git clone https://github.com/gucooing/bdstobot.git
   ```

2. 编译和运行示例代码：

   ```bash
   build.bat
   ```


## 参与贡献

### 欢迎对该项目提供反馈和建议，你可以通过以下方式参与贡献：

- 提交问题和建议：在项目的 GitHub 仓库中提交 issue。
- 提交代码：如果你有修复 bug 或者改进功能的代码，可以提交 pull 请求。
