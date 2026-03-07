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
  doc_type: 'reference' | 'typed'
  secure_notes?: string
  location_id?: number | null
  deleted_at?: string | null
  created_at: string
  updated_at: string
}

export interface DocumentFile {
  id: number
  document_id: number
  filename: string
  content_type: string
  file_size: number
  created_at: string
}

export interface AuthResponse {
  token: string
  user: User
}

export interface Contact {
  id: number
  user_id: number
  name: string
  relationship?: string
  role?: string
  phone?: string
  email?: string
  address?: string
  notes?: string
  is_primary?: boolean
  secure_notes?: string
  created_at: string
  updated_at: string
}

export interface Location {
  id: number
  user_id: number
  name: string
  type: 'physical' | 'digital'
  description?: string
  address?: string
  access_instructions?: string
  secure_notes?: string
  created_at: string
  updated_at: string
}

export interface DigitalAccess {
  id: number
  user_id: number
  type: 'computer' | 'phone' | 'password_manager'
  name: string
  username?: string
  instructions?: string
  location_id?: number | null
  secure_notes?: string
  created_at: string
  updated_at: string
}

export interface InsurancePolicy {
  id: number
  user_id: number
  party_id?: number | null
  provider: string
  policy_number?: string
  type?: string
  coverage_amount?: number | null
  beneficiary?: string
  agent_name?: string
  agent_phone?: string
  location_id?: number | null
  notes?: string
  secure_notes?: string
  created_at: string
  updated_at: string
}

export interface ServiceAccount {
  id: number
  user_id: number
  type: 'financial_tool' | 'backup_service' | 'tax_preparer'
  name: string
  provider?: string
  account_number?: string
  contact_name?: string
  contact_phone?: string
  contact_email?: string
  notes?: string
  secure_notes?: string
  created_at: string
  updated_at: string
}

export interface PartyObituaryInfo {
  id: number
  party_id: number
  type: 'survivor' | 'predeceased' | 'event'
  name: string
  relationship?: string
  details?: string
  event_date?: string | null
  created_at: string
  updated_at: string
}

export interface SurvivorLetterItem {
  id: number
  section_id: number
  source_type?: string | null
  source_id?: number | null
  content: string
  item_type: 'numbered' | 'bulleted' | 'paragraph'
  provenance: 'auto' | 'auto_edited' | 'manual'
  suppressed: boolean
  sort_order: number
  created_at: string
  updated_at: string
}

export interface SurvivorLetterSection {
  id: number
  letter_id: number
  section_key: string
  title: string
  sort_order: number
  visible: boolean
  items: SurvivorLetterItem[]
  created_at: string
  updated_at: string
}

export interface SurvivorLetter {
  id: number
  user_id: number
  greeting: string
  intro: string
  closing: string
  signature: string
  created_at: string
  updated_at: string
}

export interface FullSurvivorLetter extends SurvivorLetter {
  sections: SurvivorLetterSection[]
}

export interface ChecklistItem {
  category_id: number
  category_slug: string
  category_name: string
  status: 'pending' | 'complete' | 'not_applicable'
  has_document: boolean
}

export interface PartyChecklist {
  party_id: number
  party_name: string
  party_relationship: string
  items: ChecklistItem[]
}

export interface ConfidentialField {
  name: string
  value: string
}

export interface ConfidentialItem {
  label: string
  fields: ConfidentialField[]
}

export interface ConfidentialSection {
  title: string
  items: ConfidentialItem[]
}
