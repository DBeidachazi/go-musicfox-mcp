package mcp

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-musicfox/netease-music/service"
	musicutil "github.com/go-musicfox/netease-music/util"
	cookiejar "github.com/juju/persistent-cookiejar"

	"github.com/go-musicfox/go-musicfox/internal/configs"
	"github.com/go-musicfox/go-musicfox/internal/netease"
	"github.com/go-musicfox/go-musicfox/internal/structs"
	"github.com/go-musicfox/go-musicfox/internal/track"
	"github.com/go-musicfox/go-musicfox/utils/app"
	neteaseutil "github.com/go-musicfox/go-musicfox/utils/netease"
	_struct "github.com/go-musicfox/go-musicfox/utils/struct"
)

// Service 是 MCP server 的核心业务层，封装了网易云音乐的所有操作。
type Service struct {
	trackManager *track.Manager
	cookiePath   string
}

// NewService 创建一个新的 MCP Service。
// 它会加载配置、初始化 cookie 认证、创建 track manager。
func NewService() (*Service, error) {
	// 加载配置
	configPath := app.ConfigFilePath()
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 配置文件不存在时使用默认配置
		slog.Warn("配置文件不存在，使用默认配置", "path", configPath)
	}

	cfg, err := configs.NewConfigFromTomlFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("加载配置失败: %w", err)
	}
	configs.AppConfig = cfg

	// 更新 netease 配置
	musicutil.UNMSwitch = cfg.UNM.Enable
	musicutil.Sources = cfg.UNM.Sources
	musicutil.SearchLimit = cfg.UNM.SearchLimit
	musicutil.EnableLocalVip = cfg.UNM.EnableLocalVip
	musicutil.UnlockSoundEffects = cfg.UNM.UnlockSoundEffects
	musicutil.UNMProxyURL = cfg.UNM.ProxyURL

	// 初始化 cookie 认证
	cookiePath := filepath.Join(app.DataDir(), "cookie")
	if err := initCookie(cookiePath); err != nil {
		return nil, fmt.Errorf("初始化 cookie 失败: %w", err)
	}

	// 创建 track manager（用于下载和获取歌曲信息）
	quality := cfg.Player.SongLevel
	mgr := track.NewManager(
		track.WithSongQuality(quality),
		track.WithCacheLimit(int64(cfg.Storage.Cache.Limit)),
	)

	return &Service{
		trackManager: mgr,
		cookiePath:   cookiePath,
	}, nil
}

// initCookie 从文件加载 cookie 并设置到全局 cookie jar。
func initCookie(cookiePath string) error {
	jar, err := cookiejar.New(&cookiejar.Options{Filename: cookiePath})
	if err != nil {
		return fmt.Errorf("创建 cookie jar 失败: %w", err)
	}
	// persistent-cookiejar 在 New 时自动从文件加载 cookie
	musicutil.SetGlobalCookieJar(jar)
	return nil
}

// parseQuality 将配置中的音质字符串转换为 service.SongQualityLevel。
func parseQuality(quality string) service.SongQualityLevel {
	switch strings.ToLower(quality) {
	case "exhigh":
		return service.Exhigh
	case "lossless":
		return service.Lossless
	case "hires":
		return service.Hires
	case "higher":
		return service.Higher
	default:
		return service.Standard
	}
}

// SearchType 搜索类型常量
type SearchType string

const (
	SearchSongs     SearchType = "songs"     // 1
	SearchAlbums    SearchType = "albums"     // 10
	SearchArtists   SearchType = "artists"    // 100
	SearchPlaylists SearchType = "playlists"  // 1000
	SearchLyrics    SearchType = "lyrics"     // 1006
	SearchDjRadio   SearchType = "djradio"    // 1009
)

// searchTypeMap 搜索类型到 API type 参数的映射
var searchTypeMap = map[SearchType]string{
	SearchSongs:     "1",
	SearchAlbums:    "10",
	SearchArtists:   "100",
	SearchPlaylists: "1000",
	SearchLyrics:    "1006",
	SearchDjRadio:   "1009",
}

// SongInfo 用于 MCP 输出的歌曲信息
type SongInfo struct {
	ID       int64    `json:"id"`
	Name     string   `json:"name"`
	Artists  []string `json:"artists"`
	Album    string   `json:"album,omitempty"`
	AlbumID  int64    `json:"album_id,omitempty"`
	Duration string   `json:"duration"`
	URL      string   `json:"url"`
}

// PlaylistInfo 用于 MCP 输出的歌单信息
type PlaylistInfo struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Creator    string `json:"creator"`
	TrackCount int    `json:"track_count,omitempty"`
	URL        string `json:"url"`
}

// AlbumInfo 用于 MCP 输出的专辑信息
type AlbumInfo struct {
	ID      int64    `json:"id"`
	Name    string   `json:"name"`
	Artists []string `json:"artists"`
	URL     string   `json:"url"`
}

// ArtistInfo 用于 MCP 输出的歌手信息
type ArtistInfo struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	SongCount int    `json:"song_count,omitempty"`
	URL       string `json:"url"`
}

// LyricResult 歌词结果
type LyricResult struct {
	SongID     int64  `json:"song_id"`
	SongName   string `json:"song_name"`
	Original   string `json:"original"`
	Translated string `json:"translated,omitempty"`
	Yrc        string `json:"yrc,omitempty"`
}

// SearchResult 搜索结果
type SearchResult struct {
	Type    string         `json:"type"`
	Songs   []SongInfo     `json:"songs,omitempty"`
	Albums  []AlbumInfo    `json:"albums,omitempty"`
	Artists []ArtistInfo   `json:"artists,omitempty"`
	Playlists []PlaylistInfo `json:"playlists,omitempty"`
	Total   int            `json:"total"`
}

// SongURLResult 歌曲 URL 结果
type SongURLResult struct {
	SongID    int64  `json:"song_id"`
	URL       string `json:"url"`
	Type      string `json:"type"`
	Bitrate   int    `json:"bitrate"`
	Size      int64  `json:"size"`
	Quality   string `json:"quality"`
}

// DownloadResult 下载结果
type DownloadResult struct {
	SongID   int64  `json:"song_id"`
	FilePath string `json:"file_path"`
	Size     int64  `json:"size"`
}

// PlayResult 播放结果
type PlayResult struct {
	SongID   int64    `json:"song_id"`
	SongName string   `json:"song_name"`
	Artists  []string `json:"artists"`
	State    string   `json:"state"`
}

// PlayingInfo 当前播放信息
type PlayingInfo struct {
	SongID    int64    `json:"song_id"`
	SongName  string   `json:"song_name"`
	Artists   []string `json:"artists"`
	Album     string   `json:"album"`
	Position  string   `json:"position"`
	Duration  string   `json:"duration"`
	State     string   `json:"state"`
	Volume    int      `json:"volume"`
	PlayMode  string   `json:"play_mode"`
}

// ---------- Tier 1: 纯 API 查询 ----------

// Search 搜索网易云音乐
func (s *Service) Search(ctx context.Context, keywords string, searchType SearchType, limit, offset int) (*SearchResult, error) {
	typeStr, ok := searchTypeMap[searchType]
	if !ok {
		typeStr = "1"
		searchType = SearchSongs
	}

	searchService := service.SearchService{
		S:      keywords,
		Type:   typeStr,
		Limit:  strconv.Itoa(limit),
		Offset: strconv.Itoa(offset),
	}
	code, response := searchService.Search()
	codeType := _struct.CheckCode(code)
	if codeType != _struct.Success {
		return nil, fmt.Errorf("搜索失败，错误码: %v", code)
	}

	result := &SearchResult{Type: string(searchType)}

	switch searchType {
	case SearchSongs:
		songs := _struct.GetSongsOfSearchResult(response)
		result.Songs = convertSongs(songs)
		result.Total = len(songs)
	case SearchAlbums:
		albums := _struct.GetAlbumsOfSearchResult(response)
		result.Albums = convertAlbums(albums)
		result.Total = len(albums)
	case SearchArtists:
		artists := _struct.GetArtistsOfSearchResult(response)
		result.Artists = convertArtists(artists)
		result.Total = len(artists)
	case SearchPlaylists:
		playlists := _struct.GetPlaylistsOfSearchResult(response)
		result.Playlists = convertPlaylists(playlists)
		result.Total = len(playlists)
	}

	return result, nil
}

// GetSongURL 获取歌曲播放 URL
func (s *Service) GetSongURL(ctx context.Context, songID int64, quality string) (*SongURLResult, error) {
	q := parseQuality(quality)
	info, err := neteaseutil.FetchPlayableInfo(songID, q)
	if err != nil {
		return nil, fmt.Errorf("获取歌曲 URL 失败: %w", err)
	}

	// 修复 CDN 403 问题
	info.URL = fixCDNDomain(info.URL)

	if info.URL == "" {
		return nil, fmt.Errorf("歌曲 %d 无法获取播放链接（可能需要 VIP 或地区限制）", songID)
	}

	return &SongURLResult{
		SongID:  songID,
		URL:     info.URL,
		Type:    info.MusicType,
		Size:    info.Size,
		Quality: quality,
	}, nil
}

// GetLyrics 获取歌词
func (s *Service) GetLyrics(ctx context.Context, songID int64) (*LyricResult, error) {
	lrcData, err := s.trackManager.GetLyric(ctx, songID)
	if err != nil {
		return nil, fmt.Errorf("获取歌词失败: %w", err)
	}

	return &LyricResult{
		SongID:     songID,
		Original:   lrcData.Original,
		Translated: lrcData.Translated,
		Yrc:        lrcData.Yrc,
	}, nil
}

// DownloadSong 下载歌曲到本地
func (s *Service) DownloadSong(ctx context.Context, songID int64, quality string, outputDir string) (*DownloadResult, error) {
	// 先获取歌曲信息
	song := structs.Song{Id: songID}

	// 如果指定了输出目录，临时覆盖
	if outputDir != "" {
		s.trackManager = track.NewManager(
			track.WithSongQuality(parseQuality(quality)),
			track.WithDownloadDir(outputDir),
			track.WithCacheLimit(0),
		)
	}

	filePath, err := s.trackManager.DownloadSong(ctx, song)
	if err != nil {
		return nil, fmt.Errorf("下载歌曲失败: %w", err)
	}

	// 获取文件大小
	var size int64
	if fi, err := os.Stat(filePath); err == nil {
		size = fi.Size()
	}

	return &DownloadResult{
		SongID:   songID,
		FilePath: filePath,
		Size:     size,
	}, nil
}

// GetUserPlaylists 获取用户歌单列表
func (s *Service) GetUserPlaylists(ctx context.Context, userID int64, limit, offset int) ([]PlaylistInfo, error) {
	codeType, playlists, _ := netease.FetchUserPlaylists(userID, limit, offset)
	if codeType != _struct.Success {
		return nil, fmt.Errorf("获取歌单失败")
	}
	return convertPlaylists(playlists), nil
}

// GetPlaylistSongs 获取歌单内的歌曲
func (s *Service) GetPlaylistSongs(ctx context.Context, playlistID int64, getAll bool) ([]SongInfo, error) {
	codeType, songs := netease.FetchSongsOfPlaylist(playlistID, getAll)
	if codeType != _struct.Success {
		return nil, fmt.Errorf("获取歌单歌曲失败")
	}
	return convertSongs(songs), nil
}

// GetDailyRecommend 获取每日推荐歌曲
func (s *Service) GetDailyRecommend(ctx context.Context) ([]SongInfo, error) {
	songs, err := netease.FetchDailySongs()
	if err != nil {
		return nil, fmt.Errorf("获取每日推荐失败: %w", err)
	}
	return convertSongs(songs), nil
}

// GetPersonalFM 获取私人 FM
func (s *Service) GetPersonalFM(ctx context.Context) ([]SongInfo, error) {
	fmService := service.PersonalFmService{}
	code, response := fmService.PersonalFm()
	codeType := _struct.CheckCode(code)
	if codeType != _struct.Success {
		return nil, fmt.Errorf("获取私人FM失败")
	}
	songs := _struct.GetFmSongs(response)
	return convertSongs(songs), nil
}

// GetSimilarSongs 获取相似歌曲
func (s *Service) GetSimilarSongs(ctx context.Context, songID int64) ([]SongInfo, error) {
	simiService := service.SimiSongService{ID: strconv.FormatInt(songID, 10)}
	code, response := simiService.SimiSong()
	codeType := _struct.CheckCode(code)
	if codeType != _struct.Success {
		return nil, fmt.Errorf("获取相似歌曲失败")
	}
	songs := _struct.GetSimiSongs(response)
	return convertSongs(songs), nil
}

// LikeSong 喜欢歌曲
func (s *Service) LikeSong(ctx context.Context, songID int64) error {
	likeService := service.LikeService{ID: strconv.FormatInt(songID, 10), L: "true"}
	code, _ := likeService.Like()
	if code != 200 {
		return fmt.Errorf("喜欢歌曲失败，错误码: %v", code)
	}
	return nil
}

// UnlikeSong 取消喜欢
func (s *Service) UnlikeSong(ctx context.Context, songID int64) error {
	likeService := service.LikeService{ID: strconv.FormatInt(songID, 10), L: "false"}
	code, _ := likeService.Like()
	if code != 200 {
		return fmt.Errorf("取消喜欢失败，错误码: %v", code)
	}
	return nil
}

// DailySignin 每日签到
func (s *Service) DailySignin(ctx context.Context) (string, error) {
	signinService := service.DailySigninService{Type: "0"}
	code, response := signinService.DailySignin()
	if code != 200 {
		return "", fmt.Errorf("签到失败，错误码: %v, 响应: %s", code, string(response))
	}
	return "签到成功", nil
}

// GetAlbumDetail 获取专辑详情
func (s *Service) GetAlbumDetail(ctx context.Context, albumID int64) (interface{}, error) {
	albumService := service.AlbumDetailService{ID: strconv.FormatInt(albumID, 10)}
	code, response := albumService.AlbumDetail()
	codeType := _struct.CheckCode(code)
	if codeType != _struct.Success {
		return nil, fmt.Errorf("获取专辑详情失败")
	}
	songs := _struct.GetSongsOfAlbum(response)
	return map[string]interface{}{
		"songs": convertSongs(songs),
	}, nil
}

// GetArtistSongs 获取歌手的歌曲
func (s *Service) GetArtistSongs(ctx context.Context, artistID int64, limit, offset int) ([]SongInfo, error) {
	artistService := service.ArtistSongsService{
		ID:     strconv.FormatInt(artistID, 10),
		Limit:  strconv.Itoa(limit),
		Offset: strconv.Itoa(offset),
	}
	code, response := artistService.ArtistSongs()
	codeType := _struct.CheckCode(code)
	if codeType != _struct.Success {
		return nil, fmt.Errorf("获取歌手歌曲失败")
	}
	songs := _struct.GetSongsOfArtist(response)
	return convertSongs(songs), nil
}

// GetUserLikes 获取用户喜欢的歌曲
func (s *Service) GetUserLikes(ctx context.Context, userID int64) ([]SongInfo, error) {
	songs, err := netease.FetchLikeSongs(userID, true)
	if err != nil {
		return nil, fmt.Errorf("获取喜欢歌曲失败: %w", err)
	}
	return convertSongs(songs), nil
}

// ---------- 工具函数 ----------

// fixCDNDomain 修复 CDN 403 问题，将不可用的 CDN 域名替换为可用的。
var cdnDomainRegexp = regexp.MustCompile(`m\d+\.music\.126\.net`)

func fixCDNDomain(url string) string {
	if url == "" {
		return url
	}
	// m704 和 m10 实测返回 403，替换为 m701
	return cdnDomainRegexp.ReplaceAllString(url, "m701.music.126.net")
}

func convertSongs(songs []structs.Song) []SongInfo {
	result := make([]SongInfo, 0, len(songs))
	for _, s := range songs {
		artists := make([]string, 0, len(s.Artists))
		for _, a := range s.Artists {
			artists = append(artists, a.Name)
		}
		result = append(result, SongInfo{
			ID:       s.Id,
			Name:     s.Name,
			Artists:  artists,
			Album:    s.Album.Name,
			AlbumID:  s.Album.Id,
			Duration: s.Duration.String(),
			URL:      neteaseutil.WebUrlOfSong(s.Id),
		})
	}
	return result
}

func convertAlbums(albums []structs.Album) []AlbumInfo {
	result := make([]AlbumInfo, 0, len(albums))
	for _, a := range albums {
		artists := make([]string, 0, len(a.Artists))
		for _, ar := range a.Artists {
			artists = append(artists, ar.Name)
		}
		result = append(result, AlbumInfo{
			ID:      a.Id,
			Name:    a.Name,
			Artists: artists,
			URL:     neteaseutil.WebUrlOfAlbum(a.Id),
		})
	}
	return result
}

func convertArtists(artists []structs.Artist) []ArtistInfo {
	result := make([]ArtistInfo, 0, len(artists))
	for _, a := range artists {
		result = append(result, ArtistInfo{
			ID:   a.Id,
			Name: a.Name,
			URL:  neteaseutil.WebUrlOfArtist(a.Id),
		})
	}
	return result
}

func convertPlaylists(playlists []structs.Playlist) []PlaylistInfo {
	result := make([]PlaylistInfo, 0, len(playlists))
	for _, p := range playlists {
		result = append(result, PlaylistInfo{
			ID:      p.Id,
			Name:    p.Name,
			Creator: p.Creator.Nickname,
			URL:     neteaseutil.WebUrlOfPlaylist(p.Id),
		})
	}
	return result
}

// WebUrlOfSong 等函数在 utils/struct 包中，需要确认是否存在
// 这里使用 netease 包中的函数
