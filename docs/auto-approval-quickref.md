# Auto-Approval Quick Reference (`--yes`)

## Enable Auto-Approval

```bash
# Long form
lingti-bot --yes router --provider deepseek --api-key sk-xxx

# Short form
lingti-bot -y router --provider deepseek --api-key sk-xxx

# With debug
lingti-bot --yes --debug router --provider deepseek --api-key sk-xxx
```

## What It Does

| Operation | Without `--yes` | With `--yes` |
|-----------|-----------------|--------------|
| File write | May ask | ✅ Immediate |
| File delete | May ask | ✅ Immediate |
| Shell command | May ask | ✅ Immediate |
| Dangerous ops | ❌ Blocked | ❌ Blocked |
| Read operations | ✅ Allowed | ✅ Allowed |

## Behavior Examples

### File Operations

```bash
# Without --yes
User: Save this to config.yaml
AI: Should I save the file? (asking...)
User: Yes
AI: ✅ Saved

# With --yes
User: Save this to config.yaml
AI: ✅ Saved config.yaml (247 bytes)
```

### Batch Operations

```bash
# Without --yes
User: Rename all .txt to .md
AI: Found 15 files. Proceed? (asking...)
User: Yes
AI: ✅ Renamed 15 files

# With --yes
User: Rename all .txt to .md
AI: ✅ Renamed 15 files (.txt → .md)
```

## Safety Features (Always Active)

Even with `--yes`, these are still blocked:

- `rm -rf /`
- `mkfs`, `dd if=/dev/zero`
- System-breaking commands
- Critical file deletions

## When to Use

| ✅ Good Use Cases | ❌ Avoid Using |
|-------------------|----------------|
| Development environment | Production servers |
| Batch file processing | Shared systems |
| Code generation | First-time operations |
| Documentation updates | Sensitive data operations |
| Trusted workflows | Learning/testing |

## Best Practices

1. **Use with version control**
   ```bash
   cd /path/to/git/repo
   lingti-bot --yes router [...]
   # Easy to review: git diff
   # Easy to undo: git reset --hard
   ```

2. **Start with debug mode**
   ```bash
   lingti-bot --yes --debug router [...]
   ```

3. **Test in safe directory first**
   ```bash
   mkdir /tmp/test && cd /tmp/test
   lingti-bot --yes router [...]
   ```

4. **Review changes regularly**
   ```bash
   git log --oneline --name-only
   git diff HEAD~5
   ```

## Flag Position (Important!)

```bash
# ✅ Correct - before subcommand
lingti-bot --yes router --provider deepseek --api-key sk-xxx

# ❌ Wrong - after subcommand
lingti-bot gateway --yes --provider deepseek --api-key sk-xxx
```

## Verify It's Working

```bash
# Enable verbose logging
lingti-bot --yes --log verbose router [...]

# Look for this in output:
## 🚀 AUTO-APPROVAL MODE ENABLED
```

## Combine with Other Flags

```bash
# Auto-approval + Debug + Verbose
lingti-bot --yes --debug --log verbose router [...]

# Auto-approval + Custom debug directory
lingti-bot -y --debug-dir /tmp/debug router [...]

# All flags combined
lingti-bot -y --debug --log very-verbose --debug-dir ~/debug router \
  --provider qwen --model qwen-plus --api-key sk-xxx
```

## Aliases for Convenience

```bash
# Add to ~/.bashrc or ~/.zshrc
alias lingti='lingti-bot --yes'
alias lingti-debug='lingti-bot --yes --debug --log verbose'

# Usage
lingti router --provider deepseek --api-key sk-xxx
lingti-debug router --provider deepseek --api-key sk-xxx
```

## Troubleshooting

| Problem | Solution |
|---------|----------|
| AI still asks for confirmation | Check flag position (must be before subcommand) |
| Flag not recognized | Rebuild: `go build -o dist/lingti-bot .` |
| Want to cancel | Press Ctrl+C anytime |
| Review what changed | `git diff` or `git log` |

## Quick Decision Guide

**Should I use `--yes`?**

```
Are you in a git repository? ────────────────┐
                                             │
                                        Yes  │  No
                                             ↓    ↓
Do you trust the operation? ───────────┐    ❌ DON'T USE
                                       │
                                  Yes  │  No
                                       ↓    ↓
Is it production data? ─────────┐     ✅  USE
                                │
                           Yes  │  No
                                ↓    ↓
                               ❌   ✅
```

## Real-World Examples

### Documentation Sync
```bash
lingti-bot --yes router --provider qwen --api-key sk-xxx
User: Sync README_EN.md with README.md
AI: ✅ README_EN.md updated (43 lines changed)
```

### Code Refactoring
```bash
lingti-bot --yes router --provider deepseek --api-key sk-xxx
User: Rename all functions from snake_case to camelCase
AI: ✅ Refactored 28 functions across 7 files
```

### Batch File Cleanup
```bash
lingti-bot --yes router --provider claude --api-key sk-xxx
User: Delete all .log files older than 7 days
AI: ✅ Deleted 142 log files (saved 3.2 GB)
```

## Related Commands

```bash
# Show all global flags
lingti-bot --help

# Show router-specific flags
lingti-bot gateway --help

# Check version
lingti-bot version
```

## See Also

- [Full Documentation](auto-approval.md) - Detailed guide
- [Browser Debug](browser-debug.md) - Debug browser automation
- [README](../README.md) - Main documentation
