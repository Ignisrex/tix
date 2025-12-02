/**
 * Custom hook for managing reservations
 */

import { useState, useEffect, useCallback } from "react";
import { reserveTickets } from "@/lib/api";
import type { ReservationData } from "@/types/booking";
import {
  getValidReservation,
  calculateRemainingSeconds,
  saveReservation,
  clearReservation,
} from "@/lib/reservation-storage";

interface UseReservationReturn {
  reservation: ReservationData | null;
  remainingSeconds: number;
  ticketIds: string[];
  isExpired: boolean;
  clear: () => void;
  reserveTicketsForEvent: (
    ticketIds: string[],
    eventId: string,
    options?: {
      mergeWithExisting?: boolean;
      onSuccess?: () => void;
      onError?: (error: string) => void;
    }
  ) => Promise<void>;
}

/**
 * Hook for managing reservation state and operations
 */
export function useReservation(): UseReservationReturn {
  const [reservation, setReservation] = useState<ReservationData | null>(null);
  const [remainingSeconds, setRemainingSeconds] = useState<number>(0);

  // Load reservation on mount and update periodically
  useEffect(() => {
    const updateReservation = () => {
      const res = getValidReservation();
      setReservation(res);
      // Calculate remaining seconds from the already-fetched reservation
      setRemainingSeconds(res ? calculateRemainingSeconds(res) : 0);
    };

    // Initial load
    updateReservation();

    // Update every second
    const interval = setInterval(updateReservation, 1000);

    return () => clearInterval(interval);
  }, []);

  const clear = useCallback(() => {
    clearReservation();
    setReservation(null);
    setRemainingSeconds(0);
  }, []);

  const reserveTicketsForEvent = useCallback(
    async (
      ticketIds: string[],
      eventId: string,
      options?: {
        mergeWithExisting?: boolean;
        onSuccess?: () => void;
        onError?: (error: string) => void;
      }
    ) => {
      const { mergeWithExisting = true, onSuccess, onError } = options || {};

      try {
        // Check for existing reservation (use current state if available, otherwise fetch)
        const existing = reservation || getValidReservation();
        let ticketsToReserve = ticketIds;

        if (existing) {
          if (existing.eventId !== eventId) {
            const errorMsg =
              "You already have tickets reserved for a different event. Please complete or cancel that reservation first.";
            onError?.(errorMsg);
            return;
          }

          if (mergeWithExisting) {
            // Filter out tickets that are already reserved
            const newTicketIds = ticketIds.filter(
              (id) => !existing.ticketIds.includes(id)
            );

            // If all tickets are already reserved, just update localStorage
            if (newTicketIds.length === 0) {
              // Merge and deduplicate (no need to call mergeReservationTickets since we already have existing)
              const mergedTicketIds = [...new Set([...existing.ticketIds, ...ticketIds])];
              const updatedReservation: ReservationData = {
                ticketIds: mergedTicketIds,
                eventId: eventId,
                reservedAt: Date.now(),
              };
              saveReservation(updatedReservation);
              setReservation(updatedReservation);
              onSuccess?.();
              return;
            }

            ticketsToReserve = newTicketIds;
          }
        }

        // Reserve the tickets
        const response = await reserveTickets(ticketsToReserve);

        if (response.success) {
          // Merge with existing if applicable
          const allReservedTicketIds = mergeWithExisting && existing
            ? [...new Set([...existing.ticketIds, ...response.ticket_ids])]
            : response.ticket_ids;

          const reservationData: ReservationData = {
            ticketIds: allReservedTicketIds,
            eventId: eventId,
            reservedAt: Date.now(),
          };

          saveReservation(reservationData);
          setReservation(reservationData);
          onSuccess?.();
        } else {
          const errorMsg = response.message || "One or more seats are not available";
          onError?.(errorMsg);
        }
      } catch (err) {
        const errorMessage =
          err instanceof Error ? err.message : "Failed to reserve tickets";
        onError?.(errorMessage);
      }
    },
    [reservation]
  );

  return {
    reservation,
    remainingSeconds,
    ticketIds: reservation?.ticketIds || [],
    isExpired: remainingSeconds <= 0 && reservation !== null,
    clear,
    reserveTicketsForEvent,
  };
}

