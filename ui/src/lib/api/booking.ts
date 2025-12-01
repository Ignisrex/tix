/**
 * Booking API endpoints
 */

import { ApiException } from '@/types/api';
import type { ReserveResponse, PurchaseResponse, PurchaseDetailsResponse } from '@/types/booking';

const BASE_URL = process.env.NEXT_PUBLIC_CORE_API_URL || 'http://localhost:8080/api/v1';

import type { ReserveRequest, PurchaseRequest } from '@/types/booking';

/**
 * Reserve tickets (supports single or multiple)
 * Note: The booking service returns response body even on error status codes (409, 410, 404)
 * so we need to parse the response body even when status is not ok
 */
export async function reserveTickets(ticketIds: string[]): Promise<ReserveResponse> {
  const url = `${BASE_URL}/booking/reserve`;
  
  try {
    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ ticket_ids: ticketIds } as ReserveRequest),
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
 * Reserve a single ticket (backward compatibility)
 */
export async function reserveTicket(ticketId: string): Promise<ReserveResponse> {
  return reserveTickets([ticketId]);
}

/**
 * Purchase tickets (supports single or multiple)
 * Note: The booking service returns response body even on error status codes
 */
export async function purchaseTickets(ticketIds: string[]): Promise<PurchaseResponse> {
  const url = `${BASE_URL}/booking/purchase`;
  
  try {
    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ ticket_ids: ticketIds } as PurchaseRequest),
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

/**
 * Purchase a single ticket (backward compatibility)
 */
export async function purchaseTicket(ticketId: string): Promise<PurchaseResponse> {
  return purchaseTickets([ticketId]);
}

/**
 * Get purchase details by purchase ID
 */
export async function getPurchaseDetails(purchaseId: string): Promise<PurchaseDetailsResponse> {
  const url = `${BASE_URL}/booking/purchases/${purchaseId}`;
  
  try {
    const response = await fetch(url, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      const data = await response.json().catch(() => null);
      throw new ApiException(
        data?.error || data?.message || `Request failed with status ${response.status}`,
        response.status
      );
    }

    return await response.json() as PurchaseDetailsResponse;
  } catch (error) {
    if (error instanceof ApiException) {
      throw error;
    }
    throw new ApiException(
      error instanceof Error ? error.message : 'An unexpected error occurred'
    );
  }
}

