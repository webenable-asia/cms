'use client'

import { useState, useEffect } from 'react'
import { useRouter, useParams } from 'next/navigation'
import Link from 'next/link'
import { api } from '../../../../../lib/api'
import { ThemeToggle } from '@/components/ui/theme-toggle'

interface Contact {
  id: string
  name: string
  email: string
  company: string
  phone: string
  subject: string
  message: string
  status: string
  created_at: string
}

export default function ReplyToContact() {
  const [contact, setContact] = useState<Contact | null>(null)
  const [replyData, setReplyData] = useState({
    subject: '',
    message: ''
  })
  const [loading, setLoading] = useState(false)
  const [fetching, setFetching] = useState(true)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState(false)
  const router = useRouter()
  const params = useParams()

  useEffect(() => {
    const token = localStorage.getItem('token')
    if (!token) {
      router.push('/admin')
      return
    }

    fetchContact()
  }, [router, params.id])

  const fetchContact = async () => {
    try {
      const response = await api.get(`/contacts/${params.id}`)
      const contactData = response.data
      setContact(contactData)
      
      // Pre-fill the reply subject
      setReplyData({
        subject: `Re: ${contactData.subject}`,
        message: `Dear ${contactData.name},\n\nThank you for contacting WebEnable. \n\n\n\nBest regards,\nWebEnable Team`
      })
    } catch (error) {
      console.error('Error fetching contact:', error)
      setError('Failed to load contact')
    } finally {
      setFetching(false)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setError('')

    try {
      await api.post(`/contacts/${params.id}/reply`, {
        to_email: contact?.email,
        to_name: contact?.name,
        subject: replyData.subject,
        message: replyData.message,
        original_message: contact?.message
      })
      
      // Mark contact as replied
      await api.put(`/contacts/${params.id}`, { status: 'replied' })
      
      setSuccess(true)
    } catch (error) {
      setError('Failed to send reply')
    } finally {
      setLoading(false)
    }
  }

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    setReplyData({
      ...replyData,
      [e.target.name]: e.target.value
    })
  }

  if (fetching) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-xl">Loading contact...</div>
      </div>
    )
  }

  if (error && !contact) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-foreground mb-4">Contact Not Found</h1>
          <Link 
            href="/admin/dashboard"
            className="bg-primary hover:bg-primary/90 text-primary-foreground font-bold py-2 px-4 rounded"
          >
            Back to Dashboard
          </Link>
        </div>
      </div>
    )
  }

  if (success) {
    return (
      <div className="min-h-screen bg-background">
        <nav className="bg-card shadow-sm border-b border-border">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="flex justify-between h-16">
              <div className="flex items-center">
                <Link href="/admin/dashboard" className="text-primary hover:text-primary/80">
                  ← Back to Dashboard
                </Link>
              </div>
              <div className="flex items-center">
                <ThemeToggle />
              </div>
            </div>
          </div>
        </nav>

        <div className="flex items-center justify-center min-h-screen bg-background">
          <div className="text-center max-w-md mx-auto px-4">
            <div className="w-24 h-24 bg-green-100 dark:bg-green-900 rounded-full flex items-center justify-center mx-auto mb-6">
              <svg className="w-12 h-12 text-green-600 dark:text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
              </svg>
            </div>
            <h1 className="text-3xl font-bold text-foreground mb-4">Reply Sent!</h1>
            <p className="text-muted-foreground mb-8">
              Your reply has been sent to {contact?.email} and the contact has been marked as replied.
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <Link 
                href="/admin/dashboard"
                className="bg-primary hover:bg-primary/90 text-primary-foreground font-bold py-3 px-6 rounded-lg transition duration-300"
              >
                Back to Dashboard
              </Link>
              <Link 
                href={`/admin/contacts/${params.id}/reply`}
                className="border border-border text-foreground hover:bg-muted font-bold py-3 px-6 rounded-lg transition duration-300"
                onClick={() => setSuccess(false)}
              >
                Send Another Reply
              </Link>
            </div>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-background">
      <nav className="bg-card shadow-sm border-b border-border">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex items-center">
              <Link href="/admin/dashboard" className="text-primary hover:text-primary/80">
                ← Back to Dashboard
              </Link>
            </div>
            <div className="flex items-center">
              <ThemeToggle />
            </div>
          </div>
        </div>
      </nav>

      <div className="max-w-4xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="px-4 py-6 sm:px-0">
          <h1 className="text-3xl font-bold text-foreground mb-8">Reply to Contact</h1>

          {/* Original Message */}
          <div className="bg-card rounded-xl shadow-lg border border-border p-6 mb-8">
            <h2 className="text-xl font-bold text-foreground mb-4">Original Message</h2>
            <div className="grid md:grid-cols-2 gap-4 mb-4">
              <div>
                <span className="font-medium text-foreground">From:</span> {contact?.name}
              </div>
              <div>
                <span className="font-medium text-foreground">Email:</span> {contact?.email}
              </div>
              {contact?.company && (
                <div>
                  <span className="font-medium text-foreground">Company:</span> {contact?.company}
                </div>
              )}
              {contact?.phone && (
                <div>
                  <span className="font-medium text-foreground">Phone:</span> {contact?.phone}
                </div>
              )}
            </div>
            <div className="mb-4">
              <span className="font-medium text-foreground">Subject:</span> {contact?.subject}
            </div>
            <div className="mb-4">
              <span className="font-medium text-foreground">Date:</span>{' '}
              {contact?.created_at && new Date(contact.created_at).toLocaleString()}
            </div>
            <div className="border-t pt-4">
              <span className="font-medium text-foreground">Message:</span>
              <p className="mt-2 text-muted-foreground whitespace-pre-wrap">{contact?.message}</p>
            </div>
          </div>

          {/* Reply Form */}
          <div className="bg-card rounded-xl shadow-lg border border-border p-6">
            <h2 className="text-xl font-bold text-foreground mb-6">Your Reply</h2>
            
            <form onSubmit={handleSubmit} className="space-y-6">
              <div>
                <label htmlFor="subject" className="block text-sm font-medium text-foreground mb-2">
                  Subject
                </label>
                <input
                  type="text"
                  id="subject"
                  name="subject"
                  required
                  className="w-full px-4 py-3 border border-border rounded-lg focus:outline-none focus:ring-2 focus:ring-ring focus:border-ring bg-background text-foreground placeholder:text-muted-foreground"
                  value={replyData.subject}
                  onChange={handleChange}
                />
              </div>

              <div>
                <label htmlFor="message" className="block text-sm font-medium text-foreground mb-2">
                  Message
                </label>
                <textarea
                  id="message"
                  name="message"
                  rows={12}
                  required
                  className="w-full px-4 py-3 border border-border rounded-lg focus:outline-none focus:ring-2 focus:ring-ring focus:border-ring bg-background text-foreground placeholder:text-muted-foreground"
                  value={replyData.message}
                  onChange={handleChange}
                  placeholder="Type your reply here..."
                />
              </div>

              <div className="bg-primary/10 border border-primary/20 rounded-lg p-4">
                <div className="flex">
                  <svg className="flex-shrink-0 w-5 h-5 text-primary" fill="currentColor" viewBox="0 0 20 20">
                    <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clipRule="evenodd" />
                  </svg>
                  <div className="ml-3">
                    <p className="text-sm text-foreground">
                      <strong>Note:</strong> This will send an email reply to {contact?.email} and automatically mark this contact as "replied" in your dashboard.
                    </p>
                  </div>
                </div>
              </div>

              {error && (
                <div className="text-destructive text-sm">{error}</div>
              )}

              <div className="flex justify-end space-x-4">
                <Link
                  href="/admin/dashboard"
                  className="px-6 py-3 border border-border rounded-lg text-foreground hover:bg-muted transition duration-300"
                >
                  Cancel
                </Link>
                <button
                  type="submit"
                  disabled={loading}
                  className="px-6 py-3 bg-primary text-primary-foreground rounded-lg hover:bg-primary/90 disabled:opacity-50 transition duration-300"
                >
                  {loading ? 'Sending Reply...' : 'Send Reply'}
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
  )
}
