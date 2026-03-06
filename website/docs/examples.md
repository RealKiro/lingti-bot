---
id: examples
title: 使用示例
sidebar_position: 3
---

# 使用示例

## 内置 Web 聊天界面 {#web-chat-ui}

无需微信、无需 Slack，一个参数即可在本地浏览器中打开完整聊天界面：

```bash
lingti-bot gateway --provider deepseek --api-key sk-xxx --webapp-port 8080
# 打开浏览器访问 http://localhost:8080
```

<p align="center">
<img src="/lingti-bot/images/webapp-demo.png" alt="Web 聊天界面演示" width="800" />
</p>

**核心特性：**

| 特性 | 说明 |
|------|------|
| **多会话并行** | 左侧边栏管理多个独立会话，每个会话有独立的 AI 记忆上下文，互不干扰 |
| **并发处理** | 在会话 A 等待长任务时，立即切换到会话 B 发送新消息——两者并行处理 |
| **会话持久化** | 所有会话列表和聊天记录保存在浏览器 `localStorage`，刷新页面不丢失 |
| **Markdown 渲染** | AI 回复支持完整 Markdown：代码块、表格、列表、加粗等 |
| **端口自动切换** | 若指定端口被占用，自动递增尝试直到找到空闲端口 |
| **无需额外安装** | 纯 HTML + 原生 JS，内嵌于二进制，无 Node.js / npm 依赖 |

**多会话并行示意：**

```
浏览器
├── 会话 A（UUID-1）── AI 正在处理长任务...
├── 会话 B（UUID-2）── 刚发了新消息，AI 立即响应
└── 会话 C（UUID-3）── 历史记录，随时可查
```

每个会话对应一个独立 `channelID`，AI 记忆完全隔离——在会话 A 说的内容不会影响会话 B 的上下文。

**通过配置文件启用：**

```yaml
# ~/.lingti.yaml
platforms:
  webapp:
    port: 8080
```

**环境变量：**

```bash
export WEBAPP_PORT=8080
lingti-bot gateway --provider deepseek --api-key sk-xxx
```

---

## 智能对话、文件管理、信息检索

<table>
<tr>
<td width="33%"><img src="/lingti-bot/images/demo-chat-1.png" alt="智能助手" /></td>
<td width="33%"><img src="/lingti-bot/images/demo-chat-2.png" alt="企业微信文件传输" /></td>
<td width="33%"><img src="/lingti-bot/images/demo-chat-3.png" alt="信息搜索" /></td>
</tr>
<tr>
<td align="center"><sub>💬 智能对话</sub></td>
<td align="center"><sub>📁 企微文件传输</sub></td>
<td align="center"><sub>🔍 信息搜索</sub></td>
</tr>
</table>

<img src="/lingti-bot/images/demo-terminal.png" alt="Terminal Demo" />
<p><sub>克隆代码后直接编译运行，配合 DeepSeek 模型，实时处理钉钉消息</sub></p>

---

## 定时任务 — AI 自动创建 Cron Job {#cron-jobs}

> 用自然语言创建定时任务 — 告诉 AI 你想要什么，剩下的交给它

<p align="center">
<img src="/lingti-bot/images/demo-cron-wecom.png" alt="AI 创建定时任务演示" width="720" />
</p>

在企业微信中对 AI 说一句话，即可创建复杂的定时任务。支持 Cron 表达式调度、macOS 系统通知、Shell 脚本执行等，真正实现无人值守自动化。

---

## 企业微信 AI 文件助手 {#wecom-file}

> 用自然语言管理和传输文件 — 就像跟同事说话一样简单

<p align="center">
<img src="/lingti-bot/images/wecom-file-transfer.png" alt="企业微信 AI 文件传输" width="720" />
</p>

直接在企业微信中用自然语言浏览、查找、传输电脑上的文件。无需远程桌面，无需 U 盘，对 AI 说一句话即可。

---

## 社交平台自动化 {#social-automation}

> 通过 MCP + Chrome 浏览器能力，让 AI 代替你操作知乎、小红书等内容平台

lingti-bot 将浏览器自动化能力与 AI 深度结合，针对中国主流内容平台提供**智能化运营**支持。不是简单的脚本录制回放，而是 AI 理解页面语义后的精准操作。

**已支持平台：**

| 平台 | 发帖/回答 | 评论 | 点赞/收藏 | 状态 |
|------|:---------:|:----:|:---------:|------|
| **知乎** | ✅ | ✅ | ✅ | 已支持 |
| **小红书** | 🔜 | 🔜 | 🔜 | 规划中 |
| **微博** | 🔜 | 🔜 | 🔜 | 规划中 |
| **抖音（网页版）** | 🔜 | 🔜 | 🔜 | 规划中 |
| **B站** | 🔜 | 🔜 | 🔜 | 规划中 |
| **今日头条** | 🔜 | 🔜 | 🔜 | 规划中 |

**使用示例：**

```
"帮我在知乎上回答这个问题，内容围绕 Go 语言的优势"
"打开知乎这篇文章，帮我写一条评论"
"在小红书发一篇关于效率工具的笔记"
"帮我给这篇知乎回答点个赞"
```

**工作原理：** AI 通过 MCP 协议调用浏览器工具（CDP 或 Chrome MCP），在你已登录的 Chrome 中操作，**复用你的登录态**，无需提供账号密码。AI 读取页面无障碍树快照，理解每个按钮和输入框的含义，像人一样精准操作。

> 📖 **详细文档：[社交平台自动化指南](./social-platform-automation)**

---

## 浏览器自动化 {#browser-automation}

> 无需 Puppeteer、无需 Node.js、无需安装任何依赖——直接用自然语言驾驭 Chrome

lingti-bot 内置完整的 CDP（Chrome DevTools Protocol）自动化引擎，采用**快照-操作（Snapshot-then-Act）**模式：先读取页面无障碍树，理解每个元素的角色和名称，再精准点击、输入、截图。

**直接在你正在使用的 Chrome 里操作：**

```bash
# 1. 用调试端口启动 Chrome（一次配置，长期使用）
/Applications/Google\ Chrome.app/Contents/MacOS/Google\ Chrome \
  --remote-debugging-port=9222 \
  --user-data-dir="$HOME/.lingti-bot/chrome-profile"

# 2. ~/.lingti.yaml 中添加一行
# browser:
#   cdp_url: "127.0.0.1:9222"
```

之后，对 AI 说话即可——bot 会在你已有的 Chrome 里开新标签页操作，不再另弹窗口：

```
"打开知乎，搜索 Go 语言并截图"
"帮我登录 GitHub，查看我最新的 PR"
"打开这个表单，填写姓名张三、手机 138xxxx1234，然后提交"
"把电商后台所有待审核商品都批量通过"
"每隔5分钟截一张这个监控大屏的图"
```

**工作原理（AI 自动完成全部步骤）：**

```
browser_navigate url="https://www.zhihu.com"
browser_snapshot
→ [1] textbox "搜索" [2] button "搜索" [3] link "登录" ...

browser_type ref=1 text="Go 语言" submit=true
browser_snapshot
→ 看到搜索结果列表

browser_screenshot path="/tmp/zhihu-result.png"
→ 截图保存完成
```

14 个浏览器工具覆盖完整自动化流程：生命周期管理、导航、无障碍树快照、点击/输入/按键、JavaScript 执行、批量操作、多标签页管理、截图。

> 📖 **完整使用指南、连接配置、故障排除：[浏览器自动化文档](./browser-automation)**
>
> 📋 **AI 操作规则详解（快照法则、搜索行为、弹窗处理、批量操作等）：[浏览器 AI 操作规则](./browser-agent-rules)**

---

## 日常任务示例

配置完成后，你可以让 AI 助手执行以下操作：

### 日历与日程

```
"今天有什么日程安排？"
"这周有哪些会议？"
"帮我创建一个明天下午3点的会议，标题是'产品评审'"
"搜索所有包含'周报'的日程"
```

### 文件操作与传输

```
"列出桌面上的所有文件"
"读取 ~/Documents/notes.txt 的内容"
"将 ~/Desktop/报告.pdf 发送给我"
"把 Documents 里的产品介绍发给我"
"桌面上超过30天没动过的文件有哪些？"
"帮我把这些旧文件移到废纸篓"
```

### 系统与进程

```
"我的电脑配置是什么？"
"现在 CPU 占用多少？"
"Chrome 占用了多少内存？"
"结束 PID 1234 的进程"
```

### 网络与搜索

```
"我的 IP 地址是什么？"
"帮我搜索一下最新的 AI 新闻"
"查询 github.com 的 DNS"
```

### 音乐控制

```
"播放音乐"
"下一首"
"音量调到 50%"
"播放周杰伦的歌"
```

### 组合任务

```
"查看今天的日程，然后检查天气，最后列出待办事项"
"帮我整理桌面：列出超过60天的旧文件，然后移到废纸篓"
"搜索最近的科技新闻，整理成备忘录"
```
