# go-musicfox-mcp

基于 [go-musicfox](https://github.com/go-musicfox/go-musicfox) 的 MCP (Model Context Protocol) Server，让 AI Agent 可以直接操作网易云音乐。

## 功能

15 个 MCP Tools：

| Tool | 描述 |
|------|------|
| `search` | 搜索歌曲/专辑/歌手/歌单/歌词 |
| `get_song_url` | 获取歌曲播放链接 |
| `get_lyrics` | 获取歌词（原文/翻译/逐字） |
| `download_song` | 下载歌曲到本地 |
| `get_user_playlists` | 获取用户歌单列表 |
| `get_playlist_songs` | 获取歌单内歌曲 |
| `get_daily_recommend` | 每日推荐 |
| `get_personal_fm` | 私人FM |
| `get_similar_songs` | 相似歌曲 |
| `get_album_detail` | 专辑详情 |
| `get_artist_songs` | 歌手歌曲 |
| `get_user_likes` | 用户喜欢的歌曲 |
| `like_song` / `unlike_song` | 喜欢/取消喜欢 |
| `daily_signin` | 每日签到 |

## 安装

### 从源码编译

```bash
git clone https://github.com/DBeidachazi/go-musicfox-mcp.git
cd go-musicfox-mcp
make build-mcp
cp bin/musicfox-mcp ~/.local/bin/
```

### 直接使用

需要先在 go-musicfox TUI 中登录一次，Cookie 会保存在 `~/.local/share/go-musicfox/cookie`，MCP Server 自动复用。

## 配置

### Hermes Agent

在 `~/.hermes/config.yaml` 中添加：

```yaml
mcp_servers:
  musicfox:
    command: musicfox-mcp
    timeout: 60
```

### 其他 MCP Client

```json
{
  "mcpServers": {
    "musicfox": {
      "command": "musicfox-mcp"
    }
  }
}
```

## 架构

```
cmd/mcp-server/main.go          # MCP server 入口（stdio transport）
internal/mcp/
  service.go                     # 业务逻辑层（认证、API调用、数据转换）
  tools.go                       # MCP Tool 注册与 Handler
```

### 技术要点

- **认证**：复用 musicfox 的 cookie 文件，零配置
- **CDN 修复**：内置 m704→m701 域名回退，解决 403 问题
- **依赖**：全部复用 musicfox internal 包，无重复造轮子
- **传输**：stdio transport，标准 MCP 协议

## 基于

- [go-musicfox](https://github.com/go-musicfox/go-musicfox) - 网易云音乐 TUI 客户端
- [mcp-go](https://github.com/mark3labs/mcp-go) - Go MCP SDK
- [netease-music](https://github.com/go-musicfox/netease-music) - 网易云音乐 API SDK

## License

MIT
