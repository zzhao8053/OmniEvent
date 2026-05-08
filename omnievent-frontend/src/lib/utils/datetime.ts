export function formatUnixTime(unixTime: number, format = '%Y-%m-%d %H:%M:%S'): string {
  if (!unixTime || unixTime <= 0) {
    return ''
  }

  const date = new Date(unixTime * 1000)
  return formatDate(date, format)
}

export function formatDate(date: Date, format = '%Y-%m-%d %H:%M:%S'): string {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')

  return format
    .replace('%Y', String(year))
    .replace('%m', month)
    .replace('%d', day)
    .replace('%H', hours)
    .replace('%M', minutes)
    .replace('%S', seconds)
}

export function getUnixTime(date: Date = new Date()): number {
  return Math.floor(date.getTime() / 1000)
}

export function fromUnixTime(unixTime: number): Date {
  return new Date(unixTime * 1000)
}

export function now(): number {
  return getUnixTime(new Date())
}

export function addDays(date: Date, days: number): Date {
  const result = new Date(date)
  result.setDate(result.getDate() + days)
  return result
}

export function addHours(date: Date, hours: number): Date {
  const result = new Date(date)
  result.setHours(result.getHours() + hours)
  return result
}

export function startOfDay(date: Date): Date {
  const result = new Date(date)
  result.setHours(0, 0, 0, 0)
  return result
}

export function endOfDay(date: Date): Date {
  const result = new Date(date)
  result.setHours(23, 59, 59, 999)
  return result
}

export function isToday(date: Date): boolean {
  const today = new Date()
  return (
    date.getDate() === today.getDate() &&
    date.getMonth() === today.getMonth() &&
    date.getFullYear() === today.getFullYear()
  )
}

export function isSameDay(date1: Date, date2: Date): boolean {
  return (
    date1.getDate() === date2.getDate() &&
    date1.getMonth() === date2.getMonth() &&
    date1.getFullYear() === date2.getFullYear()
  )
}

export function getDaysBetween(date1: Date, date2: Date): number {
  const oneDay = 24 * 60 * 60 * 1000
  return Math.round(Math.abs((date1.getTime() - date2.getTime()) / oneDay))
}
