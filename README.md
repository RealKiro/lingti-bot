# lingti-bot (灵提)

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

**灵提**是一个极其易于集成的 MCP (Model Context Protocol) 服务器和多平台消息路由器，让 AI 助手能够访问你的本地计算机资源。

> **为什么叫"灵提"？** 灵缇犬（Greyhound）以其敏捷、忠诚著称。灵提 bot 同样敏捷高效，是你忠实的 AI 助手。

## 特性

- **MCP 标准协议** - 兼容所有支持 MCP 的 AI 客户端（Claude Desktop、Cursor 等）
- **多平台消息路由** - 同时支持 Slack、飞书、云中继等多个平台
- **多 AI 后端支持** - Claude、Kimi、DeepSeek 等主流 AI 服务
- **丰富的系统工具** - 文件操作、Shell 命令、系统信息、进程管理、网络工具
- **macOS 深度集成** - 日历、提醒事项、备忘录、音乐控制、截图等原生功能
- **实用工具集** - 天气查询、网页搜索、剪贴板、系统通知
- **跨平台支持** - 核心功能支持 macOS、Linux、Windows
- **极简配置** - 几行配置即可完成集成，开箱即用

## 架构

```
┌─────────────────────────────────────────────────────────────────┐
│                         lingti-bot                               │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────────┐      ┌──────────────┐      ┌──────────────┐   │
│  │  MCP Server  │      │   Message    │      │    Agent     │   │
│  │   (stdio)    │      │   Router     │      │   (Claude)   │   │
│  └──────┬───────┘      └──────┬───────┘      └──────┬───────┘   │
│         │                     │                     │            │
│         └─────────────────────┴─────────────────────┘            │
│                               │                                  │
│                               ▼                                  │
│                    ┌─────────────────────┐                       │
│                    │     MCP Tools       │                       │
│                    │  ┌───────┐ ┌─────┐  │                       │
│                    │  │ Files │ │Shell│  │                       │
│                    │  │System │ │ Net │  │                       │
│                    │  │Process│ │ Cal │  │                       │
│                    │  └───────┘ └─────┘  │                       │
│                    └─────────────────────┘                       │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘
         │                     │
         ▼                     ▼
┌─────────────────┐    ┌───────────────┐
│ Claude Desktop  │    │  Slack/飞书   │
│ Cursor / 其他   │    │  消息平台      │
└─────────────────┘    └───────────────┘
```

## 快速开始

### 安装

```bash
# 克隆仓库
git clone https://github.com/pltanton/lingti-bot.git
cd lingti-bot

# 编译
make build

# 或者编译到 dist 目录
make darwin-arm64  # Apple Silicon Mac
make darwin-amd64  # Intel Mac
make linux-amd64   # Linux x64
```

### 作为 MCP Server 使用

**最简单的使用方式** - 配置 Claude Desktop 或其他 MCP 客户端：

**Claude Desktop 配置** (`~/Library/Application Support/Claude/claude_desktop_config.json`)：

```json
{
  "mcpServers": {
    "lingti-bot": {
      "command": "/path/to/lingti-bot",
      "args": ["serve"]
    }
  }
}
```

**Cursor 配置**：

```json
{
  "mcpServers": {
    "lingti-bot": {
      "command": "/path/to/lingti-bot",
      "args": ["serve"]
    }
  }
}
```

就这么简单！重启客户端后，AI 助手即可使用所有 lingti-bot 提供的工具。

### 作为消息路由器使用

连接 Slack 或飞书等消息平台：

```bash
# 设置环境变量
export ANTHROPIC_API_KEY="sk-ant-your-api-key"

# Slack
export SLACK_BOT_TOKEN="xoxb-..."
export SLACK_APP_TOKEN="xapp-..."

# 飞书
export FEISHU_APP_ID="cli_..."
export FEISHU_APP_SECRET="..."

# 启动路由器
./lingti-bot router
```

## 项目结构

```
lingti-bot/
├── main.go                 # 程序入口
├── Makefile                # 构建脚本
├── go.mod                  # Go 模块定义
│
├── cmd/                    # 命令行接口
│   ├── root.go             # 根命令
│   ├── serve.go            # MCP 服务器命令
│   ├── service.go          # 系统服务管理
│   └── version.go          # 版本信息
│
├── internal/
│   ├── mcp/
│   │   └── server.go       # MCP 服务器实现
│   │
│   ├── tools/              # MCP 工具实现
│   │   ├── filesystem.go   # 文件读写、列表、搜索
│   │   ├── shell.go        # Shell 命令执行
│   │   ├── system.go       # 系统信息、磁盘、环境变量
│   │   ├── process.go      # 进程列表、信息、终止
│   │   ├── network.go      # 网络接口、连接、DNS
│   │   ├── calendar.go     # macOS 日历集成
│   │   ├── filemanager.go  # 文件整理（清理旧文件）
│   │   ├── reminders.go    # macOS 提醒事项
│   │   ├── notes.go        # macOS 备忘录
│   │   ├── weather.go      # 天气查询（wttr.in）
│   │   ├── websearch.go    # 网页搜索和获取
│   │   ├── clipboard.go    # 剪贴板读写
│   │   ├── notification.go # 系统通知
│   │   ├── screenshot.go   # 屏幕截图
│   │   └── music.go        # 音乐控制（Spotify/Apple Music）
│   │
│   ├── router/
│   │   └── router.go       # 多平台消息路由器
│   │
│   ├── platforms/          # 消息平台集成
│   │   ├── slack/
│   │   │   └── slack.go    # Slack Socket Mode
│   │   └── feishu/
│   │       └── feishu.go   # 飞书 WebSocket
│   │
│   ├── agent/
│   │   ├── tools.go        # Agent 工具执行
│   │   └── memory.go       # 会话记忆
│   │
│   └── service/
│       └── manager.go      # 系统服务管理
│
└── docs/                   # 文档
    ├── slack-integration.md    # Slack 集成指南
    ├── feishu-integration.md   # 飞书集成指南
    └── openclaw-reference.md   # 架构参考
```

## MCP 工具一览

lingti-bot 提供以下 MCP 工具供 AI 助手调用：

### 文件操作

| 工具 | 功能 |
|------|------|
| `file_read` | 读取文件内容 |
| `file_write` | 写入文件内容 |
| `file_list` | 列出目录内容 |
| `file_search` | 按模式搜索文件 |
| `file_info` | 获取文件详细信息 |

### Shell 命令

| 工具 | 功能 |
|------|------|
| `shell_execute` | 执行 Shell 命令 |
| `shell_which` | 查找可执行文件路径 |

### 系统信息

| 工具 | 功能 |
|------|------|
| `system_info` | 获取系统信息（CPU、内存、OS） |
| `disk_usage` | 获取磁盘使用情况 |
| `env_get` | 获取环境变量 |
| `env_list` | 列出所有环境变量 |

### 进程管理

| 工具 | 功能 |
|------|------|
| `process_list` | 列出运行中的进程 |
| `process_info` | 获取进程详细信息 |
| `process_kill` | 终止进程 |

### 网络工具

| 工具 | 功能 |
|------|------|
| `network_interfaces` | 列出网络接口 |
| `network_connections` | 列出活动网络连接 |
| `network_ping` | TCP 连接测试 |
| `network_dns_lookup` | DNS 查询 |

### 日历（macOS）

| 工具 | 功能 |
|------|------|
| `calendar_today` | 获取今日日程 |
| `calendar_list_events` | 列出未来事件 |
| `calendar_create_event` | 创建日历事件 |
| `calendar_search` | 搜索日历事件 |
| `calendar_delete_event` | 删除日历事件 |
| `calendar_list_calendars` | 列出所有日历 |

### 文件整理

| 工具 | 功能 |
|------|------|
| `file_list_old` | 列出长时间未修改的文件 |
| `file_delete_old` | 删除长时间未修改的文件 |
| `file_delete_list` | 批量删除指定文件 |
| `file_trash` | 移动文件到废纸篓（macOS） |

### 提醒事项（macOS）

| 工具 | 功能 |
|------|------|
| `reminders_today` | 获取今日待办事项 |
| `reminders_add` | 添加新提醒 |
| `reminders_complete` | 标记提醒为已完成 |
| `reminders_delete` | 删除提醒 |

### 备忘录（macOS）

| 工具 | 功能 |
|------|------|
| `notes_list_folders` | 列出备忘录文件夹 |
| `notes_list` | 列出备忘录 |
| `notes_read` | 读取备忘录内容 |
| `notes_create` | 创建新备忘录 |
| `notes_search` | 搜索备忘录 |

### 天气

| 工具 | 功能 |
|------|------|
| `weather_current` | 获取当前天气 |
| `weather_forecast` | 获取天气预报 |

### 网页搜索

| 工具 | 功能 |
|------|------|
| `web_search` | DuckDuckGo 搜索 |
| `web_fetch` | 获取网页内容 |

### 剪贴板

| 工具 | 功能 |
|------|------|
| `clipboard_read` | 读取剪贴板内容 |
| `clipboard_write` | 写入剪贴板 |

### 系统通知

| 工具 | 功能 |
|------|------|
| `notification_send` | 发送系统通知 |

### 截图

| 工具 | 功能 |
|------|------|
| `screenshot` | 截取屏幕截图 |

### 音乐控制（macOS）

| 工具 | 功能 |
|------|------|
| `music_play` | 播放音乐 |
| `music_pause` | 暂停音乐 |
| `music_next` | 下一首 |
| `music_previous` | 上一首 |
| `music_now_playing` | 获取当前播放信息 |
| `music_volume` | 设置音量 |
| `music_search` | 搜索并播放音乐 |

### 其他

| 工具 | 功能 |
|------|------|
| `open_url` | 在浏览器中打开 URL |

## Make 目标

```bash
# 开发
make build          # 编译当前平台
make run            # 本地运行
make test           # 运行测试
make fmt            # 格式化代码
make lint           # 代码检查
make clean          # 清理构建产物
make version        # 显示版本

# 跨平台编译
make darwin-arm64   # macOS Apple Silicon
make darwin-amd64   # macOS Intel
make darwin-universal # macOS 通用二进制
make linux-amd64    # Linux x64
make linux-arm64    # Linux ARM64
make linux-all      # 所有 Linux 平台
make all            # 所有平台

# 服务管理
make install        # 安装为系统服务
make uninstall      # 卸载系统服务
make start          # 启动服务
make stop           # 停止服务
make status         # 查看服务状态

# macOS 签名
make codesign       # 代码签名（需要开发者证书）
```

## 环境变量

| 变量 | 说明 | 必需 |
|------|------|------|
| `ANTHROPIC_API_KEY` | Anthropic API 密钥 | 路由器模式必需 |
| `ANTHROPIC_BASE_URL` | 自定义 API 地址 | 可选 |
| `ANTHROPIC_MODEL` | 使用的模型 | 可选 |
| `SLACK_BOT_TOKEN` | Slack Bot Token (`xoxb-...`) | Slack 集成必需 |
| `SLACK_APP_TOKEN` | Slack App Token (`xapp-...`) | Slack 集成必需 |
| `FEISHU_APP_ID` | 飞书 App ID | 飞书集成必需 |
| `FEISHU_APP_SECRET` | 飞书 App Secret | 飞书集成必需 |

## 文档

详细的集成指南请参阅：

- [Slack 集成指南](docs/slack-integration.md) - 完整的 Slack 应用配置教程
- [飞书集成指南](docs/feishu-integration.md) - 飞书/Lark 应用配置教程
- [架构参考](docs/openclaw-reference.md) - OpenClaw 架构参考

## 为什么选择 lingti-bot？

### 极简集成

与其他 MCP 服务器相比，lingti-bot 的集成极其简单：

```json
{
  "mcpServers": {
    "lingti-bot": {
      "command": "/path/to/lingti-bot",
      "args": ["serve"]
    }
  }
}
```

**无需额外配置**，无需数据库，无需 Docker，无需云服务。一个二进制文件，两行配置，即可获得完整的系统访问能力。

### 单一二进制

lingti-bot 编译为单个静态二进制文件，无外部依赖：

```bash
# 编译
make build

# 即可使用
./dist/lingti-bot serve
```

### 多平台支持

一套代码，同时支持：
- MCP stdio 模式（Claude Desktop、Cursor 等）
- Slack Socket Mode
- 飞书 WebSocket
- 更多平台陆续添加中...

### 本地优先

所有功能都在本地运行，数据不会上传到云端。你的文件、日历、进程信息都安全地保留在本地。

## 使用示例

配置完成后，你可以让 AI 助手执行以下操作：

### 日历与日程

```
"今天有什么日程安排？"
"这周有哪些会议？"
"帮我创建一个明天下午3点的会议，标题是'产品评审'"
"下周一上午10点提醒我给客户打电话"
"搜索所有包含'周报'的日程"
"删除明天的牙医预约"
```

### 提醒事项

```
"我今天有哪些待办事项？"
"添加一个提醒：明天下午2点前提交报告"
"把'买牛奶'标记为已完成"
"删除'取快递'这个提醒"
"帮我列出所有未完成的提醒"
```

### 备忘录

```
"帮我创建一个备忘录，标题是'会议纪要'，内容是..."
"读取我的'购物清单'备忘录"
"搜索包含'密码'的备忘录"
"列出我所有的备忘录文件夹"
"我的备忘录里有多少条笔记？"
```

### 天气查询

```
"北京今天天气怎么样？"
"上海未来三天的天气预报"
"深圳现在的温度是多少？"
"东京这周会下雨吗？"
```

### 网页搜索与浏览

```
"帮我搜索一下最新的 AI 新闻"
"查一下 Python 3.12 有什么新特性"
"打开 GitHub 首页"
"获取这个网页的内容：https://example.com"
"搜索附近的咖啡店"
```

### 文件操作

```
"列出桌面上的所有文件"
"读取 ~/Documents/notes.txt 的内容"
"在下载文件夹里搜索所有 PDF 文件"
"这个文件是什么时候创建的？"
"桌面上超过30天没动过的文件有哪些？"
"帮我把这些旧文件移到废纸篓"
```

### Shell 命令

```
"运行 git status"
"执行 npm install"
"看看当前目录下有什么"
"运行 docker ps 查看容器状态"
"python 安装在哪里？"
```

### 系统信息

```
"我的电脑配置是什么？"
"现在 CPU 占用多少？"
"内存还剩多少？"
"磁盘空间还有多少？"
"查看 PATH 环境变量"
```

### 进程管理

```
"现在有哪些进程在运行？"
"Chrome 占用了多少内存？"
"结束 PID 1234 的进程"
"有没有叫 node 的进程？"
```

### 网络工具

```
"我的 IP 地址是什么？"
"测试一下能不能连接到 google.com"
"查询 github.com 的 DNS"
"列出所有网络连接"
```

### 剪贴板

```
"剪贴板里有什么？"
"把这段文字复制到剪贴板：Hello World"
"读取剪贴板内容并翻译成英文"
```

### 截图

```
"帮我截个屏"
"截取当前窗口的截图"
"截图保存到桌面"
```

### 系统通知

```
"5分钟后提醒我休息一下"
"发送一个通知：任务已完成"
```

### 音乐控制

```
"播放音乐"
"暂停"
"下一首"
"上一首"
"现在在放什么歌？"
"音量调到 50%"
"播放周杰伦的歌"
"搜索并播放 Shape of You"
```

### 组合任务

```
"查看今天的日程，然后检查天气，最后列出待办事项"
"帮我整理桌面：列出超过60天的旧文件，然后移到废纸篓"
"搜索最近的科技新闻，整理成备忘录"
"截个图然后用通知提醒我已完成"
```

## 安全注意事项

- lingti-bot 提供对本地系统的访问能力，请在可信环境中使用
- Shell 命令执行有基本的危险命令过滤，但仍需谨慎
- API 密钥等敏感信息请使用环境变量，不要提交到版本控制
- 生产环境建议使用专用服务账号运行

## 依赖

- [mcp-go](https://github.com/mark3labs/mcp-go) - MCP 协议 Go 实现
- [cobra](https://github.com/spf13/cobra) - CLI 框架
- [gopsutil](https://github.com/shirou/gopsutil) - 系统信息
- [slack-go](https://github.com/slack-go/slack) - Slack SDK
- [oapi-sdk-go](https://github.com/larksuite/oapi-sdk-go) - 飞书/Lark SDK
- [go-anthropic](https://github.com/liushuangls/go-anthropic) - Anthropic API 客户端

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！

---

**灵提** - 你的敏捷 AI 助手 🐕
