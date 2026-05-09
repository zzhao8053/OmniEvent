export const BASE_API_URL_PATH: string = '/api';
export const BASE_QRCODE_PATH: string = '/qrcode';
export const BASE_PROXY_URL_PATH: string = '/proxy';

export const DEFAULT_API_TIMEOUT: number = 10000; // 10s
export const DEFAULT_UPLOAD_API_TIMEOUT: number = 30000; // 30s
export const DEFAULT_EXPORT_API_TIMEOUT: number = 180000; // 180s
export const DEFAULT_IMPORT_API_TIMEOUT: number = 1800000; // 1800s

export enum KnownErrorCode {
    ApiNotFound = 100001,
    ValidatorError = 4000,
    UserEmailNotVerified = 2001,
    UserNotFound = 2002,
    UserAlreadyExists = 2003,
    InvalidCredentials = 2004,
    TokenExpired = 3001,
    TokenInvalid = 3002,
    TokenMissing = 3003
}

export interface SpecifiedApiError {
    readonly message: string;
}

export const SPECIFIED_API_NOT_FOUND_ERRORS: Record<string, SpecifiedApiError> = {
    '/api/register.json': {
        message: 'User registration is disabled'
    },
    '/api/authorize.json': {
        message: 'Username/password login is disabled'
    },
    '/api/token/refresh.json': {
        message: 'Token refresh is disabled'
    },
    '/api/user/profile.json': {
        message: 'User profile access is disabled'
    },
    '/api/user/password/reset.json': {
        message: 'Password reset is disabled'
    }
};

export interface ParameterizedError {
    readonly localeKey: string;
    readonly regex: RegExp;
    readonly parameters: {
        readonly field: string;
        readonly localized: boolean;
    }[];
}

export const PARAMETERIZED_ERRORS: ParameterizedError[] = [
    {
        localeKey: 'parameter invalid',
        regex: /^parameter "(\w+)" is invalid$/,
        parameters: [{
            field: 'parameter',
            localized: true
        }]
    },
    {
        localeKey: 'parameter required',
        regex: /^parameter "(\w+)" is required$/,
        parameters: [{
            field: 'parameter',
            localized: true
        }]
    },
    {
        localeKey: 'parameter too large',
        regex: /^parameter "(\w+)" must be less than (\d+)$/,
        parameters: [{
            field: 'parameter',
            localized: true
        }, {
            field: 'number',
            localized: false
        }]
    },
    {
        localeKey: 'parameter too long',
        regex: /^parameter "(\w+)" must be less than (\d+) characters$/,
        parameters: [{
            field: 'parameter',
            localized: true
        }, {
            field: 'length',
            localized: false
        }]
    },
    {
        localeKey: 'parameter too small',
        regex: /^parameter "(\w+)" must be more than (\d+)$/,
        parameters: [{
            field: 'parameter',
            localized: true
        }, {
            field: 'number',
            localized: false
        }]
    },
    {
        localeKey: 'parameter too short',
        regex: /^parameter "(\w+)" must be more than (\d+) characters$/,
        parameters: [{
            field: 'parameter',
            localized: true
        }, {
            field: 'length',
            localized: false
        }]
    },
    {
        localeKey: 'parameter cannot be blank',
        regex: /^parameter "(\w+)" cannot be blank$/,
        parameters: [{
            field: 'parameter',
            localized: true
        }]
    },
    {
        localeKey: 'parameter invalid username format',
        regex: /^parameter "(\w+)" is invalid username format$/,
        parameters: [{
            field: 'parameter',
            localized: true
        }]
    },
    {
        localeKey: 'parameter invalid email format',
        regex: /^parameter "(\w+)" is invalid email format$/,
        parameters: [{
            field: 'parameter',
            localized: true
        }]
    }
];