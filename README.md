# MDXS

MDXS is a Go-based web server that renders Markdown/MDX files and serves them as web pages.

## Features

- Renders Markdown files with MDX support
- Provides a web interface to browse and view rendered content
- Watches for file changes and updates content automatically

## Installation

```bash
go install github.com/shapled/mdxs@latest
```

## Usage

```bash
mdxs -dir /path/to/markdown/files -port 3000
```

## Development

Requirements:
- Go 1.23+

To run locally:
```bash
go run main.go -dir ./docs -port 3000
```
