import React from 'react'

export const Avatar = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className = '', ...props }, ref) => (
  <div
    ref={ref}
    className={`avatar ${className}`}
    {...props}
  />
))
Avatar.displayName = 'Avatar'

export const AvatarFallback = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className = '', ...props }, ref) => (
  <div
    ref={ref}
    className={`avatar-fallback ${className}`}
    {...props}
  />
))
AvatarFallback.displayName = 'AvatarFallback'
