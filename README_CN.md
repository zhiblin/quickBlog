# QuickBlog

一个用 Go 编写的轻量级静态博客服务器。只需将 `.html` 或 `.md` 文件放入指定目录，即可自动生成博客页面——无需数据库、无需配置文件、无需框架。

## 特性

- **零配置** — 一个可执行文件 + 一个 `pages/` 目录即可运行
- **Markdown 支持** — 用 Markdown 写作，一条命令转为 HTML（支持 GFM：表格、任务列表、代码块、删除线）
- **自动排序** — 首页按文件修改时间倒序展示所有文章，最新的在最前
- **精美设计** — 温暖极简风模板，Inter 字体，响应式布局
- **部署简单** — 单二进制文件，模板已内嵌，支持交叉编译到任意平台

## 快速开始

```bash
# 编译
go build -o quickblog .

# 写一篇 Markdown 文章
echo '# Hello World\n\n我的第一篇文章！' > posts/hello.md

# 将 Markdown 转换为 HTML
./quickblog convert

# 启动服务
./quickblog
```

访问 `http://localhost:8080` 查看你的博客。

## 使用方法

### 启动服务

```bash
# 默认：端口 8080，读取 ./pages 目录
./quickblog

# 自定义端口
./quickblog -port 3000

# 自定义页面目录
./quickblog -dir /var/www/pages

# 转换 Markdown + 启动服务（一步完成）
./quickblog -md posts
```

### 转换 Markdown 为 HTML

```bash
# 默认：posts/ -> pages/
./quickblog convert

# 自定义源目录和输出目录
./quickblog convert -md src -out dist
```

## 项目结构

```
quickBlog/
├── main.go          # 服务器 + Markdown 转换器（约 230 行）
├── go.mod
├── templates/       # 内嵌的 HTML 模板
│   ├── index.html   # 首页（文章列表）
│   └── page.html    # 文章详情页
├── posts/           # Markdown 源文件（.md）
└── pages/           # 发布的 HTML 文件（.html）
```

## 工作流程

1. 在 `posts/` 目录中编写 `.md` 文件（或直接在 `pages/` 中放入 `.html` 文件）
2. 运行 `./quickblog convert` 将 Markdown 转换为 HTML
3. 运行 `./quickblog` 启动服务
4. 首页自动列出 `pages/` 中的所有 `.html` 文件，按修改时间排序

## 交叉编译

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o quickblog .

# macOS ARM
GOOS=darwin GOARCH=arm64 go build -o quickblog .

# Windows
GOOS=windows GOARCH=amd64 go build -o quickblog.exe .
```

## 依赖

- [Goldmark](https://github.com/yuin/goldmark) — 兼容 CommonMark 标准的 Markdown 解析器，支持 GFM 扩展

## 许可证

MIT
