"use client"

import * as React from "react"
import { createContext, useContext, useEffect, useState } from "react"

type Theme = "dark" | "light" | "system"

type ThemeProviderProps = {
  children: React.ReactNode
  defaultTheme?: Theme
  storageKey?: string
  attribute?: string
  defaultValue?: string
  enableSystem?: boolean
  disableTransitionOnChange?: boolean
}

type ThemeProviderState = {
  theme: Theme
  setTheme: (theme: Theme) => void
  resolvedTheme: "dark" | "light"
}

const initialState: ThemeProviderState = {
  theme: "system",
  setTheme: () => null,
  resolvedTheme: "light",
}

const ThemeProviderContext = createContext<ThemeProviderState>(initialState)

export function ThemeProvider({
  children,
  defaultTheme = "system",
  storageKey = "webenable-ui-theme",
  attribute = "class",
  defaultValue = "light",
  enableSystem = true,
  disableTransitionOnChange = false,
  ...props
}: ThemeProviderProps) {
  const [theme, setThemeState] = useState<Theme>(defaultTheme)
  const [resolvedTheme, setResolvedTheme] = useState<"dark" | "light">("light")
  const [mounted, setMounted] = useState(false)

  // Initialize theme from localStorage and sessionStorage on mount
  useEffect(() => {
    setMounted(true)
    
    // Get stored theme preference
    const stored = localStorage?.getItem(storageKey) as Theme
    const resolvedFromSession = sessionStorage?.getItem('webenable-resolved-theme') as "dark" | "light"
    
    if (stored) {
      setThemeState(stored)
    }
    
    // Set resolved theme from session storage if available (prevents hydration mismatch)
    if (resolvedFromSession) {
      setResolvedTheme(resolvedFromSession)
    }
  }, [storageKey])

  // Apply theme to document when theme changes
  useEffect(() => {
    if (!mounted) return

    const root = window.document.documentElement

    // Remove existing theme classes
    root.classList.remove("light", "dark")

    let appliedTheme: "dark" | "light"

    if (theme === "system" && enableSystem) {
      const systemTheme = window.matchMedia("(prefers-color-scheme: dark)")
        .matches
        ? "dark"
        : "light"

      appliedTheme = systemTheme
      root.classList.add(systemTheme)
      setResolvedTheme(systemTheme)
    } else {
      appliedTheme = (theme || defaultValue) as "dark" | "light"
      root.classList.add(appliedTheme)
      setResolvedTheme(appliedTheme)
    }

    // Store resolved theme in session storage
    sessionStorage?.setItem('webenable-resolved-theme', appliedTheme)
  }, [theme, defaultValue, enableSystem, mounted])

  // Listen for system theme changes
  useEffect(() => {
    if (!mounted) return

    const handleSystemThemeChange = (e: MediaQueryListEvent) => {
      if (theme === "system") {
        const newSystemTheme = e.matches ? "dark" : "light"
        const root = window.document.documentElement
        
        root.classList.remove("light", "dark")
        root.classList.add(newSystemTheme)
        setResolvedTheme(newSystemTheme)
        sessionStorage?.setItem('webenable-resolved-theme', newSystemTheme)
      }
    }

    const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)")
    mediaQuery.addEventListener("change", handleSystemThemeChange)

    return () => mediaQuery.removeEventListener("change", handleSystemThemeChange)
  }, [theme, mounted])

  // Listen for theme changes from other tabs
  useEffect(() => {
    if (!mounted) return

    const handleStorageChange = (e: StorageEvent) => {
      if (e.key === storageKey && e.newValue) {
        const newTheme = e.newValue as Theme
        setThemeState(newTheme)
      }
    }

    window.addEventListener('storage', handleStorageChange)
    return () => window.removeEventListener('storage', handleStorageChange)
  }, [storageKey, mounted])

  const value = {
    theme,
    setTheme: (newTheme: Theme) => {
      if (mounted) {
        localStorage?.setItem(storageKey, newTheme)
      }
      setThemeState(newTheme)
    },
    resolvedTheme,
  }

  return (
    <ThemeProviderContext.Provider {...props} value={value}>
      {children}
    </ThemeProviderContext.Provider>
  )
}

export const useTheme = () => {
  const context = useContext(ThemeProviderContext)

  if (context === undefined)
    throw new Error("useTheme must be used within a ThemeProvider")

  return context
}
