'use client'

import { useState, useRef } from 'react'
import { Upload, X, Image as ImageIcon, Link2 } from 'lucide-react'

interface ImageUploaderProps {
  value?: string | undefined
  onChange: (url: string, alt?: string) => void
  label?: string
  className?: string
}

export default function ImageUploader({ value, onChange, label, className }: ImageUploaderProps) {
  const [dragActive, setDragActive] = useState(false)
  const [uploading, setUploading] = useState(false)
  const [urlInput, setUrlInput] = useState('')
  const [altText, setAltText] = useState('')
  const [showUrlInput, setShowUrlInput] = useState(false)
  const fileInputRef = useRef<HTMLInputElement>(null)

  const handleDrag = (e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
    if (e.type === 'dragenter' || e.type === 'dragover') {
      setDragActive(true)
    } else if (e.type === 'dragleave') {
      setDragActive(false)
    }
  }

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
    setDragActive(false)
    
    const files = e.dataTransfer.files
    if (files && files[0]) {
      handleFile(files[0])
    }
  }

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files
    if (files && files[0]) {
      handleFile(files[0])
    }
  }

  const handleFile = async (file: File) => {
    if (!file.type.startsWith('image/')) {
      alert('Please select an image file')
      return
    }

    setUploading(true)
    
    try {
      // For demo purposes, we'll create a data URL
      // In production, you would upload to a cloud service
      const reader = new FileReader()
      reader.onload = (e) => {
        const result = e.target?.result as string
        onChange(result, altText)
        setUploading(false)
      }
      reader.readAsDataURL(file)
    } catch (error) {
      console.error('Upload error:', error)
      alert('Failed to upload image')
      setUploading(false)
    }
  }

  const handleUrlSubmit = () => {
    if (urlInput.trim()) {
      onChange(urlInput.trim(), altText)
      setUrlInput('')
      setAltText('')
      setShowUrlInput(false)
    }
  }

  const removeImage = () => {
    onChange('', '')
    setAltText('')
    if (fileInputRef.current) {
      fileInputRef.current.value = ''
    }
  }

  return (
    <div className={className}>
      {label && (
        <label className="block text-sm font-medium text-muted-foreground700 mb-2">
          {label}
        </label>
      )}

      {value ? (
        <div className="relative">
          <div className="relative group">
            <img
              src={value}
              alt={altText || 'Featured image'}
              className="w-full h-48 object-cover rounded-lg border"
            />
            <div className="absolute inset-0 bg-black bg-opacity-50 opacity-0 group-hover:opacity-100 transition-opacity rounded-lg flex items-center justify-center">
              <button
                type="button"
                onClick={removeImage}
                className="bg-red-600 hover:bg-red-700 text-white p-2 rounded-full"
              >
                <X size={16} />
              </button>
            </div>
          </div>
          
          <div className="mt-2">
            <label className="block text-sm font-medium text-muted-foreground700 mb-1">
              Alt Text (for accessibility)
            </label>
            <input
              type="text"
              value={altText}
              onChange={(e) => setAltText(e.target.value)}
              placeholder="Describe the image..."
              className="w-full px-3 py-1 border border-border rounded-md bg-background text-foreground text-sm"
            />
          </div>
        </div>
      ) : (
        <div className="space-y-4">
          {/* File Upload Area */}
          <div
            className={`relative border-2 border-dashed rounded-lg p-6 text-center transition-colors ${
              dragActive
                ? 'border-blue-500 bg-blue-50'
                : uploading
                ? 'border-border bg-muted'
                : 'border-border hover:border-muted-foreground'
            }`}
            onDragEnter={handleDrag}
            onDragLeave={handleDrag}
            onDragOver={handleDrag}
            onDrop={handleDrop}
          >
            <input
              ref={fileInputRef}
              type="file"
              accept="image/*"
              onChange={handleFileSelect}
              className="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
              disabled={uploading}
            />
            
            <div className="space-y-2">
              <ImageIcon className="mx-auto h-12 w-12 text-muted-foreground400" />
              <div>
                <p className="text-sm text-muted-foreground600">
                  {uploading ? 'Uploading...' : 'Drop an image here, or click to select'}
                </p>
                <p className="text-xs text-muted-foreground500">PNG, JPG, GIF up to 10MB</p>
              </div>
            </div>
          </div>

          {/* URL Input Toggle */}
          <div className="text-center">
            <button
              type="button"
              onClick={() => setShowUrlInput(!showUrlInput)}
              className="text-sm text-blue-600 hover:text-blue-800 flex items-center justify-center space-x-1"
            >
              <Link2 size={14} />
              <span>Or use image URL</span>
            </button>
          </div>

          {/* URL Input */}
          {showUrlInput && (
            <div className="space-y-3 p-4 bg-gray-50 rounded-lg">
              <div>
                <label className="block text-sm font-medium text-muted-foreground700 mb-1">
                  Image URL
                </label>
                <input
                  type="url"
                  value={urlInput}
                  onChange={(e) => setUrlInput(e.target.value)}
                  placeholder="https://example.com/image.jpg"
                  className="w-full px-3 py-2 border border-border rounded-md bg-background text-foreground"
                />
              </div>
              
              <div>
                <label className="block text-sm font-medium text-muted-foreground700 mb-1">
                  Alt Text
                </label>
                <input
                  type="text"
                  value={altText}
                  onChange={(e) => setAltText(e.target.value)}
                  placeholder="Describe the image..."
                  className="w-full px-3 py-2 border border-border rounded-md bg-background text-foreground"
                />
              </div>
              
              <div className="flex justify-end space-x-2">
                <button
                  type="button"
                  onClick={() => setShowUrlInput(false)}
                  className="px-3 py-1 text-sm text-muted-foreground600 hover:text-muted-foreground800"
                >
                  Cancel
                </button>
                <button
                  type="button"
                  onClick={handleUrlSubmit}
                  disabled={!urlInput.trim()}
                  className="px-3 py-1 text-sm bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50"
                >
                  Add Image
                </button>
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  )
}
