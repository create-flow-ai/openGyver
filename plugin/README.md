# openGyver Claude Code Plugin

This plugin makes openGyver's 49 commands and 180+ subcommands available as native Claude Code tools via MCP.

## Installation

### From Marketplace (recommended)

```bash
# Add the marketplace
/plugin marketplace add https://raw.githubusercontent.com/create-flow-ai/openGyver/main/plugin/marketplace.json

# Install the plugin
/plugin install opengyver
```

### Manual Installation

```bash
# Clone and build
git clone https://github.com/create-flow-ai/openGyver.git
cd openGyver/plugin
npm install  # auto-installs openGyver binary + builds MCP server

# Add to Claude Code MCP config (~/.claude/mcp.json)
{
  "mcpServers": {
    "opengyver": {
      "command": "node",
      "args": ["/path/to/openGyver/plugin/mcp-server/dist/index.js"]
    }
  }
}
```

### Prerequisite: openGyver Binary

The installer auto-downloads from GitHub Releases. You can also install manually:

```bash
# Homebrew
brew tap create-flow-ai/tap
brew install opengyver

# Go
go install github.com/mj/opengyver@latest
```

## Available Tools

Once installed, Claude Code can use these tools:

| Tool | What it does |
|------|-------------|
| `opengyver_encode` | Encode/decode: base64, hex, url, jwt, morse, binary, rot13 |
| `opengyver_hash` | Hash: sha256, md5, bcrypt, hmac, crc32, adler32 |
| `opengyver_convert` | Convert units (length, weight, temp, etc.) and 38 currencies |
| `opengyver_json` | Format, minify, validate, path query JSON |
| `opengyver_generate` | Generate passwords, API keys, UUIDs, OTP secrets |
| `opengyver_color` | Color convert (hex/rgb/hsl/cmyk), WCAG contrast, palettes |
| `opengyver_time` | Timezone convert, epoch, date math, formatting |
| `opengyver_network` | DNS, WHOIS, IP lookup, CIDR calc, URL parse |
| `opengyver_dataformat` | YAML/TOML/XML/CSV to/from JSON |
| `opengyver_validate` | Validate HTML, CSV, XML, YAML, TOML |
| `opengyver_format` | Format/minify HTML, XML, CSS, SQL |
| `opengyver_crypto` | AES encrypt, RSA/SSH keygen, TLS certificates |
| `opengyver_weather` | Weather: current, forecast, historical (no API key) |
| `opengyver_stock` | Stock prices from 35+ global exchanges (no API key) |
| `opengyver_math` | Math expressions, percentages, GCD/LCM, factorial |
| `opengyver_run` | Run any openGyver command directly |

## Examples

Once the plugin is active, Claude Code can:

```
User: "What's the SHA-256 hash of 'hello world'?"
→ Claude uses opengyver_hash(algorithm: "sha256", input: "hello world")

User: "Convert 100 USD to EUR"
→ Claude uses opengyver_convert(value: 100, from: "usd", to: "eur")

User: "What's the weather in Tokyo?"
→ Claude uses opengyver_weather(city: "Tokyo")

User: "Generate a 32-character API key with prefix sk_live_"
→ Claude uses opengyver_generate(type: "apikey", length: 32, prefix: "sk_live_")

User: "Format this SQL: select * from users where id = 1"
→ Claude uses opengyver_format(language: "sql", input: "select * from users where id = 1")
```
