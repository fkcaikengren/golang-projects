export interface UserProfile {
  id: number
  email: string
  nickname: string
  status: string
}

export interface ApiEnvelope<T> {
  code: number
  message: string
  data: T
  request_id?: string
}

export interface AuthPayload {
  token: string
  user: UserProfile
}

export interface AdminUserProfile {
  id: number
  email: string
  display_name: string
  status: string
  last_login_at: number
}

export interface AdminAuthPayload {
  token: string
  admin_user: AdminUserProfile
}

export interface AdminQuickLink {
  label: string
  path: string
}

export interface AdminDashboardPayload {
  title: string
  admin_user: AdminUserProfile
  quick_links: AdminQuickLink[]
}

export interface ProblemSetSummary {
  id: number
  name: string
  slug: string
  description: string
  problemCount: number
}

export interface ProblemSummary {
  id: number
  slug: string
  title: string
  difficulty: 'easy' | 'medium' | 'hard'
  tags: string[]
  source: string
  status: 'unsolved' | 'attempted' | 'solved'
}

export interface ProblemDetail extends ProblemSummary {
  description: string
  inputDescription: string
  outputDescription: string
  sampleInput: string
  sampleOutput: string
  hint: string
}

export interface SubmissionItem {
  id: number
  problemTitle: string
  problemSlug: string
  language: string
  result: string
  submittedAt: string
}
