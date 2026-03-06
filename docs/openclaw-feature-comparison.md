# OpenClaw vs lingti-bot: Technical Feature Comparison

> Technical comparison of capabilities between OpenClaw and lingti-bot.

## Overview

| Metric | OpenClaw | lingti-bot |
|--------|----------|------------|
| Primary Language | TypeScript | Go |
| TypeScript Files | 2,577 | 0 |
| Go Files | 43+ | ~62 |
| Documentation Files | 615 | 10 |
| Platform Extensions | 31 | 7 |
| Skills | 54 | 0 |

## Feature Comparison Matrix

| Feature | OpenClaw | lingti-bot |
|---------|----------|------------|
| CLI Tool | ✅ | ✅ |
| MCP Server | ✅ | ✅ |
| Web UI | ✅ | ✅ Built-in webapp |
| macOS App | ✅ | ❌ |
| iOS App | ✅ | ❌ |
| Android App | ✅ | ❌ |
| Browser Automation | ✅ | ✅ CDP engine (14 tools) |
| Skills System | ✅ | ✅ 8 built-in skills |
| Plugin System | ✅ | ❌ |
| Hooks System | ✅ | ❌ |
| Memory/RAG | ✅ | ❌ (session memory only) |
| Terminal UI (TUI) | ✅ | ❌ |
| Cloud Relay | ❌ | ✅ |
| Single Binary | ❌ | ✅ |
| Agent Management | ✅ | ✅ `agents add/list/bind` |
| Channel Management | ✅ | ✅ `channels add/list` |
| Multi-agent routing | ✅ | ✅ bindings by platform/channel |
| Cron/Scheduling | ✅ | ✅ AI-powered cron jobs |
| Social Platform Automation | ❌ | ✅ 知乎, 小红书, etc. |
| China Platforms | ❌ | ✅ 飞书, 企微, 钉钉, 微信 |

## Messaging Platforms

### lingti-bot Supported (19+)
- Discord, Telegram, Slack
- Feishu/Lark, DingTalk, WeCom (企业微信), WeChat (微信公众号)
- WhatsApp, LINE, Microsoft Teams, Matrix/Element
- Google Chat, Mattermost, iMessage, Signal
- Twitch, NOSTR, Zalo, Nextcloud Talk
- Built-in Web Chat UI (`--webapp-port`)

### OpenClaw Additional Platforms (24+)
- iMessage / BlueBubbles
- Signal
- WhatsApp
- LINE
- Matrix / Element
- Microsoft Teams
- Google Chat
- Mattermost
- Nextcloud Talk
- Twitch
- NOSTR
- Zalo
- And more...

## Native Applications

### OpenClaw
```
apps/
├── macos/          # SwiftUI macOS app
├── ios/            # SwiftUI iOS app (Chat, Voice, Location, Camera)
├── android/        # Gradle Android app
└── shared/         # OpenClawKit shared framework
```

### lingti-bot
Built-in web chat UI, enabled with `--webapp-port 8080`:
- Multi-session parallel chat
- Session persistence (localStorage)
- Markdown rendering
- No Node.js / no extra install — embedded in binary - CLI only

## Web UI

### OpenClaw (`ui/`)
- 113+ TypeScript/TSX files
- Vite-based build
- Features:
  - Chat interface
  - Configuration dashboard
  - Device management
  - Session management
  - Skills display
  - Theme management
  - Debug tools

### lingti-bot (`internal/browser/` + 14 MCP tools)
- Chrome DevTools Protocol (CDP)
- Snapshot-then-Act model (accessibility tree)
- Connects to existing Chrome (`--remote-debugging-port`)
- 14 tools: navigate, snapshot, click, type, press, JS execution, batch click, tabs, screenshot
- No Node.js / Playwright required

## Skills System

### OpenClaw (`skills/` - 54 total)
| Category | Skills |
|----------|--------|
| Productivity | Apple Notes, Apple Reminders, Bear Notes, Notion, Obsidian, Trello |
| Communication | iMessage, Signal, Slack, Discord, WhatsApp |
| Development | GitHub, Coding Agent |
| AI Models | Gemini, OpenAI Image Gen, OpenAI Whisper |
| Media | Spotify, Music Control |
| System | 1Password, tmux, Weather |
| Location | Local Places, GoPlaces |

### lingti-bot (`skills/` - 8 built-in)
| Skill | Description |
|-------|-------------|
| GitHub | PR/Issue management |
| Slack | Slack messaging |
| Discord | Discord messaging |
| Peekaboo | macOS UI automation |
| Tmux | Terminal session control |
| Weather | Weather queries |
| 1Password | Password management |
| Obsidian | Note-taking |

Skills are YAML-frontmatter Markdown files (`SKILL.md`). Custom and project-level skills supported.

## Browser Automation

### OpenClaw (`src/browser/`)
- Chrome DevTools Protocol (CDP)
- Playwright integration
- Features:
  - Page navigation
  - Element interaction
  - Screenshot capture
  - Profile management
  - State observation

### lingti-bot
None

## Media Understanding

### OpenClaw (`src/media-understanding/`)
- Image processing
- Video processing
- Audio transcription
- Vision model integration
- Frame extraction
- SSRF protection

### lingti-bot
Basic screenshot tool only

## Memory & Knowledge

### OpenClaw
- `src/memory/` - Core memory system
- `extensions/memory-lancedb/` - Vector database
- `src/hooks/bundled/session-memory/` - Session memory
- Persistent knowledge across sessions

### lingti-bot
In-memory session storage only (lost on restart)

## Hooks System

### OpenClaw (`src/hooks/`)
- Custom hook loading
- Internal hooks
- Plugin hooks
- Gmail integration
- Hook installation system

### lingti-bot
None

## Plugin System

### OpenClaw
- `src/plugin-sdk/` - Plugin SDK
- `src/plugins/` - Plugin runtime
- Package-based architecture
- Plugin validation

### lingti-bot
None

## Terminal UI

### OpenClaw (`src/tui/`)
- Interactive terminal interface
- Gateway chat
- Event handlers
- Theme management
- Input history
- Overlays

### lingti-bot
None - stdout only

## Configuration

### OpenClaw (`src/config/` - 126+ files)
- Channel capabilities
- Agent concurrency
- Identity/avatar management
- Multi-agent configuration
- Backup rotation
- Broadcast settings
- Platform-specific configs
- Legacy config detection
- Nix integration

### lingti-bot (`internal/config/config.go`)
- Named provider definitions (`providers:`)
- Agent definitions with per-agent model/instructions/workspace (`agents:`)
- Routing bindings by platform+channel (`bindings:`)
- Platform credentials (`platforms:`)
- Security sandboxing (`security:`)
- Browser config, logging, MCP servers

## CLI Commands

### OpenClaw (`src/commands/` - 187 files)
- Agent management (add, delete, identity, list)
- Channel management
- Auth choice handling
- Config management
- Plugin management
- Many more specialized commands

### lingti-bot (`cmd/` - 20+ files)
- `serve`, `relay`, `gateway`
- `agents add/list/info/bind/unbind` — Agent management
- `channels add/list` — Channel credential management
- `skills`, `skills check`, `skills info` — Skills discovery
- `doctor` — Health diagnostics
- `onboard` — Interactive first-run setup
- `version`, `whoami`

## Infrastructure

### OpenClaw
- Multiple Dockerfiles (dev, sandbox, sandbox-browser)
- Docker Compose
- Fly.io deployment (fly.toml)
- SystemD service files
- Launchd configuration
- E2E tests, live model tests

### lingti-bot
- Makefile
- Basic build scripts

## Advanced Features in OpenClaw Only

| Feature | Location | Description |
|---------|----------|-------------|
| Auto-reply | `src/auto-reply/` | Automated response system |
| Cron/Scheduling | `src/cron/` | Scheduled tasks |
| Daemon Mode | `src/daemon/` | Background service |
| Canvas Host | `src/canvas-host/` | Canvas rendering |
| Device Pairing | `src/pairing/` | Multi-device sync |
| Media Pipeline | `src/media/` | Media processing |
| Security Module | `src/security/` | Security features |
| Wizard | `src/wizard/` | Onboarding flow |

## MCP Tools Comparison

Both projects share the same MCP tool implementations:

| Category | Tools |
|----------|-------|
| Filesystem | file_read, file_write, file_list, file_search, file_info |
| Shell | shell_execute, shell_which |
| System | system_info, disk_usage, env_get, env_list |
| Process | process_list, process_info, process_kill |
| Network | network_interfaces, network_connections, network_ping, network_dns_lookup |
| Calendar | calendar_list_events, calendar_create_event, calendar_list_calendars, calendar_today, calendar_search, calendar_delete_event |
| File Manager | file_list_old, file_delete_old, file_delete_list, file_trash |

Additional tools in both (via agent):
- Clipboard operations
- Notes (macOS)
- Reminders (macOS)
- Screenshot
- Notifications
- Music control (macOS)
- GitHub operations
- Web search/fetch
- Weather

## Summary

**OpenClaw** is a comprehensive personal AI assistant platform with:
- Native apps for all platforms
- Full-featured web UI
- 54 specialized skills
- Browser automation
- Memory/RAG systems
- Plugin architecture
- Extensive infrastructure

**lingti-bot** is a focused CLI messaging bot with:
- Single Go binary, zero runtime dependencies
- Cloud relay for easy setup (no public server needed)
- 19+ platform support including all China platforms
- Agent + Channel management (`agents add`, `channels add`)
- Multi-agent routing via bindings
- Built-in browser automation (CDP, 14 tools)
- Built-in web chat UI
- 75+ MCP tools + 8 skills
- AI-powered cron jobs
- Social platform automation (知乎, 小红书)

## Related Documents

- [OpenClaw Reference](openclaw-reference.md)
- [vs OpenClaw Integration](vs-openclaw-integration.md)
- [Roadmap](roadmap.md)
