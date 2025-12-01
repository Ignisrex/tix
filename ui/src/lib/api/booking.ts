/**
 * Booking API endpoints
 */

import { ApiException } from '@/types/api';
import type { ReserveResponse, PurchaseResponse } from '@/types/booking';

const BASE_URL = process.env.NEXT_PUBLIC_CORE_API_URL || 'http://localhost:8080/api/v1';

/**
 * Reserve a ticket
 * Note: The booking service returns response body even on error status codes (409, 410, 404)
 * so we need to parse the response body even when status is not ok
 */
export async function reserveTicket(ticketId: string): Promise<ReserveResponse> {
  const url = `${BASE_URL}/booking/reserve/${ticketId}`;
  
  try {
    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    const data = await response.json().catch(() => null);

    // Booking service returns response body even on error status codes
    // Check if response has success field (booking response format)
    if (data && typeof data.success === 'boolean') {
      return data as ReserveResponse;
    }

    // If not a booking response format, throw error
    if (!response.ok) {
      throw new ApiException(
        data?.error || data?.message || `Request failed with status ${response.status}`,
        response.status
      );
    }

    return data as ReserveResponse;
  } catch (error) {
    if (error instanceof ApiException) {
      throw error;
    }
    throw new ApiException(
      error instanceof Error ? error.message : 'An unexpected error occurred'
    );
  }
}

/**
 * Purchase a ticket
 * Note: The booking service returns response body even on error status codes
 */
export async function purchaseTicket(ticketId: string): Promise<PurchaseResponse> {
  const url = `${BASE_URL}/booking/purchase/${ticketId}`;
  
  try {
    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    const data = await response.json().catch(() => null);

    // Booking service returns response body even on error status codes
    // Check if response has success field (booking response format)
    if (data && typeof data.success === 'boolean') {
      return data as PurchaseResponse;
    }

    // If not a booking response format, throw error
    if (!response.ok) {
      throw new ApiException(
        data?.error || data?.message || `Request failed with status ${response.status}`,
        response.status
      );
    }

    return data as PurchaseResponse;
  } catch (error) {
    if (error instanceof ApiException) {
      throw error;
    }
    throw new ApiException(
      error instanceof Error ? error.message : 'An unexpected error occurred'
    );
  }
}

