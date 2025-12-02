/**
 * Reservation storage utilities for managing reservation data in localStorage
 */

import type { ReservationData } from "@/types/booking";
import { RESERVATION_STORAGE_KEY, RESERVATION_TTL_SECONDS } from "./constants";

/**
 * Get the current reservation from localStorage
 * @returns ReservationData if valid, null otherwise
 */
export function getReservation(): ReservationData | null {
  try {
    const reservationStr = localStorage.getItem(RESERVATION_STORAGE_KEY);
    if (!reservationStr) {
      return null;
    }

    const reservation: ReservationData = JSON.parse(reservationStr);
    
    // Validate structure
    if (!reservation.ticketIds || !Array.isArray(reservation.ticketIds) || reservation.ticketIds.length === 0) {
      return null;
    }

    if (!reservation.eventId || !reservation.reservedAt) {
      return null;
    }

    return reservation;
  } catch (err) {
    console.error("Error parsing reservation data:", err);
    // Clear invalid data
    localStorage.removeItem(RESERVATION_STORAGE_KEY);
    return null;
  }
}

/**
 * Save reservation to localStorage
 */
export function saveReservation(reservation: ReservationData): void {
  try {
    localStorage.setItem(RESERVATION_STORAGE_KEY, JSON.stringify(reservation));
  } catch (err) {
    console.error("Error saving reservation:", err);
    throw new Error("Failed to save reservation");
  }
}

/**
 * Remove reservation from localStorage
 */
export function clearReservation(): void {
  localStorage.removeItem(RESERVATION_STORAGE_KEY);
}

/**
 * Check if reservation exists and is valid (not expired)
 * @returns ReservationData if valid, null if expired or missing
 */
export function getValidReservation(): ReservationData | null {
  const reservation = getReservation();
  if (!reservation) {
    return null;
  }

  const now = Date.now();
  const elapsed = Math.floor((now - reservation.reservedAt) / 1000);
  const remaining = RESERVATION_TTL_SECONDS - elapsed;

  if (remaining <= 0) {
    // Expired - clear it
    clearReservation();
    return null;
  }

  return reservation;
}

/**
 * Calculate remaining seconds for a reservation (does not validate or clear)
 * @param reservation The reservation to calculate remaining time for
 * @returns number of seconds remaining, or 0 if expired
 */
export function calculateRemainingSeconds(reservation: ReservationData): number {
  const now = Date.now();
  const elapsed = Math.floor((now - reservation.reservedAt) / 1000);
  const remaining = RESERVATION_TTL_SECONDS - elapsed;
  return Math.max(remaining, 0);
}

/**
 * Get remaining seconds for a reservation
 * @returns number of seconds remaining, or 0 if expired/missing
 */
export function getReservationRemainingSeconds(): number {
  const reservation = getValidReservation();
  if (!reservation) {
    return 0;
  }
  return calculateRemainingSeconds(reservation);
}

/**
 * Merge new ticket IDs with existing reservation
 * @param newTicketIds Array of new ticket IDs to add
 * @param eventId Event ID (must match existing reservation if one exists)
 * @returns Merged ticket IDs array
 */
export function mergeReservationTickets(
  newTicketIds: string[],
  eventId: string
): string[] {
  const existing = getValidReservation();
  
  if (!existing) {
    return newTicketIds;
  }

  // If different event, don't merge
  if (existing.eventId !== eventId) {
    return newTicketIds;
  }

  // Merge and deduplicate
  const merged = [...new Set([...existing.ticketIds, ...newTicketIds])];
  return merged;
}

