export interface ValidationRule {
  validate: (value: any) => boolean
  message: string
}

export function required(message = 'This field is required'): ValidationRule {
  return {
    validate: (value: any) => {
      if (typeof value === 'string') {
        return value.trim().length > 0
      }
      return value !== null && value !== undefined
    },
    message
  }
}

export function minLength(min: number, message?: string): ValidationRule {
  return {
    validate: (value: any) => {
      if (typeof value === 'string') {
        return value.length >= min
      }
      return false
    },
    message: message || `Minimum length is ${min} characters`
  }
}

export function maxLength(max: number, message?: string): ValidationRule {
  return {
    validate: (value: any) => {
      if (typeof value === 'string') {
        return value.length <= max
      }
      return false
    },
    message: message || `Maximum length is ${max} characters`
  }
}

export function email(message = 'Invalid email address'): ValidationRule {
  return {
    validate: (value: any) => {
      if (typeof value !== 'string') return false
      return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)
    },
    message
  }
}

export function pattern(regex: RegExp, message = 'Invalid format'): ValidationRule {
  return {
    validate: (value: any) => {
      if (typeof value !== 'string') return false
      return regex.test(value)
    },
    message
  }
}

export function minValue(min: number, message?: string): ValidationRule {
  return {
    validate: (value: any) => {
      if (typeof value !== 'number') return false
      return value >= min
    },
    message: message || `Minimum value is ${min}`
  }
}

export function maxValue(max: number, message?: string): ValidationRule {
  return {
    validate: (value: any) => {
      if (typeof value !== 'number') return false
      return value <= max
    },
    message: message || `Maximum value is ${max}`
  }
}

export function username(message = 'Username must be 3-32 characters, alphanumeric or underscore'): ValidationRule {
  return {
    validate: (value: any) => {
      if (typeof value !== 'string') return false
      return /^[a-zA-Z0-9_]{3,32}$/.test(value)
    },
    message
  }
}

export function password(message = 'Password must be at least 6 characters'): ValidationRule {
  return {
    validate: (value: any) => {
      if (typeof value !== 'string') return false
      return value.length >= 6
    },
    message
  }
}

export type Rules = Record<string, ValidationRule[]>

export function validate(rules: Rules, data: Record<string, any>): Record<string, string> {
  const errors: Record<string, string> = {}

  for (const field in rules) {
    const fieldRules = rules[field]
    const value = data[field]

    for (const rule of fieldRules) {
      if (!rule.validate(value)) {
        errors[field] = rule.message
        break
      }
    }
  }

  return errors
}

export function hasErrors(errors: Record<string, string>): boolean {
  return Object.keys(errors).length > 0
}
