import axios, { AxiosError, type AxiosRequestConfig, type AxiosRequestHeaders, type AxiosResponse } from 'axios'
import { userState } from './userstate'
import { ApiException, type ApiResponse } from './api'

const BASE_URL = '/api'

export type ApiResponsePromise<T> = Promise<AxiosResponse<ApiResponse<T>>>

interface ApiRequestConfig extends AxiosRequestConfig {
  readonly headers: AxiosRequestHeaders
  readonly noAuth?: boolean
  readonly ignoreBlocked?: boolean
  readonly ignoreError?: boolean
  readonly cancelableUuid?: string
}

let needBlockRequest = false
const blockedRequests: ((token: string | undefined) => void)[] = []
const cancelableRequests: Record<string, boolean> = {}

const axiosInstance = axios.create({
  baseURL: BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

axiosInstance.interceptors.request.use(
  (config: ApiRequestConfig) => {
    const token = userState.getToken()
    if (token && !config.noAuth) {
      config.headers.Authorization = `Bearer ${token}`
    }

    if (needBlockRequest && !config.ignoreBlocked) {
      return new Promise(resolve => {
        blockedRequests.push(newToken => {
          if (newToken) {
            config.headers.Authorization = `Bearer ${newToken}`
          }
          resolve(config)
        })
      })
    }

    return config
  },
  (error) => Promise.reject(error)
)

axiosInstance.interceptors.response.use(
  (response) => {
    if ('cancelableUuid' in response.config && response.config.cancelableUuid && cancelableRequests[response.config.cancelableUuid as string]) {
      delete cancelableRequests[response.config.cancelableUuid as string]
      return Promise.reject({ canceled: true })
    }
    return response
  },
  (error: AxiosError<ApiResponse>) => {
    const config = error.response?.config as ApiRequestConfig | undefined
    if (config && 'cancelableUuid' in config && config.cancelableUuid && cancelableRequests[config.cancelableUuid as string]) {
      delete cancelableRequests[config.cancelableUuid as string]
      return Promise.reject({ canceled: true })
    }

    if (error.response && config && !config.ignoreError && error.response.data && error.response.data.errorCode) {
      const errorCode = error.response.data.errorCode
      if (errorCode === 202001 || errorCode === 202002 || errorCode === 202003 ||
          errorCode === 202004 || errorCode === 202005 || errorCode === 202006 ||
          errorCode === 202012) {
        userState.clearToken()
        window.location.href = '/login'
        return Promise.reject({ processed: true })
      }
    }

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

function handleResponse<T>(response: AxiosResponse<ApiResponse<T>>): T {
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

  async authorize2FA(passcode: string, token: string): Promise<LoginResponse> {
    const response = await axiosInstance.post<ApiResponse<LoginResponse>>('2fa/authorize.json', { passcode }, {
      noAuth: true,
      headers: {
        Authorization: `Bearer ${token}`
      }
    } as ApiRequestConfig)
    return handleResponse(response)
  }

  async authorize2FAByBackupCode(recoveryCode: string, token: string): Promise<LoginResponse> {
    const response = await axiosInstance.post<ApiResponse<LoginResponse>>('2fa/recovery.json', { recoveryCode }, {
      noAuth: true,
      headers: {
        Authorization: `Bearer ${token}`
      }
    } as ApiRequestConfig)
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

  async refreshTokenBlocking(): Promise<string | undefined> {
    needBlockRequest = true
    try {
      const response = await axiosInstance.post<ApiResponse<{ newToken: string }>>('/v1/tokens/refresh.json', {}, {
        ignoreBlocked: true
      } as ApiRequestConfig)
      const newToken = response.data.result?.newToken
      blockedRequests.forEach(func => func(newToken))
      blockedRequests.length = 0
      return newToken
    } finally {
      needBlockRequest = false
    }
  }
}

export const userService = new UserService()
export const tokenService = new TokenService()

export const apiService = {
  setLocale: (locale: string) => {
    axiosInstance.defaults.headers.common['Accept-Language'] = locale
  },
  cancelRequest: (cancelableUuid: string) => {
    cancelableRequests[cancelableUuid] = true
  }
}