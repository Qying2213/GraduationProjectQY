export interface User {
    id: number
    username: string
    email: string
    role: 'admin' | 'hr' | 'candidate'
    avatar?: string
    phone?: string
    real_name?: string
    status: string
    created_at: string
    updated_at: string
}

export interface Talent {
    id: number
    name: string
    email: string
    phone: string
    skills: string[]
    experience: number
    education: string
    status: 'active' | 'hired' | 'rejected' | 'pending'
    tags: string[]
    resume_id?: number
    location: string
    salary: string
    summary: string
    user_id: number
    created_at: string
    updated_at: string
}

export interface Job {
    id: number
    title: string
    description: string
    requirements: string[]
    salary: string
    location: string
    type: 'full-time' | 'part-time' | 'contract' | 'internship'
    status: 'open' | 'closed' | 'filled'
    created_by: number
    department: string
    level: string
    skills: string[]
    benefits: string[]
    created_at: string
    updated_at: string
}

export interface Resume {
    id: number
    talent_id: number
    file_name: string
    file_url: string
    file_size: number
    parsed_data: string
    status: string
    created_at: string
    updated_at: string
}

export interface Application {
    id: number
    job_id: number
    talent_id: number
    resume_id: number
    status: 'pending' | 'reviewed' | 'interview' | 'rejected' | 'accepted'
    cover_letter: string
    notes: string
    created_at: string
    updated_at: string
}

export interface Message {
    id: number
    from_id: number
    to_id: number
    title: string
    content: string
    type: 'system' | 'user' | 'notification'
    is_read: boolean
    related_id?: number
    created_at: string
    updated_at: string
}

export interface Recommendation {
    id: number
    name: string
    score: number
    reason: string
    match_level: 'high' | 'medium' | 'low'
}

export interface LoginRequest {
    username: string
    password: string
}

export interface RegisterRequest {
    username: string
    email: string
    password: string
    role?: 'hr' | 'candidate'
    real_name?: string
    phone?: string
}

export interface ApiResponse<T = any> {
    code: number
    message: string
    data?: T
}

export interface PaginatedResponse<T> {
    total: number
    page: number
    page_size: number
    [key: string]: any
}
