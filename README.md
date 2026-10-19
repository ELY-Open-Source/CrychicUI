# CrychicUI

MyGo 原生 UI 的 Go 组件库。沿用 MyGo 自己的外观：锌灰底色，跟随系统的强调色，6 DIP 圆角，4 DIP 间距单位，系统字体。每个组件都会动，每个可交互组件都有合成音效。

界面由 MyGo 在 GPU 上自绘，不用 HTML、JavaScript 或 WebView。

## 安装

```sh
go get github.com/ELY-Open-Source/CrychicUI
```

要求 Go 1.27.1 与 MyGo v0.5.5。

## 音效

可交互组件在动作时播放一个音效，默认开启。

```go
sound.SetEnabled(false) // 全部静音
```

## Gallery

```sh
go run ./cmd/gallery
```

同一个程序编译成 WebAssembly，就是网站：https://crychicui.zhangyanghaha0407.workers.dev

## 许可

MIT。第三方声明见 [THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md)。
