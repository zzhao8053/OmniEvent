const LOG_PREFIX = '[OmniEvent]'

export enum LogLevel {
  DEBUG = 0,
  INFO = 1,
  WARN = 2,
  ERROR = 3
}

let currentLogLevel = LogLevel.INFO

export function setLogLevel(level: LogLevel): void {
  currentLogLevel = level
}

export function isDebugEnabled(): boolean {
  return currentLogLevel <= LogLevel.DEBUG
}

function formatMessage(level: string, message: string, ...args: any[]): string {
  const timestamp = new Date().toISOString()
  const formattedArgs = args.length > 0 ? ' ' + args.map(arg =>
    typeof arg === 'object' ? JSON.stringify(arg) : String(arg)
  ).join(' ') : ''
  return `${timestamp} ${LOG_PREFIX} [${level}] ${message}${formattedArgs}`
}

export function debug(message: string, ...args: any[]): void {
  if (currentLogLevel <= LogLevel.DEBUG) {
    console.debug(formatMessage('DEBUG', message, ...args))
  }
}

export function info(message: string, ...args: any[]): void {
  if (currentLogLevel <= LogLevel.INFO) {
    console.info(formatMessage('INFO', message, ...args))
  }
}

export function warn(message: string, ...args: any[]): void {
  if (currentLogLevel <= LogLevel.WARN) {
    console.warn(formatMessage('WARN', message, ...args))
  }
}

export function error(message: string, ...args: any[]): void {
  if (currentLogLevel <= LogLevel.ERROR) {
    console.error(formatMessage('ERROR', message, ...args))
  }
}

export const logger = {
  debug,
  info,
  warn,
  error,
  setLogLevel,
  isDebugEnabled
}

export default logger
