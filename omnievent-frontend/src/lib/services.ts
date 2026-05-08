import axios, { AxiosError, type AxiosRequestConfig } from 'axios'
import { userState } from './userstate'
import { ApiException, type ApiResponse } from './api'

const BASE_URL = '/api'

const axiosInstance = axios.create({
  baseURL: BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

axiosInstance.interceptors.request.use(
  (config) => {
    const token = userState.getToken()
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

axiosInstance.interceptors.response.use(
  (response) => response,
  (error: AxiosError<ApiResponse>) => {
    if (error.response) {
      const data = error.response.data
      if (data?.errorCode) {
        throw new ApiException(data.errorCode, data.errorMessage || 'Unknown error', data.path)
      }
      if (error.response.status === 401) {
        userState.clearToken()
        window.location.href = '/login'
        return Promise.reject(new ApiException(0, 'Unauthorized'))
      }
    }
    return Promise.reject(error)
  }
)

function handleResponse<T>(response: axios.AxiosResponse<ApiResponse<T>>): T {
  if (response.data.success && response.data.result !== undefined) {
    return response.data.result
  }
  throw new ApiException(response.data.errorCode || 0, response.data.errorMessage || 'Unknown error')
}

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  need_2fa: boolean
  user: {
    uid: number
    username: string
    email: string
    nickname: string
  }
}

export interface RegisterRequest {
  username: string
  email: string
  nickname: string
  password: string
}

export interface ProfileResponse {
  uid: number
  username: string
  email: string
  nickname: string
  avatar_url?: string
}

export interface TokenResponse {
  uid: number
  user_token_id: number
  token_type: number
  user_agent: string
  created_unix_time: number
  expired_unix_time: number
  last_seen_unix_time: number
}

export interface TokenListResponse {
  tokens: TokenResponse[]
}

class UserService {
  async login(request: LoginRequest): Promise<LoginResponse> {
    const response = await axiosInstance.post<ApiResponse<LoginResponse>>('/authorize.json', request)
    return handleResponse(response)
  }

  async register(request: RegisterRequest): Promise<void> {
    const response = await axiosInstance.post('/register.json', request)
    return handleResponse(response)
  }

  async getProfile(): Promise<ProfileResponse> {
    const response = await axiosInstance.get<ApiResponse<ProfileResponse>>('/v1/users/profile/get.json')
    return handleResponse(response)
  }

  async updateProfile(data: { nickname: string; email: string }): Promise<ProfileResponse> {
    const response = await axiosInstance.post<ApiResponse<ProfileResponse>>('/v1/users/profile/update.json', data)
    return handleResponse(response)
  }

  async updateAvatar(avatarType: string): Promise<{ avatar_url: string }> {
    const response = await axiosInstance.post<ApiResponse<{ avatar_url: string }>>('/v1/users/avatar/update.json', { avatar_type: avatarType })
    return handleResponse(response)
  }

  async removeAvatar(): Promise<void> {
    const response = await axiosInstance.post('/v1/users/avatar/remove.json')
    return handleResponse(response)
  }
}

class TokenService {
  async listTokens(): Promise<TokenListResponse> {
    const response = await axiosInstance.get<ApiResponse<TokenListResponse>>('/v1/tokens/list.json')
    return handleResponse(response)
  }

  async refreshToken(userTokenId: number): Promise<{ token: string }> {
    const response = await axiosInstance.post<ApiResponse<{ token: string }>>('/v1/tokens/refresh.json', { user_token_id: userTokenId })
    return handleResponse(response)
  }

  async revokeToken(userTokenId: number): Promise<void> {
    const response = await axiosInstance.post('/v1/tokens/revoke.json', { user_token_id: userTokenId })
    return handleResponse(response)
  }

  async revokeAllTokens(): Promise<void> {
    const response = await axiosInstance.post('/v1/tokens/revoke_all.json')
    return handleResponse(response)
  }
}

export const userService = new UserService()
export const tokenService = new TokenService()
