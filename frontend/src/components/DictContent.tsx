import { useMemo } from 'react'

interface DictContentProps {
  html: string
}

/**
 * 解析并美化字典 HTML 内容
 * 支持常见 MDX 字典的 CSS 类名
 */
export function DictContent({ html }: DictContentProps) {
  // 处理 HTML 内容，移除不需要的元素，添加样式
  const processedHtml = useMemo(() => {
    let content = html

    // 移除内联样式表链接（我们用自己的样式）
    content = content.replace(/<link[^>]*>/gi, '')

    // 移除 style 标签
    content = content.replace(/<style[^>]*>[\s\S]*?<\/style>/gi, '')

    // 处理音频标签，添加类名以便隐藏
    content = content.replace(/<audio/gi, '<audio class="dict-audio"')

    return content
  }, [html])

  return (
    <div
      className="dict-content"
      dangerouslySetInnerHTML={{ __html: processedHtml }}
    />
  )
}
