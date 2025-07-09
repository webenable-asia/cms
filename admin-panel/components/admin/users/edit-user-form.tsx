'use client'

import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Switch } from '@/components/ui/switch'
import { User } from '@/types/api'
import { usersApi } from '@/lib/api'

interface EditUserFormProps {
  user: User
  onSuccess: () => void
  onCancel: () => void
}

export function EditUserForm({ user, onSuccess, onCancel }: EditUserFormProps) {
  const [formData, setFormData] = useState({
    username: user.username,
    email: user.email,
    password: '',
    role: user.role,
    active: user.active
  })
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!formData.username || !formData.email || !formData.role) {
      setError('Username, email, and role are required')
      return
    }

    try {
      setLoading(true)
      setError('')
      
      const updates: any = {
        username: formData.username,
        email: formData.email,
        role: formData.role,
        active: formData.active
      }
      
      // Only include password if it's provided
      if (formData.password) {
        updates.password = formData.password
      }
      
      await usersApi.update(user.id!, updates)
      onSuccess()
    } catch (err: any) {
      setError(err.message || 'Failed to update user')
    } finally {
      setLoading(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {error && (
        <div className="text-sm text-red-600 bg-red-50 p-3 rounded-md">
          {error}
        </div>
      )}
      
      <div className="space-y-2">
        <Label htmlFor="edit-username">Username</Label>
        <Input
          id="edit-username"
          value={formData.username}
          onChange={(e) => setFormData({ ...formData, username: e.target.value })}
          placeholder="Enter username"
          required
        />
      </div>

      <div className="space-y-2">
        <Label htmlFor="edit-email">Email</Label>
        <Input
          id="edit-email"
          type="email"
          value={formData.email}
          onChange={(e) => setFormData({ ...formData, email: e.target.value })}
          placeholder="Enter email address"
          required
        />
      </div>

      <div className="space-y-2">
        <Label htmlFor="edit-password">Password</Label>
        <Input
          id="edit-password"
          type="password"
          value={formData.password}
          onChange={(e) => setFormData({ ...formData, password: e.target.value })}
          placeholder="Leave empty to keep current password"
        />
        <p className="text-xs text-muted-foreground">
          Leave empty to keep the current password
        </p>
      </div>

      <div className="space-y-2">
        <Label htmlFor="edit-role">Role</Label>
        <Select value={formData.role} onValueChange={(value: string) => setFormData({ ...formData, role: value as 'admin' | 'editor' | 'author' })}>
          <SelectTrigger>
            <SelectValue placeholder="Select a role" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="admin">Admin</SelectItem>
            <SelectItem value="editor">Editor</SelectItem>
            <SelectItem value="author">Author</SelectItem>
          </SelectContent>
        </Select>
      </div>

      <div className="flex items-center space-x-2">
        <Switch
          id="edit-active"
          checked={formData.active}
          onCheckedChange={(checked) => setFormData({ ...formData, active: checked })}
        />
        <Label htmlFor="edit-active">Active</Label>
      </div>

      <div className="flex justify-end space-x-2 pt-4">
        <Button type="button" variant="outline" onClick={onCancel}>
          Cancel
        </Button>
        <Button type="submit" disabled={loading}>
          {loading ? 'Updating...' : 'Update User'}
        </Button>
      </div>
    </form>
  )
}
