"use client"

import * as React from "react"
import { Moon, Sun } from "lucide-react"
import { useTheme } from "next-themes"
import { Button } from "./button"

export function ThemeToggleSimpleNextThemes() {
  const { theme, setTheme } = useTheme()
  const [mounted, setMounted] = React.useState(false)

  React.useEffect(() => {
    setMounted(true)
  }, [])

  if (!mounted) {
    return (
      <Button variant="ghost" size="sm" disabled className="h-9 w-9 p-0">
        <Sun className="h-4 w-4" />
        <span className="sr-only">Toggle theme</span>
      </Button>
    )
  }

  const toggleTheme = () => {
    console.log('Current theme:', theme)
    if (theme === "light") {
      console.log('Switching to dark')
      setTheme("dark")
    } else {
      console.log('Switching to light')  
      setTheme("light")
    }
  }

  return (
    <Button
      variant="ghost"
      size="sm"
      onClick={toggleTheme}
      className="h-9 w-9 p-0 border"
      title={`Current theme: ${theme} - Click to toggle`}
    >
      {theme === "dark" ? (
        <Moon className="h-4 w-4 text-blue-400" />
      ) : (
        <Sun className="h-4 w-4 text-yellow-500" />
      )}
      <span className="sr-only">Toggle theme</span>
    </Button>
  )
}
