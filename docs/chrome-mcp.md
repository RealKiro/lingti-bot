# Chrome DevTools MCP Server

> 使用 Chrome DevTools Protocol (CDP) 通过 MCP 协议控制浏览器的完整指南。

---

## 什么是 Chrome DevTools MCP Server？

Chrome DevTools MCP Server 是一个基于 [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) 的浏览器自动化工具。它允许 AI 通过标准化的 MCP 接口控制 Chrome 浏览器，执行导航、点击、输入、截图等操作。

### 与 lingti-bot 内置浏览器的区别

| | Chrome DevTools MCP | lingti-bot 内置浏览器 |
|--|--|--|
| **协议** | MCP (Model Context Protocol) | 直接 CDP + go-rod |
| **安装** | 需要单独安装 MCP server | 开箱即用 |
| **适用范围** | 任何 MCP 兼容的 AI 客户端 | lingti-bot 专用 |
| **连接方式** | 连接已运行的 Chrome | 可启动新实例或连接 |

---

## 安装

### 前置条件

- Node.js 18+
- Chrome/Brave/Edge 浏览器

### 安装 MCP Server

```bash
# 全局安装
npm install -g chrome-devtools-mcp

# 或使用 npx 直接运行
npx chrome-devtools-mcp --help
```

---

## 配置

### 方式一：自动启动 Chrome（推荐）

`chrome-devtools-mcp` 可以自动启动 Chrome，无需手动配置调试端口：

```json
{
  "mcpServers": {
    "chrome-devtools": {
      "command": "npx",
      "args": ["chrome-devtools-mcp"]
    }
  }
}
```

**常用选项：**

| 参数 | 说明 |
|------|------|
| `--headless` | 无头模式运行 |
| `--viewport 1280x720` | 设置视口大小 |
| `--channel stable` | 使用稳定版 Chrome（默认） |
| `--channel canary` | 使用 Chrome Canary |
| `--user-data-dir` | 指定用户数据目录 |
| `--autoConnect` | 连接到正在运行的 Chrome（Chrome 144+） |

### 完整配置选项

| 参数 | 缩写 | 说明 | 默认值 |
|------|------|------|--------|
| `--browserUrl` | `-u` | 连接已运行的 Chrome（如 `http://127.0.0.1:9222`） | - |
| `--wsEndpoint` | `-w` | WebSocket 端点连接 | - |
| `--wsHeaders` | - | WebSocket 自定义请求头（JSON 格式） | - |
| `--headless` | - | 无头模式运行 | false |
| `--executablePath` | `-e` | 指定 Chrome 可执行文件路径 | 自动检测 |
| `--isolated` | - | 创建临时用户数据目录，关闭后自动清理 | false |
| `--userDataDir` | - | 用户数据目录路径 | `~/.cache/chrome-devtools-mcp/chrome-profile` |
| `--channel` | - | Chrome 渠道：`stable`, `canary`, `beta`, `dev` | stable |
| `--viewport` | - | 视口大小（如 `1280x720`） | - |
| `--proxyServer` | - | 代理服务器配置 | - |
| `--acceptInsecureCerts` | - | 忽略自签名证书错误 | false |
| `--chromeArg` | - | 额外 Chrome 参数（可多次使用） | - |
| `--ignoreDefaultChromeArg` | - | 禁用默认 Chrome 参数 | - |
| `--categoryEmulation` | - | 启用模拟工具 | true |
| `--categoryPerformance` | - | 启用性能工具 | true |
| `--categoryNetwork` | - | 启用网络工具 | true |
| `--performanceCrux` | - | 启用 CrUX 性能数据 | true |
| `--usageStatistics` | - | 发送使用统计 | true |
| `--autoConnect` | - | 自动连接到运行中的 Chrome（Chrome 144+） | false |
| `--logFile` | - | 日志文件路径 | - |

### 方式二：连接已有 Chrome

如果你已有运行中的 Chrome，可以用调试端口连接：

**1. 启动 Chrome（带调试端口）：**

```bash
# macOS
/Applications/Google\ Chrome.app/Contents/MacOS/Google\ Chrome \
  --remote-debugging-port=9222 \
  --user-data-dir="$HOME/.chrome-mcp-profile"

# Linux
google-chrome --remote-debugging-port=9222 --user-data-dir="$HOME/.chrome-mcp-profile"
```

**2. 验证端口：**

```bash
curl http://localhost:9222/json/version
```

**3. 配置 MCP 客户端：**

```json
{
  "mcpServers": {
    "chrome-devtools": {
      "command": "npx",
      "args": ["chrome-devtools-mcp", "--browserUrl", "http://127.0.0.1:9222"]
    }
  }
}
```

### 配置示例

**Cursor / Claude Code 设置：**

在 `~/.cursor/user/settings.json` 或 `~/.claude/settings.json` 中添加：

```json
{
  "mcpServers": {
    "chrome-devtools": {
      "command": "npx",
      "args": ["chrome-devtools-mcp"]
    }
  }
}
```

**无头模式：**

```json
{
  "mcpServers": {
    "chrome-devtools": {
      "command": "npx",
      "args": ["chrome-devtools-mcp", "--headless", "--viewport", "1920x1080"]
    }
  }
}
```

**连接到已有 Chrome（需先手动启动）：**

```json
{
  "mcpServers": {
    "chrome-devtools": {
      "command": "npx",
      "args": ["chrome-devtools-mcp", "--browserUrl", "http://127.0.0.1:9222"]
    }
  }
}
```

**使用 Chrome Canary：**

```json
{
  "mcpServers": {
    "chrome-devtools": {
      "command": "npx",
      "args": ["chrome-devtools-mcp", "--channel", "canary"]
    }
  }
}
```

**使用自定义用户数据目录：**

```json
{
  "mcpServers": {
    "chrome-devtools": {
      "command": "npx",
      "args": ["chrome-devtools-mcp", "--user-data-dir", "/path/to/profile"]
    }
  }
}
```

**忽略证书错误（自签名证书）：**

```json
{
  "mcpServers": {
    "chrome-devtools": {
      "command": "npx",
      "args": ["chrome-devtools-mcp", "--accept-insecure-certs"]
    }
  }
}
```

**自动连接 Chrome 144+（无需调试端口）：**

```json
{
  "mcpServers": {
    "chrome-devtools": {
      "command": "npx",
      "args": ["chrome-devtools-mcp", "--auto-connect"]
    }
  }
}
```

**禁用默认 Chrome 参数：**

```json
{
  "mcpServers": {
    "chrome-devtools": {
      "command": "npx",
      "args": ["chrome-devtools-mcp", "--ignore-default-chrome-arg=--disable-extensions"]
    }
  }
}
```

**使用自定义 Chrome 路径：**

```json
{
  "mcpServers": {
    "chrome-devtools": {
      "command": "npx",
      "args": ["chrome-devtools-mcp", "--executable-path", "/Applications/Brave Browser.app/Contents/MacOS/Brave Browser"]
    }
  }
}
```

### 重启 MCP 客户端

重启 Cursor/Claude Code 后，MCP 工具即可使用。

---

## 工具完整参考

### 页面管理

#### `list_pages` — 列出所有标签页

返回浏览器中所有打开的页面：

```json
[
  {"id": 1, "url": "https://www.google.com", "title": "Google"},
  {"id": 2, "url": "https://github.com", "title": "GitHub"}
]
```

#### `new_page` — 打开新标签页

| 参数 | 类型 | 说明 |
|------|------|------|
| `url` | string | 打开的 URL（可选） |
| `background` | bool | 是否在后台打开，默认 false |

```bash
new_page url="https://example.com"
new_page                              # 打开空白页
```

#### `select_page` — 切换标签页

| 参数 | 类型 | 说明 |
|------|------|------|
| `pageId` | number | 页面 ID（从 `list_pages` 获取） |
| `bringToFront` | bool | 是否聚焦到该页面 |

```bash
select_page pageId=2
```

#### `close_page` — 关闭标签页

```bash
close_page pageId=2
```

#### `navigate_page` — 导航到 URL

| 参数 | 类型 | 说明 |
|------|------|------|
| `url` | string | 目标 URL |
| `type` | string | 类型：`url`（默认）, `back`, `forward`, `reload` |
| `ignoreCache` | bool | 忽略缓存（reload 时） |
| `timeout` | number | 超时毫秒数 |

```bash
navigate_page url="https://example.com"
navigate_page type="back"
navigate_page type="reload" ignoreCache=true
```

---

### 元素交互

> 所有交互工具都需要先执行 `take_snapshot` 获取元素的 `uid`。

#### `take_snapshot` — 获取页面快照

返回页面的可访问性树，每个元素带唯一标识符（uid）：

```
[1] link "首页"
[2] link "发现"
[3] textbox "搜索"
[4] button "搜索"
```

**参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| `verbose` | bool | 是否包含完整信息，默认 false |

#### `click` — 点击元素

```bash
click uid=4
click uid=4 dblClick=true    # 双击
```

#### `hover` — 悬停元素

```bash
hover uid=3
```

#### `fill` — 输入文本

| 参数 | 类型 | 说明 |
|------|------|------|
| `uid` | string | 元素 uid |
| `value` | string | 输入的文本 |

```bash
fill uid=3 value="搜索内容"
```

#### `fill_form` — 填写表单

```bash
fill_form elements=[{"uid": "1", "value": "张三"}, {"uid": "2", "value": "13812345678"}]
```

#### `press_key` — 按键

支持的按键组合：
- 单键：`Enter`, `Tab`, `Escape`, `Backspace`, `Delete`, `Space`
- 方向键：`ArrowUp`, `ArrowDown`, `ArrowLeft`, `ArrowRight`
- 修饰键：`Control`, `Shift`, `Alt`, `Meta`

```bash
press_key key="Enter"
press_key key="Control+A"
press_key key="Control+Shift+T"
```

#### `drag` — 拖拽

```bash
drag from_uid="1" to_uid="5"
```

#### `upload_file` — 上传文件

```bash
upload_file uid="1" filePath="/path/to/file.pdf"
```

---

### 内容获取

#### `screenshot` — 截图

| 参数 | 类型 | 说明 |
|------|------|------|
| `uid` | string | 元素 uid（可选，截图特定元素） |
| `fullPage` | bool | 是否截取整页 |
| `format` | string | 格式：`png`（默认）, `jpeg`, `webp` |
| `quality` | number | 质量 0-100（仅 jpeg/webp） |
| `filePath` | string | 保存路径 |

```bash
screenshot
screenshot fullPage=true
screenshot filePath="/tmp/screenshot.png"
```

#### `evaluate_script` — 执行 JavaScript

```bash
evaluate_script function="() => document.title"
evaluate_script function="() => window.scrollTo(0, document.body.scrollHeight)"
```

---

### 等待与同步

#### `wait_for` — 等待文本出现

```bash
wait_for text="加载完成"
wait_for text="登录成功" timeout=30000
```

---

### 网络请求

#### `list_network_requests` — 列出网络请求

```bash
list_network_requests
list_network_requests resourceTypes=["xhr", "fetch"]
```

#### `get_network_request` — 获取请求详情

```bash
get_network_request reqid=1
get_network_request reqid=1 requestFilePath="/tmp/request.json" responseFilePath="/tmp/response.json"
```

---

### 浏览器设置

#### `emulate` — 模拟设备和网络

| 参数 | 类型 | 说明 |
|------|------|------|
| `viewport` | object | 视口设置 |
| `userAgent` | string | 用户代理 |
| `colorScheme` | string | 颜色方案：`dark`, `light`, `auto` |
| `networkConditions` | string | 网络节流：`Slow 3G`, `Fast 3G`, `Offline` 等 |
| `geolocation` | object | 地理位置 |

```bash
emulate viewport={"width": 375, "height": 667, "isMobile": true}
emulate colorScheme="dark"
emulate networkConditions="Slow 3G"
```

#### `resize_page` — 调整窗口大小

```bash
resize_page width=1920 height=1080
```

---

## 典型使用场景

### 场景一：信息查询

```
用户: "帮我查一下 BTC 价格"

AI:
1. navigate_page url="https://www.coindesk.com"
2. take_snapshot
3. 从快照中提取价格信息返回
```

### 场景二：登录操作

```
用户: "帮我登录 GitHub"

AI:
1. navigate_page url="https://github.com/login"
2. take_snapshot
3. fill uid=<username输入框> value="your@email.com"
4. fill uid=<password输入框> value="yourpassword"
5. click uid=<登录按钮>
6. take_snapshot 确认登录成功
```

### 场景三：表单填写

```
AI:
1. navigate_page url="https://example.com/form"
2. take_snapshot
3. fill_form elements=[{"uid": "1", "value": "张三"}, {"uid": "2", "value": "工程师"}]
4. click uid=<提交按钮>
```

### 场景四：批量操作

```
用户: "点击所有弹窗的关闭按钮"

AI:
1. navigate_page url="https://example.com"
2. take_snapshot
3. evaluate_script function="() => document.querySelectorAll('.close-btn').forEach(el => el.click())"
```

---

## 故障排除

### 连接失败

```
Error: Cannot connect to Chrome at http://localhost:9222
```

**解决：**
1. 确认 Chrome 以调试端口运行：`curl http://localhost:9222/json/version`
2. 如果未运行，重新启动 Chrome：
   ```bash
   /Applications/Google\ Chrome.app/Contents/MacOS/Google\ Chrome \
     --remote-debugging-port=9222 \
     --user-data-dir="$HOME/.chrome-mcp-profile"
   ```

### uid 无效

```
Error: Element with uid 5 not found
```

**原因：** 页面已导航或刷新，uid 已失效。

**解决：** 重新执行 `take_snapshot` 获取新的 uid。

### 页面加载超时

```
Error: Navigation timeout
```

**解决：**
```bash
navigate_page url="https://slow-site.com" timeout=60000
```

### 无头模式

如果需要无头模式（服务器环境），使用 Chrome 的无头模式：

```bash
/Applications/Google\ Chrome.app/Contents/MacOS/Google\ Chrome \
  --remote-debugging-port=9222 \
  --user-data-dir="$HOME/.chrome-mcp-profile" \
  --headless=new
```

---

## 快速命令参考

| 操作 | 命令 |
|------|------|
| 打开页面 | `navigate_page url="..."` |
| 获取元素 | `take_snapshot` |
| 点击 | `click uid=...` |
| 输入 | `fill uid=... value="..."` |
| 截图 | `screenshot` |
| 执行 JS | `evaluate_script function="..."` |
| 后退/前进 | `navigate_page type="back"` / `type="forward"` |
| 刷新 | `navigate_page type="reload"` |
| 等待 | `wait_for text="..."` |

---

## 相关文档

- [MCP 官方文档](https://modelcontextprotocol.io/)
- [Chrome DevTools Protocol](https://chromedevtools.devtools.protocols.dev/)
- [lingti-bot 浏览器自动化](browser-automation.md)
- [浏览器 AI 操作规则](browser-agent-rules.md)
