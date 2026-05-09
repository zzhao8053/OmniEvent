export interface User {
  uid: number
  username: string
  email: string
  nickname: string
  avatar?: string
  avatar_url?: string
  avatar_provider?: string
  email_verified: boolean
  last_login_unix_time?: number
}

export interface UserBasicInfo {
  readonly username: string;
  readonly email: string;
  readonly nickname: string;
  readonly avatar: string;
  readonly avatar_provider?: string;
  readonly email_verified: boolean;
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
