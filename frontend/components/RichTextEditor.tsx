'use client'

import { useState, useRef, useEffect } from 'react'
import { 
  Bold, 
  Italic, 
  Underline, 
  Link2, 
  List, 
  ListOrdered, 
  Quote, 
  Code, 
  Heading1, 
  Heading2, 
  Heading3,
  Image as ImageIcon,
  Eye,
  FileText
} from 'lucide-react'

interface RichTextEditorProps {
  value: string
  onChange: (value: string) => void
  placeholder?: string
  className?: string
}

export default function RichTextEditor({ value, onChange, placeholder, className }: RichTextEditorProps) {
  const [isPreview, setIsPreview] = useState(false)
  const textareaRef = useRef<HTMLTextAreaElement>(null)

  const insertText = (before: string, after: string = '') => {
    const textarea = textareaRef.current
    if (!textarea) return

    const start = textarea.selectionStart
    const end = textarea.selectionEnd
    const selectedText = value.substring(start, end)
    
    const newText = value.substring(0, start) + before + selectedText + after + value.substring(end)
    onChange(newText)

    // Restore cursor position
    setTimeout(() => {
      textarea.focus()
      textarea.setSelectionRange(start + before.length, start + before.length + selectedText.length)
    }, 0)
  }

  const insertAtCursor = (text: string) => {
    const textarea = textareaRef.current
    if (!textarea) return

    const start = textarea.selectionStart
    const newText = value.substring(0, start) + text + value.substring(start)
    onChange(newText)

    setTimeout(() => {
      textarea.focus()
      textarea.setSelectionRange(start + text.length, start + text.length)
    }, 0)
  }

  const formatMarkdown = (content: string) => {
    return content
      .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
      .replace(/\*(.*?)\*/g, '<em>$1</em>')
      .replace(/`(.*?)`/g, '<code class="bg-muted text-muted-foreground px-1 rounded">$1</code>')
      .replace(/^# (.*$)/gm, '<h1 class="text-2xl font-bold mb-4 text-foreground">$1</h1>')
      .replace(/^## (.*$)/gm, '<h2 class="text-xl font-bold mb-3 text-foreground">$1</h2>')
      .replace(/^### (.*$)/gm, '<h3 class="text-lg font-bold mb-2 text-foreground">$1</h3>')
      .replace(/^> (.*$)/gm, '<blockquote class="border-l-4 border-primary pl-4 italic text-muted-foreground">$1</blockquote>')
      .replace(/^\* (.*$)/gm, '<li class="text-foreground">$1</li>')
      .replace(/^(\d+)\. (.*$)/gm, '<li class="text-foreground">$1. $2</li>')
      .replace(/\n/g, '<br>')
  }

  const calculateReadingTime = (text: string) => {
    const wordsPerMinute = 200
    const words = text.trim().split(/\s+/).length
    return Math.ceil(words / wordsPerMinute)
  }

  const tools = [
    { icon: Bold, action: () => insertText('**', '**'), title: 'Bold (Ctrl+B)' },
    { icon: Italic, action: () => insertText('*', '*'), title: 'Italic (Ctrl+I)' },
    { icon: Underline, action: () => insertText('<u>', '</u>'), title: 'Underline' },
    { icon: Heading1, action: () => insertAtCursor('# '), title: 'Heading 1' },
    { icon: Heading2, action: () => insertAtCursor('## '), title: 'Heading 2' },
    { icon: Heading3, action: () => insertAtCursor('### '), title: 'Heading 3' },
    { icon: Quote, action: () => insertAtCursor('> '), title: 'Quote' },
    { icon: Code, action: () => insertText('`', '`'), title: 'Inline Code' },
    { icon: List, action: () => insertAtCursor('* '), title: 'Bullet List' },
    { icon: ListOrdered, action: () => insertAtCursor('1. '), title: 'Numbered List' },
    { icon: Link2, action: () => insertText('[Link Text](', ')'), title: 'Link' },
    { icon: ImageIcon, action: () => insertText('![Alt Text](', ')'), title: 'Image' },
  ]

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.ctrlKey || e.metaKey) {
      switch (e.key) {
        case 'b':
          e.preventDefault()
          insertText('**', '**')
          break
        case 'i':
          e.preventDefault()
          insertText('*', '*')
          break
      }
    }
  }

    return (
    <div className={`border border-border rounded-lg overflow-hidden bg-card ${className}`}>
      {/* Toolbar */}
      <div className="bg-muted border-b border-border p-2">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-1">
            {tools.map((tool, index) => (
              <button
                key={index}
                type="button"
                onClick={tool.action}
                title={tool.title}
                className="p-2 hover:bg-muted-foreground/10 rounded transition-colors text-foreground"
              >
                <tool.icon size={16} />
              </button>
            ))}
          </div>
          
          <div className="flex items-center space-x-2">
            <span className="text-sm text-muted-foreground">
              {calculateReadingTime(value)} min read
            </span>
            <div className="h-4 w-px bg-border" />
            <button
              type="button"
              onClick={() => setIsPreview(!isPreview)}
              className={`p-2 rounded transition-colors ${
                isPreview ? 'bg-primary/10 text-primary' : 'hover:bg-muted-foreground/10 text-foreground'
              }`}
              title={isPreview ? 'Edit Mode' : 'Preview Mode'}
            >
              {isPreview ? <FileText size={16} /> : <Eye size={16} />}
            </button>
          </div>
        </div>
      </div>

      {/* Editor/Preview */}
      <div className="relative">
        {isPreview ? (
          <div 
            className="p-4 min-h-[400px] prose prose-sm max-w-none dark:prose-invert text-foreground"
            dangerouslySetInnerHTML={{ __html: formatMarkdown(value) }}
          />
        ) : (
          <textarea
            ref={textareaRef}
            value={value}
            onChange={(e) => onChange(e.target.value)}
            onKeyDown={handleKeyDown}
            placeholder={placeholder}
            className="w-full min-h-[400px] p-4 border-none outline-none resize-none font-mono text-sm leading-relaxed bg-card text-foreground placeholder:text-muted-foreground"
          />
        )}
      </div>

      {/* Status Bar */}
      <div className="bg-muted border-t border-border px-4 py-2 text-xs text-muted-foreground flex justify-between">
        <span>{value.length} characters</span>
        <span>{value.split(/\s+/).filter(word => word.length > 0).length} words</span>
      </div>
    </div>
  )
}
