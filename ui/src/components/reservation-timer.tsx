"use client";

import { useEffect, useState } from "react";
import { useRouter, usePathname } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Clock } from "lucide-react";
import type { ReservationData } from "@/types/booking";

const RESERVATION_TTL_SECONDS = parseInt(
  process.env.NEXT_PUBLIC_RESERVATION_TTL_SECONDS || "180",
  10
);

export function ReservationTimer() {
  const router = useRouter();
  const pathname = usePathname();
  const [remainingSeconds, setRemainingSeconds] = useState<number | null>(null);
  const [ticketIds, setTicketIds] = useState<string[]>([]);
  
  const isOnCheckoutPage = pathname?.startsWith("/checkout");

  useEffect(() => {
    // Check localStorage for active reservation
    const checkReservation = () => {
      const reservationStr = localStorage.getItem("tix_reservation");
      if (!reservationStr) {
        setRemainingSeconds(null);
        setTicketIds([]);
        return;
      }

      try {
        const reservation: ReservationData = JSON.parse(reservationStr);
        const now = Date.now();
        const elapsed = Math.floor((now - reservation.reservedAt) / 1000);
        const remaining = RESERVATION_TTL_SECONDS - elapsed;

        if (remaining <= 0) {
          // Reservation expired
          localStorage.removeItem("tix_reservation");
          setRemainingSeconds(null);
          setTicketIds([]);
          return;
        }

        setRemainingSeconds(remaining);
        setTicketIds(reservation.ticketIds || []);
      } catch (err) {
        console.error("Error parsing reservation data:", err);
        localStorage.removeItem("tix_reservation");
        setRemainingSeconds(null);
        setTicketIds([]);
      }
    };

    checkReservation();

    // Update every second 
    // - ideally we would have webhook to the server that would hit us with the reservation time
    const interval = setInterval(() => {
      checkReservation();
    }, 1000);

    return () => clearInterval(interval);
  }, []);

  const handleGoToCheckout = () => {
    if (ticketIds.length > 0 && !isOnCheckoutPage) {
      router.push("/checkout");
    }
  };

  // Don't render if no active reservation
  if (remainingSeconds === null || ticketIds.length === 0) {
    return null;
  }

  // Format time as M:SS
  const minutes = Math.floor(remainingSeconds / 60);
  const seconds = remainingSeconds % 60;
  const timeString = `${minutes}:${seconds.toString().padStart(2, "0")}`;

  
  const isUrgent = remainingSeconds < 30;

  return (
    <div className="fixed bottom-6 left-6 z-50 animate-in fade-in slide-in-from-bottom-4 duration-300">
      <div className="rounded-2xl bg-background/95 backdrop-blur-xl border border-border/50 shadow-xl shadow-black/5 p-5 min-w-[280px]">
        <div className="flex flex-col gap-4">
          {/* Header with icon */}
          <div className="flex items-center gap-3">
            <div className={`rounded-full p-2 ${isUrgent ? 'bg-red-100 dark:bg-red-900/20' : 'bg-indigo-100 dark:bg-indigo-900/20'}`}>
              <Clock className={`h-4 w-4 ${isUrgent ? 'text-red-600 dark:text-red-400' : 'text-indigo-600 dark:text-indigo-400'}`} />
            </div>
            <div className="flex-1">
              <p className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
                Reservation Active
              </p>
              <p className="text-sm text-muted-foreground mt-0.5">
                {ticketIds.length} {ticketIds.length === 1 ? 'ticket' : 'tickets'} • Expires in <span className={`font-semibold ${isUrgent ? 'text-red-600 dark:text-red-400' : 'text-foreground'}`}>{timeString}</span>
              </p>
            </div>
          </div>

          {/* Button */}
          <Button 
            onClick={handleGoToCheckout} 
            disabled={isOnCheckoutPage}
            className={`w-full rounded-lg font-medium transition-all duration-200 ${
              isOnCheckoutPage 
                ? "bg-muted text-muted-foreground cursor-not-allowed opacity-60" 
                : "bg-indigo-500 hover:bg-indigo-600 text-white shadow-md hover:shadow-lg"
            }`}
          >
            Checkout
          </Button>
        </div>
      </div>
    </div>
  );
}

