/**
 * Booking API endpoints
 */

import { request } from './client';
import type { ReserveResponse, PurchaseResponse } from '@/types/booking';

/**
 * Reserve a ticket
 */
export async function reserveTicket(ticketId: string): Promise<ReserveResponse> {
  return request<ReserveResponse>(`/booking/reserve/${ticketId}`, {
    method: 'POST',
  });
}

/**
 * Purchase a ticket
 */
export async function purchaseTicket(ticketId: string): Promise<PurchaseResponse> {
  return request<PurchaseResponse>(`/booking/purchase/${ticketId}`, {
    method: 'POST',
  });
}

