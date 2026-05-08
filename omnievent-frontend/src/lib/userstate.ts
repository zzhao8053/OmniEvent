const TOKEN_KEY = 'omnievent_token'

class UserState {
  private encryptedToken: string | null = null

  getToken(): string | null {
    if (this.encryptedToken) {
      return this.encryptedToken
    }

    if (typeof localStorage !== 'undefined') {
      this.encryptedToken = localStorage.getItem(TOKEN_KEY)
      return this.encryptedToken
    }

    return null
  }

  setToken(token: string): void {
    this.encryptedToken = token
    if (typeof localStorage !== 'undefined') {
      localStorage.setItem(TOKEN_KEY, token)
    }
  }

  clearToken(): void {
    this.encryptedToken = null
    if (typeof localStorage !== 'undefined') {
      localStorage.removeItem(TOKEN_KEY)
    }
  }
}

export const userState = new UserState()
