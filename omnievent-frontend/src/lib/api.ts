export interface ApiResponse<T> {
    readonly success: boolean;
    readonly result: T;
}

export class ApiException extends Error {
    constructor(
        public code: number,
        public message: string,
        public data?: any
    ) {
        super(message);
    }
}

export interface ErrorResponse {
    readonly success: boolean;
    readonly errorCode: number;
    readonly errorMessage: string;
    readonly path: string;
}

export function buildErrorResponse(errorCode: number, errorMessage: string): ErrorResponse {
    const errorResponse: ErrorResponse = {
        success: false,
        errorCode: errorCode,
        errorMessage: errorMessage,
        path: ''
    };

    return errorResponse;
}