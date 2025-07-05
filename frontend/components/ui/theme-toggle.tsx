"use client"

import * as React from "react"
import { Moon, Sun, Monitor, ChevronDown } from "lucide-react"
import { useTheme } from "./theme-provider"
import { Button } from "./button"

export function ThemeToggle() {
  const { theme, setTheme, resolvedTheme } = useTheme()
  const [mounted, setMounted] = React.useState(false)
  const [isOpen, setIsOpen] = React.useState(false)

  React.useEffect(() => {
    setMounted(true)
  }, [])

  React.useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (isOpen && !(event.target as Element).closest('.theme-toggle-dropdown')) {
        setIsOpen(false)
      }
    }

    document.addEventListener('click', handleClickOutside)
    return () => document.removeEventListener('click', handleClickOutside)
  }, [isOpen])

  if (!mounted) {
    return (
      <Button variant="ghost" size="sm" disabled className="h-9 w-9 p-0">
        <Sun className="h-4 w-4" />
        <span className="sr-only">Toggle theme</span>
      </Button>
    )
  }

  const getIcon = () => {
    if (theme === "system") {
      return <Monitor className="h-4 w-4" />
    }
    return resolvedTheme === "dark" ? (
      <Moon className="h-4 w-4" />
    ) : (
      <Sun className="h-4 w-4" />
    )
  }

  const options = [
    { value: "light", label: "Light", icon: Sun },
    { value: "dark", label: "Dark", icon: Moon },
    { value: "system", label: "System", icon: Monitor },
  ]

  return (
    <div className="relative theme-toggle-dropdown">
      <Button
        variant="ghost"
        size="sm"
        onClick={() => setIsOpen(!isOpen)}
        className="h-9 px-2 py-0 gap-1 text-foreground/80 hover:text-foreground"
        aria-expanded={isOpen}
        aria-haspopup="true"
      >
        {getIcon()}
        <ChevronDown className={`h-3 w-3 transition-transform ${isOpen ? 'rotate-180' : ''}`} />
        <span className="sr-only">Toggle theme</span>
      </Button>

      {isOpen && (
        <div className="absolute right-0 top-full mt-1 w-36 rounded-md border bg-popover p-1 shadow-md z-50 animate-in fade-in-0 zoom-in-95">
          {options.map((option) => {
            const Icon = option.icon
            const isSelected = theme === option.value
            
            return (
              <button
                key={option.value}
                onClick={() => {
                  setTheme(option.value as "light" | "dark" | "system")
                  setIsOpen(false)
                }}
                className={`
                  flex w-full items-center gap-2 rounded-sm px-2 py-1.5 text-sm
                  hover:bg-accent hover:text-accent-foreground
                  ${isSelected 
                    ? 'bg-accent text-accent-foreground font-medium' 
                    : 'text-popover-foreground'
                  }
                `}
              >
                <Icon className="h-4 w-4" />
                <span>{option.label}</span>
                {option.value === "system" && (
                  <span className="ml-auto text-xs text-muted-foreground">
                    ({resolvedTheme})
                  </span>
                )}
              </button>
            )
          })}
        </div>
      )}
    </div>
  )
}

// Keep a simple cycle version for minimal use cases
export function ThemeToggleSimple() {
  const { theme, setTheme, resolvedTheme } = useTheme()
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

  const cycleTheme = () => {
    console.log('Current theme:', theme, 'Resolved theme:', resolvedTheme)
    if (theme === "light") {
      console.log('Switching to dark')
      setTheme("dark")
    } else if (theme === "dark") {
      console.log('Switching to system')
      setTheme("system")
    } else {
      console.log('Switching to light')
      setTheme("light")
    }
  }

  const getIcon = () => {
    if (theme === "system") {
      return <Monitor className="h-4 w-4" />
    }
    return resolvedTheme === "dark" ? (
      <Moon className="h-4 w-4" />
    ) : (
      <Sun className="h-4 w-4" />
    )
  }

  return (
    <Button
      variant="ghost"
      size="sm"
      onClick={cycleTheme}
      className="h-9 w-9 p-0"
      title={`Current: ${theme === "system" ? `System (${resolvedTheme})` : theme}`}
    >
      {getIcon()}
      <span className="sr-only">Toggle theme</span>
    </Button>
  )
}

// Remove the old ThemeToggleDropdown - it's replaced by the new ThemeToggle
export const ThemeToggleDropdown = ThemeToggle
