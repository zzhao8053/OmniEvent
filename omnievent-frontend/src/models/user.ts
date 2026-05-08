export interface User {
  uid: number
  username: string
  email: string
  nickname: string
  avatar_url?: string
  email_verified: boolean
  last_login_unix_time?: number
}

export interface LoginRequest {
  username: string
  password: string
}

export interface RegisterRequest {
  username: string
  email: string
  nickname: string
  password: string
}

export interface TokenRecord {
  uid: number
  user_token_id: number
  token_type: number
  user_agent: string
  created_unix_time: number
  expired_unix_time: number
  last_seen_unix_time: number
}
