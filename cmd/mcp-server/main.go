package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/mark3labs/mcp-go/server"

	"github.com/go-musicfox/go-musicfox/internal/mcp"
)

func main() {
	// 初始化日志
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	// 创建 MCP Service
	svc, err := mcp.NewService()
	if err != nil {
		fmt.Fprintf(os.Stderr, "初始化失败: %v\n", err)
		os.Exit(1)
	}

	// 创建 MCP Server
	mcpServer := server.NewMCPServer(
		"musicfox",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	// 注册所有 Tool
	svc.RegisterTools(mcpServer)

	// 启动 stdio transport
	slog.Info("musicfox MCP server 启动")
	if err := server.ServeStdio(mcpServer); err != nil {
		fmt.Fprintf(os.Stderr, "MCP server 错误: %v\n", err)
		os.Exit(1)
	}
}
