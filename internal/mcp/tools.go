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
			Name:        "search",
			Description: "Search NetEase Cloud Music (songs/albums/artists/playlists/lyrics)",
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"keywords": map[string]interface{}{"type": "string", "description": "Search keywords"},
					"type":     map[string]interface{}{"type": "string", "description": "Search type", "enum": []string{"songs", "albums", "artists", "playlists", "lyrics", "djradio"}, "default": "songs"},
					"limit":    map[string]interface{}{"type": "integer", "description": "Max results (default 20)", "default": 20},
					"offset":   map[string]interface{}{"type": "integer", "description": "Offset (default 0)", "default": 0},
				},
				Required: []string{"keywords"},
			},
		},
		s.handleSearch,
	)

	mcpServer.AddTool(
		mcp.Tool{
			Name:        "get_song_url",
			Description: "Get the playback URL for a song (temporary, expires)",
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"song_id": map[string]interface{}{"type": "integer", "description": "Song ID"},
					"quality": map[string]interface{}{"type": "string", "description": "Audio quality", "enum": []string{"standard", "higher", "exhigh", "lossless", "hires"}, "default": "exhigh"},
				},
				Required: []string{"song_id"},
			},
		},
		s.handleGetSongURL,
	)

	mcpServer.AddTool(
		mcp.Tool{
			Name:        "get_lyrics",
			Description: "Get lyrics for a song (original, translated, word-by-word)",
			InputSchema: mcp.ToolInputSchema{
				Type:       "object",
				Properties: map[string]interface{}{"song_id": map[string]interface{}{"type": "integer", "description": "Song ID"}},
				Required:   []string{"song_id"},
			},
		},
		s.handleGetLyrics,
	)

	mcpServer.AddTool(
		mcp.Tool{
			Name:        "download_song",
			Description: "Download a song to local filesystem",
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"song_id":    map[string]interface{}{"type": "integer", "description": "Song ID"},
					"quality":    map[string]interface{}{"type": "string", "description": "Audio quality", "enum": []string{"standard", "higher", "exhigh", "lossless"}, "default": "exhigh"},
					"output_dir": map[string]interface{}{"type": "string", "description": "Output directory (optional)"},
				},
				Required: []string{"song_id"},
			},
		},
		s.handleDownloadSong,
	)

	mcpServer.AddTool(
		mcp.Tool{
			Name:        "get_user_playlists",
			Description: "Get user's playlist list",
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"user_id": map[string]interface{}{"type": "integer", "description": "User ID"},
					"limit":   map[string]interface{}{"type": "integer", "description": "Max results (default 30)", "default": 30},
					"offset":  map[string]interface{}{"type": "integer", "description": "Offset (default 0)", "default": 0},
				},
				Required: []string{"user_id"},
			},
		},
		s.handleGetUserPlaylists,
	)

	mcpServer.AddTool(
		mcp.Tool{
			Name:        "get_playlist_songs",
			Description: "Get songs in a playlist",
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"playlist_id": map[string]interface{}{"type": "integer", "description": "Playlist ID"},
					"get_all":     map[string]interface{}{"type": "boolean", "description": "Get all songs (default false, max 1000)", "default": false},
				},
				Required: []string{"playlist_id"},
			},
		},
		s.handleGetPlaylistSongs,
	)

	mcpServer.AddTool(
		mcp.Tool{
			Name:        "get_daily_recommend",
			Description: "Get daily recommended songs (requires login)",
			InputSchema: mcp.ToolInputSchema{Type: "object", Properties: map[string]interface{}{}},
		},
		s.handleGetDailyRecommend,
	)

	mcpServer.AddTool(
		mcp.Tool{
			Name:        "get_personal_fm",
			Description: "Get personal FM songs (requires login)",
			InputSchema: mcp.ToolInputSchema{Type: "object", Properties: map[string]interface{}{}},
		},
		s.handleGetPersonalFM,
	)

	mcpServer.AddTool(
		mcp.Tool{
			Name:        "get_similar_songs",
			Description: "Get songs similar to a given song",
			InputSchema: mcp.ToolInputSchema{
				Type:       "object",
				Properties: map[string]interface{}{"song_id": map[string]interface{}{"type": "integer", "description": "Song ID"}},
				Required:   []string{"song_id"},
			},
		},
		s.handleGetSimilarSongs,
	)

	mcpServer.AddTool(
		mcp.Tool{
			Name:        "like_song",
			Description: "Add a song to liked songs",
			InputSchema: mcp.ToolInputSchema{
				Type:       "object",
				Properties: map[string]interface{}{"song_id": map[string]interface{}{"type": "integer", "description": "Song ID"}},
				Required:   []string{"song_id"},
			},
		},
		s.handleLikeSong,
	)

	mcpServer.AddTool(
		mcp.Tool{
			Name:        "unlike_song",
			Description: "Remove a song from liked songs",
			InputSchema: mcp.ToolInputSchema{
				Type:       "object",
				Properties: map[string]interface{}{"song_id": map[string]interface{}{"type": "integer", "description": "Song ID"}},
				Required:   []string{"song_id"},
			},
		},
		s.handleUnlikeSong,
	)

	mcpServer.AddTool(
		mcp.Tool{
			Name:        "daily_signin",
			Description: "Perform daily sign-in for NetEase Cloud Music",
			InputSchema: mcp.ToolInputSchema{Type: "object", Properties: map[string]interface{}{}},
		},
		s.handleDailySignin,
	)

	mcpServer.AddTool(
		mcp.Tool{
			Name:        "get_album_detail",
			Description: "Get album details and song list",
			InputSchema: mcp.ToolInputSchema{
				Type:       "object",
				Properties: map[string]interface{}{"album_id": map[string]interface{}{"type": "integer", "description": "Album ID"}},
				Required:   []string{"album_id"},
			},
		},
		s.handleGetAlbumDetail,
	)

	mcpServer.AddTool(
		mcp.Tool{
			Name:        "get_artist_songs",
			Description: "Get songs by an artist",
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"artist_id": map[string]interface{}{"type": "integer", "description": "Artist ID"},
					"limit":     map[string]interface{}{"type": "integer", "description": "Max results (default 30)", "default": 30},
					"offset":    map[string]interface{}{"type": "integer", "description": "Offset (default 0)", "default": 0},
				},
				Required: []string{"artist_id"},
			},
		},
		s.handleGetArtistSongs,
	)

	mcpServer.AddTool(
		mcp.Tool{
			Name:        "get_user_likes",
			Description: "Get user's liked songs",
			InputSchema: mcp.ToolInputSchema{
				Type:       "object",
				Properties: map[string]interface{}{"user_id": map[string]interface{}{"type": "integer", "description": "User ID"}},
				Required:   []string{"user_id"},
			},
		},
		s.handleGetUserLikes,
	)
}

// ========== Tool Handlers ==========

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
	userID := toInt64(args["user_id"])
	if userID == 0 {
		return toolError("user_id is required"), nil
	}
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
	return toolText("Added to liked songs"), nil
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
	return toolText("Removed from liked songs"), nil
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
	userID := toInt64(args["user_id"])
	if userID == 0 {
		return toolError("user_id is required"), nil
	}
	result, err := s.GetUserLikes(ctx, userID)
	if err != nil {
		return toolError(err.Error()), nil
	}
	return toolJSON(result)
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
