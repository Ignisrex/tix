/**
 * Shared API types and utilities
 */

export interface ApiError {
  error: string;
  status?: number;
}

export interface RequestOptions extends RequestInit {
  params?: Record<string, string | number>;
}

export class ApiException extends Error {
  status?: number;

  constructor(message: string, status?: number) {
    super(message);
    this.name = 'ApiException';
    this.status = status;
  }
}

