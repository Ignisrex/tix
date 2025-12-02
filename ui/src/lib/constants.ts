/**
 * Shared constants across the application
 */

export const RESERVATION_STORAGE_KEY = "tix_reservation";

export const RESERVATION_TTL_SECONDS = parseInt(
  process.env.NEXT_PUBLIC_RESERVATION_TTL_SECONDS || "180",
  10
);

export const URGENT_RESERVATION_THRESHOLD_SECONDS = 30;

