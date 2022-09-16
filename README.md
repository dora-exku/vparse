## 视频播放链接解析

本项目仅供学习交流

#### 目前已实现

- ✅  腾讯视频
- ✅  爱奇艺

## 测试方式

- 腾讯视频
  > 1、打开 https://v.qq.com 登录 通过 F12 获取到 cookie 将ck保存到根目录的 cookie_tencent.ck 文件中
  >
  > 2、获取  [v-algorithm](https://github.com/dora-exku/v-algorithm) 并使用node 运行
  > 
  > 3、运行 go run example/tencent.go

- 爱奇艺
  > 1、打开 https://www.iqiyi.com 登录 通过 F12 获取到 cookie 将ck保存到根目录的 cookie_iqiyi.ck 文件中
  >
  > 2、获取  [v-algorithm](https://github.com/dora-exku/v-algorithm) 并使用node 运行
  >
  > 3、运行 go run example/iqiyi.go