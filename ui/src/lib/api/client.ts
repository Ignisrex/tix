/**
 * Base HTTP client for API requests
 */

import { ApiException } from '@/types/api';

const BASE_URL = process.env.NEXT_PUBLIC_CORE_API_URL || 'http://localhost:8080/api/v1';

interface RequestConfig extends RequestInit {
  params?: Record<string, string | number>;
}

/**
 * Generic request function with error handling
 */
export async function request<T>(path: string, config?: RequestConfig): Promise<T> {
  const { params, ...fetchOptions } = config || {};

  // Build URL with query params
  let url = `${BASE_URL}${path}`;
  if (params && Object.keys(params).length > 0) {
    const searchParams = new URLSearchParams();
    Object.entries(params).forEach(([key, value]) => {
      searchParams.append(key, String(value));
    });
    url += `?${searchParams.toString()}`;
  }

  try {
    const response = await fetch(url, {
      ...fetchOptions,
      headers: {
        'Content-Type': 'application/json',
        ...(fetchOptions.headers || {}),
      },
    });

    const data = await response.json().catch(() => null);

    if (!response.ok) {
      throw new ApiException(
        data?.error || `Request failed with status ${response.status}`,
        response.status
      );
    }

    return data as T;
  } catch (error) {
    if (error instanceof ApiException) {
      throw error;
    }
    throw new ApiException(
      error instanceof Error ? error.message : 'An unexpected error occurred'
    );
  }
}

