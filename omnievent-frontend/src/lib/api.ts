export interface ApiResponse<T = any> {
  success: boolean
  result?: T
  errorCode?: number
  errorMessage?: string
  path?: string
}

export interface ApiError {
  code: number
  message: string
  path?: string
}

export class ApiException extends Error {
  constructor(
    public code: number,
    public message: string,
    public path?: string
  ) {
    super(message)
    this.name = 'ApiException'
  }
}

export function isApiError(error: any): error is ApiException {
  return error instanceof ApiException
}

export function getErrorMessage(error: any): string {
  if (isApiError(error)) {
    return error.message
  }
  if (error.response?.data?.errorMessage) {
    return error.response.data.errorMessage
  }
  if (error.message) {
    return error.message
  }
  return 'An unknown error occurred'
}

export function getErrorCode(error: any): number {
  if (isApiError(error)) {
    return error.code
  }
  if (error.response?.data?.errorCode) {
    return error.response.data.errorCode
  }
  return 0
}

export const ErrorCodes = {
  // User errors (2000-2999)
  USER_NOT_FOUND: 2001,
  USER_EXISTS: 2002,
  INVALID_CREDENTIALS: 2003,
  EMAIL_EXISTS: 2004,
  USERNAME_EXISTS: 2005,
  USER_DISABLED: 2006,
  EMAIL_NOT_VERIFIED: 2007,

  // Token errors (3000-3999)
  TOKEN_NOT_FOUND: 3001,
  TOKEN_EXPIRED: 3002,
  TOKEN_INVALID: 3003,
  TOKEN_TYPE_MISMATCH: 3004,
  REFRESH_TOKEN_EXPIRED: 3005,

  // Validation errors (4000-4999)
  VALIDATION_FAILED: 4001,
  INVALID_REQUEST: 4002,

  // System errors (100000-999999)
  SYSTEM_ERROR: 100000,
  DATABASE_ERROR: 100001,
}

export function isUserError(code: number): boolean {
  return code >= 2000 && code < 3000
}

export function isTokenError(code: number): boolean {
  return code >= 3000 && code < 4000
}

export function isValidationError(code: number): boolean {
  return code >= 4000 && code < 5000
}

export function isSystemError(code: number): boolean {
  return code >= 100000
}
