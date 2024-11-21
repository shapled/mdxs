package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

// Version information
var (
	Version = "v0.1.0"
	Commit  = "unknown"
	Date    = "unknown"
)

var (
	dir  string
	port int
	host string
	ver  bool
)

func init() {
	flag.StringVar(&dir, "dir", ".", "Directory containing markdown files")
	flag.IntVar(&port, "port", 3000, "Port to serve on")
	flag.StringVar(&host, "host", "127.0.0.1", "Host address to bind")
	flag.BoolVar(&ver, "version", false, "Show version information")
}

func main() {
	flag.Parse()

	if ver {
		fmt.Printf("mdxs %s (%s) %s\n", Version, Commit, Date)
		return
	}

	// 验证目录是否存在
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Fatalf("Directory %s does not exist", dir)
	}

	r := gin.Default()

	// 处理markdown文件
	r.GET("/", handleMarkdown)
	r.GET("/:path", handleMarkdown)

	// 确定服务器地址
	addr := fmt.Sprintf("%s:%d", host, port)

	log.Printf("mdxs %s starting server on http://%s", Version, addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}

func handleMarkdown(c *gin.Context) {
	path := c.Param("path")
	if path == "" {
		path = "index.md"
	}
	if !strings.HasSuffix(path, ".md") {
		path += ".md"
	}

	fullPath := filepath.Join(dir, path)

	// 检查文件是否存在
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		c.String(http.StatusNotFound, "File not found")
		return
	}

	// 读取文件内容
	content, err := os.ReadFile(fullPath)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error reading file")
		return
	}

	// 解析markdown
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)
	html := markdown.ToHTML(content, p, nil)

	// 返回HTML页面
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, `
<!DOCTYPE html>
<html>
<head>
    <title>MDXS</title>
    <meta charset="utf-8">
    <style>
        body {
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            font-family: system-ui, -apple-system, sans-serif;
            line-height: 1.6;
        }
        pre {
            background: #f6f8fa;
            padding: 16px;
            border-radius: 6px;
            overflow-x: auto;
        }
        code {
            font-family: ui-monospace, monospace;
        }
    </style>
</head>
<body>
    %s
</body>
</html>`, string(html))
}
