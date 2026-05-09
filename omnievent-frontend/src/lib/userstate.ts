import CryptoJS from 'crypto-js';

import { isString, isObject } from './common.ts';

const appLockSecretBaseStringPrefix: string = 'OMNIEVENT_LOCK_SECRET_';

const tokenLocalStorageKey: string = 'omnievent_user_token';
const userInfoLocalStorageKey: string = 'omnievent_user_info';

const tokenSessionStorageKey: string = 'omnievent_user_session_token';
const encryptedTokenSessionStorageKey: string = 'omnievent_user_session_encrypted_token';
const appLockStateSessionStorageKey: string = 'omnievent_user_app_lock_state'; // { 'username': '', secret: '' }

function getAppLockSecret(pinCode: string): string {
    const hashedPinCode = CryptoJS.SHA256(appLockSecretBaseStringPrefix + pinCode).toString();
    return hashedPinCode.substring(0, 24);
}

function getEncryptedToken(token: string, appLockState: { username: string; secret: string }): string {
    const key = CryptoJS.SHA256(`${appLockSecretBaseStringPrefix}|${appLockState.username}|${appLockState.secret}`).toString();
    return CryptoJS.AES.encrypt(token, key).toString();
}

function getDecryptedToken(encryptedToken: string, appLockState: { username: string; secret: string }): string {
    const key = CryptoJS.SHA256(`${appLockSecretBaseStringPrefix}|${appLockState.username}|${appLockState.secret}`).toString();
    const bytes = CryptoJS.AES.decrypt(encryptedToken, key);
    return bytes.toString(CryptoJS.enc.Utf8);
}

export function isUserLogined(): boolean {
    return !!localStorage.getItem(tokenLocalStorageKey);
}

export function getCurrentToken(): string | null {
    return localStorage.getItem(tokenLocalStorageKey);
}

export function updateCurrentToken(token: string): void {
    if (isString(token)) {
        localStorage.setItem(tokenLocalStorageKey, token);
    }
}

export function clearCurrentToken(): void {
    localStorage.removeItem(tokenLocalStorageKey);
    sessionStorage.removeItem(tokenSessionStorageKey);
    sessionStorage.removeItem(encryptedTokenSessionStorageKey);
    sessionStorage.removeItem(appLockStateSessionStorageKey);
}

export function getCurrentUserInfo(): Record<string, unknown> | null {
    const data = localStorage.getItem(userInfoLocalStorageKey);

    if (!data) {
        return null;
    }

    return JSON.parse(data) as Record<string, unknown>;
}

export function updateCurrentUserInfo(user: Record<string, unknown>): void {
    if (isObject(user)) {
        localStorage.setItem(userInfoLocalStorageKey, JSON.stringify(user));
    }
}

export function clearCurrentUserInfo(): void {
    localStorage.removeItem(userInfoLocalStorageKey);
}


// ========== OmniEvent UserState class (保持原有接口兼容) ==========

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
