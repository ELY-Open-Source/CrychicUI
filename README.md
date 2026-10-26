# CrychicUI

MyGo 原生 UI 的 Go 组件库：1500 个组件，分 38 类，每类一个包。沿用 MyGo 自己的外观：锌灰底色，跟随系统的强调色，6 DIP 圆角，4 DIP 间距单位，系统字体。每个组件都会动，每个可交互组件都有合成音效。

界面由 MyGo 在 GPU 上自绘，不用 HTML、JavaScript 或 WebView。

## 安装

```sh
go get github.com/ELY-Open-Source/CrychicUI
```

要求 Go 1.27.1 与 MyGo v0.5.5。

## 使用

组件是函数：第一个参数是 `*ui.Context`，返回 `ui.Element`；可选项放在最后的 `NameOptions` 里。容器把内容作为最后一个 `func()`，放在选项之后。有多个动作的组件返回多个命名的 `ui.Element`，`root` 在前，例如 `button.SplitButton`。

```go
import (
	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/egoist/mygo/ui"
)

func view(c *ui.Context) {
	if button.Solid(c, "Save").Clicked() {
		save() // 按钮已经播放了自己的音效
	}
}
```

## 分类

| 类 | 包 | 数量 |
| --- | --- | ---: |
| Layout | `layout` | 50 |
| Typography | `typography` | 40 |
| Button | `button` | 45 |
| Display | `display` | 55 |
| Input | `input` | 55 |
| Choice | `choice` | 50 |
| Feedback | `feedback` | 55 |
| Motion | `motion` | 55 |
| Effects | `fx` | 45 |
| Form | `form` | 40 |
| Date & Time | `datetime` | 45 |
| Color | `color` | 20 |
| Navigation | `navigation` | 55 |
| Overlay | `overlay` | 45 |
| Data | `data` | 55 |
| Chart | `chart` | 60 |
| Media | `media` | 45 |
| Chat | `chat` | 45 |
| Agent | `agent` | 40 |
| Code | `code` | 40 |
| Files | `files` | 35 |
| Shop | `shop` | 40 |
| Finance | `finance` | 40 |
| Dashboard | `dashboard` | 35 |
| Schedule | `schedule` | 30 |
| Project | `project` | 40 |
| Mail | `mail` | 30 |
| Social | `social` | 35 |
| Account | `account` | 40 |
| Devtools | `devtools` | 35 |
| Game | `game` | 30 |
| Learn | `learn` | 20 |
| Health | `health` | 20 |
| Home | `home` | 20 |
| Marketing | `marketing` | 40 |
| Content | `content` | 25 |
| Desktop | `desktop` | 30 |
| Geo | `geo` | 15 |

`sound`、`icon`、`kit` 是组件共用的底层包。

## 音效

可交互组件在动作时播放一个音效，默认开启。音效在 `sound` 包里合成，不读音频文件。

```go
sound.SetEnabled(false) // 全部静音
sound.SetVolume(0.5)
```

## Gallery

```sh
go run ./cmd/gallery
```

同一个程序编译成 WebAssembly，就是网站：https://crychicui.zhangyanghaha0407.workers.dev

## 检查

```sh
go vet ./... && go tool mygo vet ./... && go test ./...
GOOS=js GOARCH=wasm go vet ./...
cd web && npm install && npm test   # 构建 web/dist，Playwright 本地跑
```

网站检查需要 Node 22+（Wrangler 的要求）、PATH 里的 `python3`（本地测试服务器），以及一次 `npx playwright install chromium`。

`npm run deploy`（在 `web/` 里）先在预览地址测试新版，再切到线上并重测。

## 许可

MIT。第三方声明见 [THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md)。
