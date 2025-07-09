'use client';

import { useState } from 'react';
import { useContacts, useUpdateContactStatus, useReplyToContact } from '@/hooks/use-api';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Mail, Phone, Building, Calendar, User, MessageSquare } from 'lucide-react';
import { Contact } from '@/types/api';

export default function ContactsPage() {
  const { data: contacts, loading, error, refetch } = useContacts();
  const { mutate: updateStatus } = useUpdateContactStatus();
  const { mutate: replyToContact } = useReplyToContact();
  
  const [selectedContact, setSelectedContact] = useState<Contact | null>(null);
  const [replyDialogOpen, setReplyDialogOpen] = useState(false);
  const [replySubject, setReplySubject] = useState('');
  const [replyMessage, setReplyMessage] = useState('');
  const [isReplying, setIsReplying] = useState(false);
  const [replyError, setReplyError] = useState('');

  const handleStatusChange = async (contactId: string, status: string) => {
    try {
      await updateStatus({ id: contactId, status });
      refetch();
    } catch (error) {
      console.error('Failed to update status:', error);
    }
  };

  const handleReply = async () => {
    if (!selectedContact || !selectedContact.id) return;
    
    setIsReplying(true);
    setReplyError('');
    
    try {
      await replyToContact({
        id: selectedContact.id,
        subject: replySubject,
        message: replyMessage
      });
      
      setReplyDialogOpen(false);
      setReplySubject('');
      setReplyMessage('');
      setSelectedContact(null);
      refetch();
    } catch (error) {
      setReplyError('Failed to send reply. Please try again.');
    } finally {
      setIsReplying(false);
    }
  };

  const openReplyDialog = (contact: Contact) => {
    setSelectedContact(contact);
    setReplySubject(`Re: ${contact.subject}`);
    setReplyMessage(`Hi ${contact.name},\n\nThank you for contacting us. \n\nBest regards,\nWebEnable Team`);
    setReplyDialogOpen(true);
    setReplyError('');
  };

  const getStatusBadge = (status: string) => {
    const statusConfig = {
      new: { variant: 'default' as const, label: 'New' },
      read: { variant: 'secondary' as const, label: 'Read' },
      replied: { variant: 'outline' as const, label: 'Replied' }
    };
    
    const config = statusConfig[status as keyof typeof statusConfig] || statusConfig.new;
    return <Badge variant={config.variant}>{config.label}</Badge>;
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  if (loading) {
    return (
      <div className="container mx-auto p-6">
        <div className="text-center">Loading contacts...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="container mx-auto p-6">
        <Alert variant="destructive">
          <AlertDescription>
            Failed to load contacts: {error}
          </AlertDescription>
        </Alert>
      </div>
    );
  }

  return (
    <div className="container mx-auto p-6">
      <div className="mb-6">
        <h1 className="text-3xl font-bold">Contact Management</h1>
        <p className="text-gray-600 mt-2">Manage and reply to customer inquiries</p>
      </div>

      <div className="grid gap-6">
        {contacts && contacts.length > 0 ? (
          contacts.map((contact) => (
            <Card key={contact.id} className="w-full">
              <CardHeader>
                <div className="flex justify-between items-start">
                  <div className="space-y-2">
                    <CardTitle className="flex items-center gap-2">
                      <User className="h-5 w-5" />
                      {contact.name}
                      {getStatusBadge(contact.status)}
                    </CardTitle>
                    <CardDescription className="text-lg font-medium">
                      {contact.subject}
                    </CardDescription>
                  </div>
                  <div className="flex gap-2">
                    <Select 
                      value={contact.status} 
                      onValueChange={(status) => contact.id && handleStatusChange(contact.id, status)}
                    >
                      <SelectTrigger className="w-32">
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="new">New</SelectItem>
                        <SelectItem value="read">Read</SelectItem>
                        <SelectItem value="replied">Replied</SelectItem>
                      </SelectContent>
                    </Select>
                    
                    <Button 
                      onClick={() => openReplyDialog(contact)}
                      variant="outline"
                      size="sm"
                      className="gap-2"
                    >
                      <MessageSquare className="h-4 w-4" />
                      Reply
                    </Button>
                  </div>
                </div>
              </CardHeader>
              
              <CardContent className="space-y-4">
                <div className="grid md:grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <div className="flex items-center gap-2 text-sm">
                      <Mail className="h-4 w-4 text-gray-500" />
                      <span className="font-medium">Email:</span>
                      <a href={`mailto:${contact.email}`} className="text-blue-600 hover:underline">
                        {contact.email}
                      </a>
                    </div>
                    
                    {contact.phone && (
                      <div className="flex items-center gap-2 text-sm">
                        <Phone className="h-4 w-4 text-gray-500" />
                        <span className="font-medium">Phone:</span>
                        <a href={`tel:${contact.phone}`} className="text-blue-600 hover:underline">
                          {contact.phone}
                        </a>
                      </div>
                    )}
                    
                    {contact.company && (
                      <div className="flex items-center gap-2 text-sm">
                        <Building className="h-4 w-4 text-gray-500" />
                        <span className="font-medium">Company:</span>
                        <span>{contact.company}</span>
                      </div>
                    )}
                  </div>
                  
                  <div className="space-y-2">
                    <div className="flex items-center gap-2 text-sm">
                      <Calendar className="h-4 w-4 text-gray-500" />
                      <span className="font-medium">Received:</span>
                      <span>{formatDate(contact.created_at)}</span>
                    </div>
                    
                    {contact.replied_at && (
                      <div className="flex items-center gap-2 text-sm">
                        <MessageSquare className="h-4 w-4 text-gray-500" />
                        <span className="font-medium">Replied:</span>
                        <span>{formatDate(contact.replied_at)}</span>
                      </div>
                    )}
                  </div>
                </div>
                
                <div className="border-t pt-4">
                  <div className="font-medium text-sm mb-2">Message:</div>
                  <div className="bg-gray-50 p-4 rounded-md text-sm whitespace-pre-wrap">
                    {contact.message}
                  </div>
                </div>
              </CardContent>
            </Card>
          ))
        ) : (
          <div className="text-center py-8">
            <p className="text-gray-500">No contacts found.</p>
          </div>
        )}
      </div>

      {/* Reply Dialog */}
      <Dialog open={replyDialogOpen} onOpenChange={setReplyDialogOpen}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>Reply to {selectedContact?.name}</DialogTitle>
            <DialogDescription>
              Send an email reply to {selectedContact?.email}
            </DialogDescription>
          </DialogHeader>
          
          <div className="space-y-4">
            {replyError && (
              <Alert variant="destructive">
                <AlertDescription>{replyError}</AlertDescription>
              </Alert>
            )}
            
            <div className="space-y-2">
              <Label htmlFor="subject">Subject</Label>
              <Input
                id="subject"
                value={replySubject}
                onChange={(e) => setReplySubject(e.target.value)}
                placeholder="Email subject"
              />
            </div>
            
            <div className="space-y-2">
              <Label htmlFor="message">Message</Label>
              <Textarea
                id="message"
                value={replyMessage}
                onChange={(e) => setReplyMessage(e.target.value)}
                placeholder="Your reply message..."
                rows={10}
                className="resize-none"
              />
            </div>
            
            <div className="flex justify-end gap-2">
              <Button 
                variant="outline" 
                onClick={() => setReplyDialogOpen(false)}
                disabled={isReplying}
              >
                Cancel
              </Button>
              <Button 
                onClick={handleReply}
                disabled={isReplying || !replySubject.trim() || !replyMessage.trim()}
              >
                {isReplying ? 'Sending...' : 'Send Reply'}
              </Button>
            </div>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
}
