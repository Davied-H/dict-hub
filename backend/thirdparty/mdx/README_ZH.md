# MDX/MDD Parser for Go
[EN](README.md)

这是一个使用 Go 语言编写的高性能 MDict (.mdx/.mdd) 文件解析库。它支持查询词典内容、获取词典元信息，并提供了符合 `io/fs.FS` 接口的文件系统封装，方便与其他 Go 生态库（如 HTTP 服务器）集成。

该库最初基于 [terasum/medict](https://github.com/terasum/medict) 项目，并在其基础上进行了大量的 bug 修复、性能优化和代码重构。

## 特性

- **高性能查询**: 使用二分查找算法，对词条进行 O(log n) 复杂度的快速查询。
- **MDX/MDD 支持**: 同时支持 .mdx（文本词典）和 .mdd（资源文件）格式。
- **标准接口**: 实现 `io/fs.FS` 接口，可以轻松地将词典作为文件系统提供服务。
- **健壮的错误处理**: 完善的错误处理和日志记录。
- **完整的元信息**: 提供访问词典所有元数据（如标题、描述、创建日期等）的 API。

## 安装

```bash
go get github.com/lib-x/mdx
```
*(注意: 请将 `github.com/lib-x/mdx` 替换为实际的仓库路径)*

## 用法示例

### 示例 1: 查询 MDX 词典

下面是一个如何加载 MDX 词典并查询单词的简单示例。

```go
package main

import (
	"fmt"
	"log"

	"github.com/lib-x/mdx" // 假设模块路径为此
)

func main() {
	// 1. 创建一个新的 Mdict 实例
	// 将 "path/to/your/dictionary.mdx" 替换为你的 MDX 文件路径
	mdict, err := mdx.New("path/to/your/dictionary.mdx")
	if err != nil {
		log.Fatalf("无法加载词典文件: %v", err)
	}

	// 2. 构建索引（建议在程序启动时执行一次）
	err = mdict.BuildIndex()
	if err != nil {
		log.Fatalf("无法构建词典索引: %v", err)
	}

	// 3. 打印词典信息
	fmt.Printf("词典名称: %s\n", mdict.Title())
	fmt.Printf("词典描述: %s\n", mdict.Description())

	// 4. 查询一个单词
	word := "hello"
	definition, err := mdict.Lookup(word)
	if err != nil {
		log.Fatalf("查询单词 '%s' 失败: %v", word, err)
	}

	fmt.Printf("'%s' 的释义:\n%s\n", word, string(definition))

	// 5. 查询一个不存在的单词
	word = "nonexistentword"
	_, err = mdict.Lookup(word)
	if err != nil {
		fmt.Printf("查询一个不存在的单词 '%s' 时按预期出错: %v\n", word, err)
	}
}
```

### 示例 2: 列出 MDD 资源文件

MDD 文件通常包含音频、图片等资源。下面的示例展示了如何列出 MDD 文件中的所有资源。

```go
package main

import (
	"fmt"
	"log"

	"github.com/lib-x/mdx" // 假设模块路径为此
)

func main() {
	// 1. 加载 MDD 文件
	// 将 "path/to/your/resource.mdd" 替换为你的 MDD 文件路径
	mdd, err := mdx.New("path/to/your/resource.mdd")
	if err != nil {
		log.Fatalf("无法加载 MDD 文件: %v", err)
	}

	// 2. 构建索引
	err = mdd.BuildIndex()
	if err != nil {
		log.Fatalf("无法构建 MDD 索引: %v", err)
	}

	// 3. 获取并打印所有关键字条目（在 MDD 中，这通常是文件路径）
	entries, err := mdd.GetKeyWordEntries()
	if err != nil {
		log.Fatalf("无法获取关键字条目: %v", err)
	}

	fmt.Printf("在 '%s' 中找到 %d 个资源文件:\n", mdd.Name(), len(entries))
	for i, entry := range entries {
		// 只打印前 10 个作为示例
		if i >= 10 {
			break
		}
		fmt.Println(entry.KeyWord) // KeyWord 字段存储了资源的文件路径
	}
}
```

## 贡献

欢迎提交问题（Issues）和拉取请求（Pull Requests）。

## 许可

本项目基于 [GNU General Public License v3.0](LICENSE) 授权。
