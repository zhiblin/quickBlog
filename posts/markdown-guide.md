# Markdown 快速指南

Markdown 是一种轻量级标记语言，让你用纯文本写出格式丰富的文档。

## 基本语法

**粗体**、*斜体*、~~删除线~~

## 列表

- 项目一
- 项目二
  - 子项目

## 任务列表

- [x] 搭建 QuickBlog
- [x] 支持 Markdown 转换
- [ ] 添加更多主题

## 代码块

```go
func main() {
    fmt.Println("Hello, QuickBlog!")
}
```

## 表格

| 功能 | 状态 |
|------|------|
| HTML 页面 | 已完成 |
| Markdown 转换 | 已完成 |
| 自定义端口 | 已完成 |

> 只需将 `.md` 文件放入 `posts/` 目录，运行 `convert` 即可发布！
