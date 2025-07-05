// Theme initialization script to prevent flash of unstyled content
(function() {
  const storageKey = 'webenable-ui-theme'
  
  // Function to apply theme
  function applyTheme(theme) {
    const root = document.documentElement
    root.classList.remove('light', 'dark')
    
    if (theme === 'system') {
      const systemTheme = window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
      root.classList.add(systemTheme)
      // Store the resolved theme for SSR hydration
      sessionStorage.setItem('webenable-resolved-theme', systemTheme)
    } else {
      root.classList.add(theme)
      sessionStorage.setItem('webenable-resolved-theme', theme)
    }
  }
  
  // Get stored theme or default to system
  const theme = localStorage.getItem(storageKey) || 'system'
  
  // Apply theme immediately
  applyTheme(theme)
  
  // Listen for system theme changes if using system theme
  if (theme === 'system') {
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
    mediaQuery.addEventListener('change', function(e) {
      if (localStorage.getItem(storageKey) === 'system') {
        applyTheme('system')
      }
    })
  }
  
  // Listen for storage changes from other tabs
  window.addEventListener('storage', function(e) {
    if (e.key === storageKey) {
      const newTheme = e.newValue || 'system'
      applyTheme(newTheme)
    }
  })
})()
