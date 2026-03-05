export interface User {
  id: number
  email: string
  name: string
}

export interface Party {
  id: number
  user_id: number
  name: string
  relationship: 'self' | 'spouse' | 'dependent' | 'other'
  notes?: string
  created_at: string
  updated_at: string
}

export interface DocumentCategory {
  id: number
  slug: string
  name: string
  description?: string
  sort_order: number
}

export interface Document {
  id: number
  party_id: number
  category_id: number
  title: string
  content?: string
  status: 'draft' | 'complete'
  created_at: string
  updated_at: string
}

export interface AuthResponse {
  token: string
  user: User
}
