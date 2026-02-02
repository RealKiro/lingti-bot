# Feishu (飞书) Integration

This guide explains how to set up Feishu/Lark integration for lingti-bot's message router.

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    lingti-bot Router                     │
├─────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐   │
│  │    Router    │  │    Agent     │  │  MCP Tools   │   │
│  │  (messages)  │──│  (Claude)    │──│  (actions)   │   │
│  └──────────────┘  └──────────────┘  └──────────────┘   │
└─────────────────────────────────────────────────────────┘
          │
          ▼
┌─────────────────┐
│     Feishu      │
│   (WebSocket)   │
└─────────────────┘
```

## Prerequisites

- A Feishu account with permission to create apps (企业自建应用)
- An Anthropic API key for Claude
- lingti-bot binary built and available

## Step 1: Create a Feishu App

1. Go to [Feishu Open Platform](https://open.feishu.cn/app)
2. Click **"创建企业自建应用"** (Create Enterprise App)
3. Fill in app details:
   - App Name: `lingti-bot`
   - App Description: Your bot description
   - App Icon: Upload an icon
4. Click **"创建"** (Create)

## Step 2: Get App Credentials

1. In your app settings, go to **"凭证与基础信息"** (Credentials & Basic Info)
2. Find and copy:
   - **App ID** - Your application's unique identifier
   - **App Secret** - Click "查看" to reveal, then copy

> **Important**: Keep your App Secret secure. Never commit it to version control.

## Step 3: Configure Bot Capabilities

### Enable Bot

1. Go to **"应用能力"** → **"机器人"** (Capabilities → Bot)
2. Toggle **"启用机器人"** (Enable Bot) to ON
3. Configure bot settings:
   - Bot Name: `lingti-bot`
   - Bot Description: Your bot description

### Enable WebSocket Mode

1. In Bot settings, find **"消息接收方式"** (Message Receiving Method)
2. Select **"使用长连接接收消息"** (Use WebSocket to receive messages)

## Step 4: Configure Permissions

1. Go to **"权限管理"** (Permission Management)
2. Add the following scopes:

| Permission | Scope ID | Description |
|------------|----------|-------------|
| 获取与发送单聊、群组消息 | `im:message` | Send and receive messages |
| 获取用户基本信息 | `contact:user.base:readonly` | Read user info for usernames |
| 获取群组信息 | `im:chat:readonly` | Read chat info |

3. Click **"批量开通"** (Batch Enable)

## Step 5: Subscribe to Events

1. Go to **"事件订阅"** (Event Subscriptions)
2. Add the following events:

| Event | Event ID | Description |
|-------|----------|-------------|
| 接收消息 | `im.message.receive_v1` | Receive messages sent to bot |

3. Ensure **"使用长连接接收事件"** (Use WebSocket for events) is enabled

## Step 6: Publish the App

1. Go to **"版本管理与发布"** (Version Management)
2. Click **"创建版本"** (Create Version)
3. Fill in version details
4. Click **"申请发布"** (Request Publish)
5. If you're the admin, approve the publish request

> For development testing, you can use the app in "开发中" (Development) status within your own account.

## Step 7: Run lingti-bot Router

### Using Environment Variables

```bash
export FEISHU_APP_ID="cli_your_app_id"
export FEISHU_APP_SECRET="your_app_secret"
export ANTHROPIC_API_KEY="sk-ant-your-api-key"
# Optional: custom API base URL (for proxies or alternative endpoints)
# export ANTHROPIC_BASE_URL="https://your-proxy.com/v1"

lingti-bot router
```

### Using Command-Line Flags

```bash
lingti-bot router \
  --feishu-app-id "cli_your_app_id" \
  --feishu-app-secret "your_app_secret" \
  --api-key "sk-ant-your-api-key"
  # --base-url "https://your-proxy.com/v1"  # optional
```

### Using a .env File

Create a `.env` file:

```bash
FEISHU_APP_ID=cli_your_app_id
FEISHU_APP_SECRET=your_app_secret
ANTHROPIC_API_KEY=sk-ant-your-api-key
# ANTHROPIC_BASE_URL=https://your-proxy.com/v1  # optional
```

Then run:

```bash
source .env && lingti-bot router
```

### Running Both Slack and Feishu

You can run multiple platforms simultaneously:

```bash
export SLACK_BOT_TOKEN="xoxb-..."
export SLACK_APP_TOKEN="xapp-..."
export FEISHU_APP_ID="cli_..."
export FEISHU_APP_SECRET="..."
export ANTHROPIC_API_KEY="sk-ant-..."

lingti-bot router
```

## Step 8: Test the Integration

1. Open Feishu app
2. Find your bot:
   - Search for the bot name in the search bar
   - Or go to **"工作台"** → find your bot app
3. Start a conversation:
   - **Direct message**: Just send a message to the bot
   - **Group chat**: Add the bot to a group, then @mention it

Example messages:
- `what's on my calendar today?`
- `@lingti-bot list files in ~/Desktop`
- `@lingti-bot what's my system info?`

## Available Commands

Once connected, the bot can:

| Category | Examples |
|----------|----------|
| **Calendar** | "What's on my calendar today?", "Schedule a meeting tomorrow at 2pm" |
| **Files** | "List files in ~/Downloads", "Find old files on my Desktop" |
| **System** | "What's my CPU usage?", "Show disk space" |
| **Shell** | "Run `ls -la`", "Check git status" |
| **Process** | "List running processes", "What's using the most memory?" |

## Troubleshooting

### Bot not responding

1. Check that App ID and App Secret are set correctly
2. Verify the bot is running: check router logs
3. Ensure WebSocket mode is enabled in Feishu app settings
4. Check that the app has been published or you're testing with the correct account

### "failed to verify credentials" error

Your `FEISHU_APP_ID` or `FEISHU_APP_SECRET` is invalid. Double-check the credentials in the Feishu Open Platform.

### Bot only works in DMs, not in groups

Make sure:
1. The bot has been added to the group
2. You @mention the bot in group messages
3. The `im:message` permission is enabled

### Messages not being received

1. Verify **"使用长连接接收消息"** is enabled
2. Check that `im.message.receive_v1` event is subscribed
3. Ensure the app version with bot capabilities is published

### Permission denied errors

1. Go to Permission Management in Feishu Open Platform
2. Ensure all required permissions are enabled
3. If permissions were recently added, you may need to create a new app version

## Message Format Notes

- **DMs**: Bot responds to all direct messages
- **Group chats**: Bot only responds when @mentioned
- **@mentions**: The `@botname` is automatically removed from the message text before processing

## Security Considerations

- Never commit App Secret to version control
- Use environment variables or a secrets manager
- Restrict app installation to trusted organizations
- Review bot permissions regularly
- Consider IP whitelist for production deployments

## Running as a Service

### macOS (launchd)

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.lingti.bot.router</string>
    <key>ProgramArguments</key>
    <array>
        <string>/usr/local/bin/lingti-bot</string>
        <string>router</string>
    </array>
    <key>EnvironmentVariables</key>
    <dict>
        <key>FEISHU_APP_ID</key>
        <string>cli_...</string>
        <key>FEISHU_APP_SECRET</key>
        <string>...</string>
        <key>ANTHROPIC_API_KEY</key>
        <string>sk-ant-...</string>
    </dict>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>/tmp/lingti-bot-router.log</string>
    <key>StandardErrorPath</key>
    <string>/tmp/lingti-bot-router.log</string>
</dict>
</plist>
```

### Linux (systemd)

```ini
[Unit]
Description=Lingti Bot Router
After=network.target

[Service]
Type=simple
Environment=FEISHU_APP_ID=cli_...
Environment=FEISHU_APP_SECRET=...
Environment=ANTHROPIC_API_KEY=sk-ant-...
ExecStart=/usr/local/bin/lingti-bot router
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

## Lark (International) vs Feishu (China)

This integration works with both:
- **Feishu (飞书)** - China version at `open.feishu.cn`
- **Lark** - International version at `open.larksuite.com`

The SDK automatically handles the regional differences. Use the appropriate developer console for your version.

## References

- [Feishu Open Platform](https://open.feishu.cn/)
- [Lark Developer Documentation](https://open.larksuite.com/document)
- [Bot Development Guide](https://open.feishu.cn/document/home/develop-a-bot-in-5-minutes/create-an-app)
- [Event Subscription Guide](https://open.feishu.cn/document/ukTMukTMukTM/uUTNz4SN1MjL1UzM)
- [Lark SDK for Go](https://github.com/larksuite/oapi-sdk-go)
