"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { getEvent, getEventTickets, purchaseTickets } from "@/lib/api";
import type { Event, Ticket } from "@/types/events";
import type { ReservationData } from "@/types/booking";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { ArrowLeft } from "lucide-react";
import Link from "next/link";
import { formatPrice } from "@/lib/utils";

const RESERVATION_TTL_SECONDS = parseInt(
  process.env.NEXT_PUBLIC_RESERVATION_TTL_SECONDS || "180",
  10
);

export default function CheckoutPage() {
  const router = useRouter();

  const [event, setEvent] = useState<Event | null>(null);
  const [tickets, setTickets] = useState<Ticket[]>([]);
  const [loading, setLoading] = useState(true);
  const [purchasing, setPurchasing] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [reservationValid, setReservationValid] = useState<boolean | null>(null);

  useEffect(() => {
    async function validateAndFetch() {
      const reservationStr = localStorage.getItem("tix_reservation");
      if (!reservationStr) {
        setReservationValid(false);
        setLoading(false);
        return;
      }

      try {
        const reservation: ReservationData = JSON.parse(reservationStr);
        
        if (!reservation.ticketIds || reservation.ticketIds.length === 0) {
          setReservationValid(false);
          setLoading(false);
          return;
        }

        const now = Date.now();
        const elapsed = Math.floor((now - reservation.reservedAt) / 1000);
        const remaining = RESERVATION_TTL_SECONDS - elapsed;

        if (remaining <= 0) {
          localStorage.removeItem("tix_reservation");
          setReservationValid(false);
          setLoading(false);
          return;
        }

        setReservationValid(true);

        // Fetch event and tickets
        try {
          setLoading(true);
          const eventId = reservation.eventId;

          const [eventData, ticketsData] = await Promise.all([
            getEvent(eventId),
            getEventTickets(eventId),
          ]);

          setEvent(eventData);
          
          // Find all reserved tickets
          // TODO: Store ticket data in a store and pass it to the checkout page instead of fetching it
          //    -> introducing zustand for this would be a good idea
          const foundTickets = ticketsData.filter((t) => 
            reservation.ticketIds.includes(t.id)
          );
          
          if (foundTickets.length !== reservation.ticketIds.length) {
            setError("Some tickets were not found");
            setLoading(false);
            return;
          }

          setTickets(foundTickets);
        } catch (err) {
          setError(err instanceof Error ? err.message : "Failed to load ticket details");
        } finally {
          setLoading(false);
        }
      } catch (err) {
        console.error("Error validating reservation:", err);
        localStorage.removeItem("tix_reservation");
        setReservationValid(false);
        setLoading(false);
      }
    }

    validateAndFetch();
  }, []);

  const handlePurchase = async () => {
    if (tickets.length === 0) return;

    const reservationStr = localStorage.getItem("tix_reservation");
    if (!reservationStr) {
      alert("Reservation not found");
      return;
    }

    try {
      const reservation: ReservationData = JSON.parse(reservationStr);
      const ticketIds = reservation.ticketIds;

      setPurchasing(true);
      const response = await purchaseTickets(ticketIds);
      
      if (response.success && response.purchase_id) {
        // Clear reservation
        localStorage.removeItem("tix_reservation");
        
        // Redirect to success page with purchase ID
        router.push(`/purchases/${response.purchase_id}`);
      } else {
        alert(response.message || "Purchase failed");
        // Redirect back to event page
        if (event) {
          router.push(`/events/${event.id}`);
        }
      }
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to purchase tickets";
      alert(errorMessage);
      // Redirect back to event page
      if (event) {
        router.push(`/events/${event.id}`);
      }
    } finally {
      setPurchasing(false);
    }
  };

  const totalCents = tickets.reduce((sum, ticket) => sum + (ticket.ticket_type_price_cents || 0), 0);

  if (loading) {
    return (
      <div className="min-h-screen bg-background">
        <div className="container mx-auto px-4 py-8 max-w-2xl">
          <Skeleton className="h-10 w-32 mb-6" />
          <Skeleton className="h-64 w-full" />
        </div>
      </div>
    );
  }

  if (!reservationValid) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <Card className="max-w-md">
          <CardHeader>
            <CardTitle>Reservation Expired</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-muted-foreground mb-4">
              Your reservation has expired or is invalid. Please select seats again.
            </p>
            <Link href="/">
              <Button variant="outline">Go Home</Button>
            </Link>
          </CardContent>
        </Card>
      </div>
    );
  }

  if (error || !tickets.length || !event) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <Card className="max-w-md">
          <CardHeader>
            <CardTitle>Error</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-muted-foreground mb-4">
              {error || "Failed to load ticket details"}
            </p>
            <Link href="/">
              <Button variant="outline">Go Home</Button>
            </Link>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background">
      <div className="container mx-auto px-4 py-8 max-w-2xl">
        {event && (
          <Link href={`/events/${event.id}`}>
            <Button variant="ghost" className="mb-6">
              <ArrowLeft className="mr-2 h-4 w-4" />
              Back to Event
            </Button>
          </Link>
        )}

        <Card>
          <CardHeader>
            <CardTitle>Checkout</CardTitle>
          </CardHeader>
          <CardContent className="space-y-6">
            <div>
              <h2 className="text-2xl font-bold mb-2">{event.title}</h2>
              <p className="text-muted-foreground">{event.description}</p>
            </div>

            <div className="border-t pt-4">
              <h3 className="font-semibold mb-4">Ticket Details ({tickets.length} {tickets.length === 1 ? 'ticket' : 'tickets'})</h3>
              <div className="space-y-3">
                {tickets.map((ticket) => (
                  <div key={ticket.id} className="flex justify-between items-center p-3 rounded-lg bg-muted/50">
                    <div>
                      <p className="font-medium">{ticket.ticket_type_display_name || ticket.ticket_type_name}</p>
                      </div>
                    <span className="font-semibold">{formatPrice(ticket.ticket_type_price_cents || 0)}</span>
                  </div>
                ))}
              </div>
            </div>

            <div className="border-t pt-4">
              <div className="flex justify-between items-center mb-4">
                <span className="text-lg font-semibold">Total:</span>
                <span className="text-2xl font-bold">{formatPrice(totalCents)}</span>
              </div>
              <Button
                onClick={handlePurchase}
                disabled={purchasing}
                className="w-full bg-indigo-500 hover:bg-indigo-600 text-white"
                size="lg"
              >
                {purchasing ? "Processing..." : "Complete Purchase"}
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

