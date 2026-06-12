package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterTools registers all MCP tools with the server.
func (s *Service) RegisterTools(mcpServer *server.MCPServer) {
	mcpServer.AddTool(
		mcp.Tool{
			Name:        "get_current_user",
			Description: "获取当前登录用户的 user_id 和昵称",
			InputSchema: mcp.ToolInputSchema{Type: "object", Properties: map[string]interface{}{}},
		},
		s.handleGetCurrentUser,
	)

	mcpServer.AddTool(
		mcp.Tool{
			Name:        "search",
			Description: "搜索网易云音乐（歌曲/专辑/歌手/歌单/歌词/电台）",
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"keywords": map[string]interface{}{"type": "string", "description": "搜索关键词"},
					"type":     map[string]interface{}{"type": "string", "description": "搜索类型", "enum": []string{"songs", "albums", "artists", "playlists", "lyrics", "djradio"}, "default": "songs"},
					"limit":    map[string]interface{}{"type": "integer", "description": "最大结果数（默认20）", "default": 20},
					"offset":   map[string]interface{}{"type": "integer", "description": "偏移量（默认0）", "default": 0},
				},
				Required: []string{"keywords"},
			},
		},
		s.handleSearch,
	)

	mcpServer.AddTool(
		mcp.Tool{
			Name:        "search_suggest",
			Description: "搜索建议/自动补全",
			InputSchema: mcp.ToolInputSchema{
				Type:       "object",
				Properties: map[string]interface{}{"keywords": map[string]interface{}{"type": "string", "description": "搜索关键词"}},
				Required:   []string{"keywords"},
			},
		},
		s.handleSearchSuggest,
	)

	mcpServer.AddTool(
		mcp.Tool{
			Name:        "get_song_url",
			Description: "获取歌曲播放链接（临时，会过期）",
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"song_id": map[string]interface{}{"type": "integer", "description": "歌曲ID"},
					"quality": map[string]interface{}{"type": "string", "description": "音质", "enum": []string{"standard", "higher", "exhigh", "lossless", "hires"}, "default": "exhigh"},
				},
				Required: []string{"song_id"},
			},
		},
		s.handleGetSongURL,
	)

	mcpServer.AddTool(
		mcp.Tool{
			Name:        "get_song_detail",
			Description: "获取歌曲详情（封面、专辑、时长等元信息）",
			InputSchema: mcp.ToolInputSchema{
				Type:       "object",
				Properties: map[string]interface{}{"song_id": map[string]interface{}{"type": "integer", "description": "歌曲ID"}},
				Required:   []string{"song_id"},
			},
		},
		s.handleGetSongDetail,
	)

	mcpServer.AddTool(
		mcp.Tool{
			Name:        "get_lyrics",
			Description: "获取歌词（原文、翻译、逐字）",
			InputSchema: mcp.ToolInputSchema{
				Type:       "object",
				Properties: map[string]interface{}{"song_id": map[string]interface{}{"type": "integer", "description": "歌曲ID"}},
				Required:   []string{"song_id"},
			},
		},
		s.handleGetLyrics,
	)

	mcpServer.AddTool(
		mcp.Tool{
			Name:        "download_song",
			Description: "下载歌曲到本地文件",
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"song_id":    map[string]interface{}{"type": "integer", "description": "歌曲ID"},
					"quality":    map[string]interface{}{"type": "string", "description": "音质", "enum": []string{"standard", "higher", "exhigh", "lossless"}, "default": "exhigh"},
					"output_dir": map[string]interface{}{"type": "string", "description": "输出目录（可选）"},
				},
				Required: []string{"song_id"},
			},
		},
		s.handleDownloadSong,
	)

	mcpServer.AddTool(
		mcp.Tool{
			Name:        "get_user_playlists",
			Description: "获取用户歌单列表。不传 user_id 则获取当前登录用户的歌单",
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"user_id": map[string]interface{}{"type": "integer", "description": "用户ID（可选，默认当前登录用户）"},
					"limit":   map[string]interface{}{"type": "integer", "description": "最大结果数（默认30）", "default": 30},
					"offset":  map[string]interface{}{"type": "integer", "description": "偏移量（默认0）", "default": 0},
				},
				Required: []string{},
			},
		},
		s.handleGetUserPlaylists,
	)

	mcpServer.AddTool(
		mcp.Tool{
			Name:        "get_playlist_detail",
			Description: "获取歌单详情（封面、描述、标签、歌曲列表）",
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"playlist_id":    map[string]interface{}{"type": "integer", "description": "歌单ID"},
					"include_songs":  map[string]interface{}{"type": "boolean", "description": "是否包含歌曲列表（默认true）", "default": true},
				},
				Required: []string{"playlist_id"},
			},
		},
		s.handleGetPlaylistDetail,
	)

	mcpServer.AddTool(
		mcp.Tool{
			Name:        "get_playlist_songs",
			Description: "获取歌单内的歌曲列表",
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"playlist_id": map[string]interface{}{"type": "integer", "description": "歌单ID"},
					"get_all":     map[string]interface{}{"type": "boolean", "description": "获取全部歌曲（默认false，最多1000首）", "default": false},
				},
				Required: []string{"playlist_id"},
			},
		},
		s.handleGetPlaylistSongs,
	)

	mcpServer.AddTool(
		mcp.Tool{
			Name:        "get_daily_recommend",
			Description: "获取每日推荐歌曲（需要登录）",
			InputSchema: mcp.ToolInputSchema{Type: "object", Properties: map[string]interface{}{}},
		},
		s.handleGetDailyRecommend,
	)

	mcpServer.AddTool(
		mcp.Tool{
			Name:        "get_personal_fm",
			Description: "获取私人FM歌曲（需要登录）",
			InputSchema: mcp.ToolInputSchema{Type: "object", Properties: map[string]interface{}{}},
		},
		s.handleGetPersonalFM,
	)

	mcpServer.AddTool(
		mcp.Tool{
			Name:        "get_similar_songs",
			Description: "获取相似歌曲推荐",
			InputSchema: mcp.ToolInputSchema{
				Type:       "object",
				Properties: map[string]interface{}{"song_id": map[string]interface{}{"type": "integer", "description": "歌曲ID"}},
				Required:   []string{"song_id"},
			},
		},
		s.handleGetSimilarSongs,
	)

	mcpServer.AddTool(
		mcp.Tool{
			Name:        "like_song",
			Description: "喜欢歌曲",
			InputSchema: mcp.ToolInputSchema{
				Type:       "object",
				Properties: map[string]interface{}{"song_id": map[string]interface{}{"type": "integer", "description": "歌曲ID"}},
				Required:   []string{"song_id"},
			},
		},
		s.handleLikeSong,
	)

	mcpServer.AddTool(
		mcp.Tool{
			Name:        "unlike_song",
			Description: "取消喜欢歌曲",
			InputSchema: mcp.ToolInputSchema{
				Type:       "object",
				Properties: map[string]interface{}{"song_id": map[string]interface{}{"type": "integer", "description": "歌曲ID"}},
				Required:   []string{"song_id"},
			},
		},
		s.handleUnlikeSong,
	)

	mcpServer.AddTool(
		mcp.Tool{
			Name:        "daily_signin",
			Description: "每日签到",
			InputSchema: mcp.ToolInputSchema{Type: "object", Properties: map[string]interface{}{}},
		},
		s.handleDailySignin,
	)

	mcpServer.AddTool(
		mcp.Tool{
			Name:        "get_album_detail",
			Description: "获取专辑详情（封面、歌手、描述）",
			InputSchema: mcp.ToolInputSchema{
				Type:       "object",
				Properties: map[string]interface{}{"album_id": map[string]interface{}{"type": "integer", "description": "专辑ID"}},
				Required:   []string{"album_id"},
			},
		},
		s.handleGetAlbumDetail,
	)

	mcpServer.AddTool(
		mcp.Tool{
			Name:        "get_artist_songs",
			Description: "获取歌手的歌曲列表",
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"artist_id": map[string]interface{}{"type": "integer", "description": "歌手ID"},
					"limit":     map[string]interface{}{"type": "integer", "description": "最大结果数（默认30）", "default": 30},
					"offset":    map[string]interface{}{"type": "integer", "description": "偏移量（默认0）", "default": 0},
				},
				Required: []string{"artist_id"},
			},
		},
		s.handleGetArtistSongs,
	)

	mcpServer.AddTool(
		mcp.Tool{
			Name:        "get_user_likes",
			Description: "获取用户喜欢的歌曲。不传 user_id 则获取当前登录用户喜欢的歌曲",
			InputSchema: mcp.ToolInputSchema{
				Type:       "object",
				Properties: map[string]interface{}{"user_id": map[string]interface{}{"type": "integer", "description": "用户ID（可选，默认当前登录用户）"}},
				Required:   []string{},
			},
		},
		s.handleGetUserLikes,
	)

	mcpServer.AddTool(
		mcp.Tool{
			Name:        "create_playlist",
			Description: "创建歌单",
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"name":    map[string]interface{}{"type": "string", "description": "歌单名称"},
					"privacy": map[string]interface{}{"type": "boolean", "description": "是否私密（默认false）", "default": false},
				},
				Required: []string{"name"},
			},
		},
		s.handleCreatePlaylist,
	)

	mcpServer.AddTool(
		mcp.Tool{
			Name:        "delete_playlist",
			Description: "删除歌单",
			InputSchema: mcp.ToolInputSchema{
				Type:       "object",
				Properties: map[string]interface{}{"playlist_id": map[string]interface{}{"type": "integer", "description": "歌单ID"}},
				Required:   []string{"playlist_id"},
			},
		},
		s.handleDeletePlaylist,
	)

	mcpServer.AddTool(
		mcp.Tool{
			Name:        "rename_playlist",
			Description: "重命名歌单",
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"playlist_id": map[string]interface{}{"type": "integer", "description": "歌单ID"},
					"name":        map[string]interface{}{"type": "string", "description": "新名称"},
				},
				Required: []string{"playlist_id", "name"},
			},
		},
		s.handleRenamePlaylist,
	)

	mcpServer.AddTool(
		mcp.Tool{
			Name:        "add_to_playlist",
			Description: "向歌单添加歌曲",
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"playlist_id": map[string]interface{}{"type": "integer", "description": "歌单ID"},
					"song_ids":    map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "integer"}, "description": "歌曲ID列表"},
				},
				Required: []string{"playlist_id", "song_ids"},
			},
		},
		s.handleAddToPlaylist,
	)

	mcpServer.AddTool(
		mcp.Tool{
			Name:        "remove_from_playlist",
			Description: "从歌单删除歌曲",
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"playlist_id": map[string]interface{}{"type": "integer", "description": "歌单ID"},
					"song_ids":    map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "integer"}, "description": "歌曲ID列表"},
				},
				Required: []string{"playlist_id", "song_ids"},
			},
		},
		s.handleRemoveFromPlaylist,
	)
}

// ========== Tool Handlers ==========

func (s *Service) handleGetCurrentUser(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	info := s.GetCurrentUserInfo()
	return toolJSON(info)
}

func (s *Service) handleSearch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	keywords, _ := args["keywords"].(string)
	if keywords == "" {
		return toolError("keywords is required"), nil
	}
	searchType := SearchSongs
	if t, ok := args["type"].(string); ok {
		searchType = SearchType(t)
	}
	limit, offset := 20, 0
	if l, ok := args["limit"]; ok {
		limit = toInt(l)
	}
	if o, ok := args["offset"]; ok {
		offset = toInt(o)
	}
	result, err := s.Search(ctx, keywords, searchType, limit, offset)
	if err != nil {
		return toolError(err.Error()), nil
	}
	return toolJSON(result)
}

func (s *Service) handleSearchSuggest(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	keywords, _ := args["keywords"].(string)
	if keywords == "" {
		return toolError("keywords is required"), nil
	}
	result, err := s.SearchSuggest(ctx, keywords)
	if err != nil {
		return toolError(err.Error()), nil
	}
	return toolJSON(result)
}

func (s *Service) handleGetSongURL(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	songID := toInt64(args["song_id"])
	if songID == 0 {
		return toolError("song_id is required"), nil
	}
	quality := "exhigh"
	if q, ok := args["quality"].(string); ok {
		quality = q
	}
	result, err := s.GetSongURL(ctx, songID, quality)
	if err != nil {
		return toolError(err.Error()), nil
	}
	return toolJSON(result)
}

func (s *Service) handleGetSongDetail(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	songID := toInt64(args["song_id"])
	if songID == 0 {
		return toolError("song_id is required"), nil
	}
	result, err := s.GetSongDetail(ctx, songID)
	if err != nil {
		return toolError(err.Error()), nil
	}
	return toolJSON(result)
}

func (s *Service) handleGetLyrics(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	songID := toInt64(args["song_id"])
	if songID == 0 {
		return toolError("song_id is required"), nil
	}
	result, err := s.GetLyrics(ctx, songID)
	if err != nil {
		return toolError(err.Error()), nil
	}
	return toolJSON(result)
}

func (s *Service) handleDownloadSong(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	songID := toInt64(args["song_id"])
	if songID == 0 {
		return toolError("song_id is required"), nil
	}
	quality := "exhigh"
	if q, ok := args["quality"].(string); ok {
		quality = q
	}
	outputDir := ""
	if d, ok := args["output_dir"].(string); ok {
		outputDir = d
	}
	result, err := s.DownloadSong(ctx, songID, quality, outputDir)
	if err != nil {
		return toolError(err.Error()), nil
	}
	return toolJSON(result)
}

func (s *Service) handleGetUserPlaylists(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	userID := toInt64(args["user_id"]) // 可选，默认0表示当前用户
	limit, offset := 30, 0
	if l, ok := args["limit"]; ok {
		limit = toInt(l)
	}
	if o, ok := args["offset"]; ok {
		offset = toInt(o)
	}
	result, err := s.GetUserPlaylists(ctx, userID, limit, offset)
	if err != nil {
		return toolError(err.Error()), nil
	}
	return toolJSON(result)
}

func (s *Service) handleGetPlaylistDetail(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	playlistID := toInt64(args["playlist_id"])
	if playlistID == 0 {
		return toolError("playlist_id is required"), nil
	}
	includeSongs := true
	if v, ok := args["include_songs"].(bool); ok {
		includeSongs = v
	}
	result, err := s.GetPlaylistDetail(ctx, playlistID, includeSongs)
	if err != nil {
		return toolError(err.Error()), nil
	}
	return toolJSON(result)
}

func (s *Service) handleGetPlaylistSongs(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	playlistID := toInt64(args["playlist_id"])
	if playlistID == 0 {
		return toolError("playlist_id is required"), nil
	}
	getAll := false
	if g, ok := args["get_all"].(bool); ok {
		getAll = g
	}
	result, err := s.GetPlaylistSongs(ctx, playlistID, getAll)
	if err != nil {
		return toolError(err.Error()), nil
	}
	return toolJSON(result)
}

func (s *Service) handleGetDailyRecommend(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result, err := s.GetDailyRecommend(ctx)
	if err != nil {
		return toolError(err.Error()), nil
	}
	return toolJSON(result)
}

func (s *Service) handleGetPersonalFM(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result, err := s.GetPersonalFM(ctx)
	if err != nil {
		return toolError(err.Error()), nil
	}
	return toolJSON(result)
}

func (s *Service) handleGetSimilarSongs(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	songID := toInt64(args["song_id"])
	if songID == 0 {
		return toolError("song_id is required"), nil
	}
	result, err := s.GetSimilarSongs(ctx, songID)
	if err != nil {
		return toolError(err.Error()), nil
	}
	return toolJSON(result)
}

func (s *Service) handleLikeSong(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	songID := toInt64(args["song_id"])
	if songID == 0 {
		return toolError("song_id is required"), nil
	}
	if err := s.LikeSong(ctx, songID); err != nil {
		return toolError(err.Error()), nil
	}
	return toolText("已添加到喜欢"), nil
}

func (s *Service) handleUnlikeSong(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	songID := toInt64(args["song_id"])
	if songID == 0 {
		return toolError("song_id is required"), nil
	}
	if err := s.UnlikeSong(ctx, songID); err != nil {
		return toolError(err.Error()), nil
	}
	return toolText("已取消喜欢"), nil
}

func (s *Service) handleDailySignin(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result, err := s.DailySignin(ctx)
	if err != nil {
		return toolError(err.Error()), nil
	}
	return toolText(result), nil
}

func (s *Service) handleGetAlbumDetail(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	albumID := toInt64(args["album_id"])
	if albumID == 0 {
		return toolError("album_id is required"), nil
	}
	result, err := s.GetAlbumDetail(ctx, albumID)
	if err != nil {
		return toolError(err.Error()), nil
	}
	return toolJSON(result)
}

func (s *Service) handleGetArtistSongs(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	artistID := toInt64(args["artist_id"])
	if artistID == 0 {
		return toolError("artist_id is required"), nil
	}
	limit, offset := 30, 0
	if l, ok := args["limit"]; ok {
		limit = toInt(l)
	}
	if o, ok := args["offset"]; ok {
		offset = toInt(o)
	}
	result, err := s.GetArtistSongs(ctx, artistID, limit, offset)
	if err != nil {
		return toolError(err.Error()), nil
	}
	return toolJSON(result)
}

func (s *Service) handleGetUserLikes(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	userID := toInt64(args["user_id"]) // 可选，默认0表示当前用户
	result, err := s.GetUserLikes(ctx, userID)
	if err != nil {
		return toolError(err.Error()), nil
	}
	return toolJSON(result)
}

func (s *Service) handleCreatePlaylist(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	name, _ := args["name"].(string)
	if name == "" {
		return toolError("name is required"), nil
	}
	privacy := false
	if p, ok := args["privacy"].(bool); ok {
		privacy = p
	}
	id, playlistName, err := s.CreatePlaylist(ctx, name, privacy)
	if err != nil {
		return toolError(err.Error()), nil
	}
	return toolJSON(map[string]interface{}{
		"playlist_id": id,
		"name":        playlistName,
		"url":         fmt.Sprintf("https://music.163.com/#/playlist?id=%d", id),
	})
}

func (s *Service) handleDeletePlaylist(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	playlistID := toInt64(args["playlist_id"])
	if playlistID == 0 {
		return toolError("playlist_id is required"), nil
	}
	if err := s.DeletePlaylist(ctx, playlistID); err != nil {
		return toolError(err.Error()), nil
	}
	return toolText("歌单已删除"), nil
}

func (s *Service) handleRenamePlaylist(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	playlistID := toInt64(args["playlist_id"])
	if playlistID == 0 {
		return toolError("playlist_id is required"), nil
	}
	name, _ := args["name"].(string)
	if name == "" {
		return toolError("name is required"), nil
	}
	if err := s.RenamePlaylist(ctx, playlistID, name); err != nil {
		return toolError(err.Error()), nil
	}
	return toolText("歌单已重命名"), nil
}

func (s *Service) handleAddToPlaylist(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	playlistID := toInt64(args["playlist_id"])
	if playlistID == 0 {
		return toolError("playlist_id is required"), nil
	}
	songIDs := toInt64Slice(args["song_ids"])
	if len(songIDs) == 0 {
		return toolError("song_ids is required"), nil
	}
	if err := s.AddToPlaylist(ctx, playlistID, songIDs); err != nil {
		return toolError(err.Error()), nil
	}
	return toolText(fmt.Sprintf("已添加 %d 首歌曲到歌单", len(songIDs))), nil
}

func (s *Service) handleRemoveFromPlaylist(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	playlistID := toInt64(args["playlist_id"])
	if playlistID == 0 {
		return toolError("playlist_id is required"), nil
	}
	songIDs := toInt64Slice(args["song_ids"])
	if len(songIDs) == 0 {
		return toolError("song_ids is required"), nil
	}
	if err := s.RemoveFromPlaylist(ctx, playlistID, songIDs); err != nil {
		return toolError(err.Error()), nil
	}
	return toolText(fmt.Sprintf("已从歌单删除 %d 首歌曲", len(songIDs))), nil
}

// ========== Utility ==========

func toolText(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{mcp.TextContent{Type: "text", Text: text}},
	}
}

func toolError(msg string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{mcp.TextContent{Type: "text", Text: fmt.Sprintf("Error: %s", msg)}},
		IsError: true,
	}
}

func toolJSON(v interface{}) (*mcp.CallToolResult, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return toolError(fmt.Sprintf("JSON marshal failed: %v", err)), nil
	}
	return toolText(string(data)), nil
}

func toInt(v interface{}) int {
	switch val := v.(type) {
	case float64:
		return int(val)
	case int:
		return val
	case string:
		n, _ := strconv.Atoi(val)
		return n
	default:
		return 0
	}
}

func toInt64(v interface{}) int64 {
	switch val := v.(type) {
	case float64:
		return int64(val)
	case int64:
		return val
	case int:
		return int64(val)
	case string:
		n, _ := strconv.ParseInt(val, 10, 64)
		return n
	case json.Number:
		n, _ := val.Int64()
		return n
	default:
		return 0
	}
}

func toInt64Slice(v interface{}) []int64 {
	arr, ok := v.([]interface{})
	if !ok {
		return nil
	}
	result := make([]int64, 0, len(arr))
	for _, item := range arr {
		if id := toInt64(item); id != 0 {
			result = append(result, id)
		}
	}
	return result
}
