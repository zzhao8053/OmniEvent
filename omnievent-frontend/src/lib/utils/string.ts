export function isNullOrEmpty(value: string | null | undefined): boolean {
  return value === null || value === undefined || value === ''
}

export function isNullOrWhitespace(value: string | null | undefined): boolean {
  return isNullOrEmpty(value) || value!.trim().length === 0
}

export function trim(value: string): string {
  return value?.trim() || ''
}

export function capitalize(value: string): string {
  if (!value) return ''
  return value.charAt(0).toUpperCase() + value.slice(1).toLowerCase()
}

export function truncate(value: string, maxLength: number, suffix = '...'): string {
  if (!value || value.length <= maxLength) {
    return value
  }
  return value.substring(0, maxLength - suffix.length) + suffix
}

export function randomString(length: number): string {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
  let result = ''
  for (let i = 0; i < length; i++) {
    result += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  return result
}

export function escapeHtml(unsafe: string): string {
  return unsafe
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#039;')
}

export function unescapeHtml(safe: string): string {
  return safe
    .replace(/&amp;/g, '&')
    .replace(/&lt;/g, '<')
    .replace(/&gt;/g, '>')
    .replace(/&quot;/g, '"')
    .replace(/&#039;/g, "'")
}

export function slugify(text: string): string {
  return text
    .toLowerCase()
    .trim()
    .replace(/[^\w\s-]/g, '')
    .replace(/[\s_-]+/g, '-')
    .replace(/^-+|-+$/g, '')
}
