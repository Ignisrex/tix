import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

/**
 * Format price in cents to dollar string
 * @param cents - Price in cents
 * @returns Formatted price string (e.g., "$10.00")
 */
export function formatPrice(cents: number): string {
  return `$${(cents / 100).toFixed(2)}`;
}

/**
 * Format date string to short format (e.g., "Jan 15, 2024")
 * @param dateString - ISO date string
 * @returns Formatted date string
 */
export function formatDateShort(dateString: string): string {
  return new Date(dateString).toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
    year: "numeric",
  });
}

/**
 * Format date string to long format with time (e.g., "Monday, January 15, 2024 at 7:00 PM")
 * @param dateString - ISO date string
 * @returns Formatted date string
 */
export function formatDateLong(dateString: string): string {
  return new Date(dateString).toLocaleDateString("en-US", {
    weekday: "long",
    year: "numeric",
    month: "long",
    day: "numeric",
    hour: "numeric",
    minute: "2-digit",
  });
}
