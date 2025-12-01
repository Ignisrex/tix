"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { getEvent, getEventTickets, reserveTicket } from "@/lib/api";
import type { Event, Ticket } from "@/types/events";
import type { ReservationData } from "@/types/booking";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { SeatView } from "@/components/seat-view";
import { ArrowLeft } from "lucide-react";
import Link from "next/link";
import { formatDateLong } from "@/lib/utils";

export default function EventDetailPage() {
  const params = useParams();
  const router = useRouter();
  const eventId = params.id as string;

  const [event, setEvent] = useState<Event | null>(null);
  const [tickets, setTickets] = useState<Ticket[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchData() {
      try {
        setLoading(true);
        const [eventData, ticketsData] = await Promise.all([
          getEvent(eventId),
          getEventTickets(eventId),
        ]);
        setEvent(eventData);
        setTickets(ticketsData);
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to load event");
      } finally {
        setLoading(false);
      }
    }

    if (eventId) {
      fetchData();
    }
  }, [eventId]);

  if (loading) {
    return (
      <div className="min-h-screen bg-background">
        <div className="container mx-auto px-4 py-8 max-w-6xl">
          <Skeleton className="h-10 w-32 mb-6" />
          <Skeleton className="h-12 w-3/4 mb-4" />
          <Skeleton className="h-64 w-full mb-8" />
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {[1, 2, 3].map((i) => (
              <Skeleton key={i} className="h-48 w-full" />
            ))}
          </div>
        </div>
      </div>
    );
  }

  if (error || !event) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <Card className="max-w-md">
          <CardHeader>
            <CardTitle>Error</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-muted-foreground mb-4">
              {error || "Event not found"}
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
      <div className="container mx-auto px-4 py-8 max-w-6xl">
        {/* Back button */}
        <Link href="/">
          <Button variant="ghost" className="mb-6">
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to Search
          </Button>
        </Link>

        {/* Event Header */}
        <div className="mb-8">
          <h1 className="text-4xl font-bold mb-4">{event.title}</h1>
          <div className="text-lg text-muted-foreground mb-4">
            {formatDateLong(event.start_date)}
          </div>
          <p className="text-base leading-relaxed">{event.description}</p>
        </div>

        {/* Seat View Section */}
        <div className="mb-8">
          <h2 className="text-2xl font-semibold mb-6">Select Your Seat</h2>

          {tickets.length === 0 ? (
            <Card>
              <CardContent className="pt-6">
                <p className="text-center text-muted-foreground">
                  No tickets available for this event
                </p>
              </CardContent>
            </Card>
          ) : (
            <SeatView
              tickets={tickets}
              onSeatSelect={async (ticketId) => {
                try {
                  const response = await reserveTicket(ticketId);
                  
                  if (response.success) {
                    // Store reservation in localStorage
                    const reservationData: ReservationData = {
                      ticketId: response.ticket_id,
                      eventId: eventId,
                      reservedAt: Date.now(),
                    };
                    localStorage.setItem("tix_reservation", JSON.stringify(reservationData));
                    
                    // Route to checkout
                    router.push(`/checkout/${response.ticket_id}`);
                  } else {
                    // Show alert for failure
                    alert(response.message || "This seat is not available");
                  }
                } catch (err) {
                  const errorMessage = err instanceof Error ? err.message : "Failed to reserve ticket";
                  alert(errorMessage || "This seat is not available");
                }
              }}
            />
          )}
        </div>
      </div>
    </div>
  );
}


