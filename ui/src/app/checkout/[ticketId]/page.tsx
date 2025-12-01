"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { getEvent, getEventTickets, purchaseTicket } from "@/lib/api";
import type { Event, Ticket } from "@/types/events";
import type { ReservationData } from "@/types/booking";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { ArrowLeft } from "lucide-react";
import Link from "next/link";
import { formatPrice } from "@/lib/utils";

const RESERVATION_TTL_SECONDS = 180; // 3 minutes

export default function CheckoutPage() {
  const params = useParams();
  const router = useRouter();
  const ticketId = params.ticketId as string;

  const [event, setEvent] = useState<Event | null>(null);
  const [ticket, setTicket] = useState<Ticket | null>(null);
  const [loading, setLoading] = useState(true);
  const [purchasing, setPurchasing] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [reservationValid, setReservationValid] = useState<boolean | null>(null);

  useEffect(() => {
    // Validate reservation and fetch data
    async function validateAndFetch() {
      const reservationStr = localStorage.getItem("tix_reservation");
      if (!reservationStr) {
        setReservationValid(false);
        setLoading(false);
        return;
      }

      try {
        const reservation: ReservationData = JSON.parse(reservationStr);
        if (reservation.ticketId !== ticketId) {
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
          
          // Find the specific ticket
          const foundTicket = ticketsData.find((t) => t.id === ticketId);
          if (!foundTicket) {
            setError("Ticket not found");
            setLoading(false);
            return;
          }

          setTicket(foundTicket);
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
  }, [ticketId]);

  const handlePurchase = async () => {
    if (!ticket) return;

    try {
      setPurchasing(true);
      const response = await purchaseTicket(ticketId);
      
      if (response.success) {
        // Clear reservation
        localStorage.removeItem("tix_reservation");
        
        if (event) {
          router.push(`/success/${event.id}/${ticketId}`);
        } else {
          alert("Purchase successful but unable to redirect. Please go to home page.");
        }
      } else {
        alert(response.message || "Purchase failed");
      }
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to purchase ticket";
      alert(errorMessage);
    } finally {
      setPurchasing(false);
    }
  };


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
              Your reservation has expired or is invalid. Please select a seat again.
            </p>
            <Link href="/">
              <Button variant="outline">Go Home</Button>
            </Link>
          </CardContent>
        </Card>
      </div>
    );
  }

  if (error || !ticket || !event) {
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
        <Link href="/">
          <Button variant="ghost" className="mb-6">
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to Search
          </Button>
        </Link>

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
              <h3 className="font-semibold mb-4">Ticket Details</h3>
              <div className="space-y-2">
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Ticket Type:</span>
                  <span className="font-medium">{ticket.ticket_type_display_name}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Price:</span>
                  <span className="font-medium">{formatPrice(ticket.ticket_type_price_cents)}</span>
                </div>
              </div>
            </div>

            <div className="border-t pt-4">
              <div className="flex justify-between items-center mb-4">
                <span className="text-lg font-semibold">Total:</span>
                <span className="text-2xl font-bold">{formatPrice(ticket.ticket_type_price_cents)}</span>
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

