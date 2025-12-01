"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { getEvent, getEventTickets, reserveTickets } from "@/lib/api";
import type { Event, Ticket } from "@/types/events";
import type { ReservationData } from "@/types/booking";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { SeatView } from "@/components/seat-view";
import { SeatSelectionDrawer } from "@/components/seat-selection-drawer";
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
  const [selectedTickets, setSelectedTickets] = useState<Ticket[]>([]);
  const [isReserving, setIsReserving] = useState(false);

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

        <SeatSelectionDrawer
          selectedTickets={selectedTickets}
          onRemoveTicket={(ticketId) => {
            setSelectedTickets((prev) => prev.filter((t) => t.id !== ticketId));
          }}
          onProceedToCheckout={async () => {
            if (selectedTickets.length === 0) return;

            // Check if all selected tickets are from the same event
            const allSameEvent = selectedTickets.every((t) => t.event_id === eventId);
            if (!allSameEvent) {
              alert("All selected tickets must be from the same event");
              return;
            }

            setIsReserving(true);
            try {
              // Check for existing reservation
              const existingReservationStr = localStorage.getItem("tix_reservation");
              let existingTicketIds: string[] = [];
              
              if (existingReservationStr) {
                try {
                  const existingReservation: ReservationData = JSON.parse(existingReservationStr);
                  
                  // Check if existing reservation is for the same event
                  if (existingReservation.eventId === eventId) {
                    // Get existing ticket IDs
                    existingTicketIds = existingReservation.ticketIds || [];
                  } else {
                    // Different event - user should clear existing reservation first
                    alert("You already have tickets reserved for a different event. Please complete or cancel that reservation first.");
                    setIsReserving(false);
                    return;
                  }
                } catch (err) {
                  // Invalid existing reservation, ignore it
                  console.error("Error parsing existing reservation:", err);
                }
              }

              // Get new ticket IDs (only tickets that aren't already reserved)
              const newTicketIds = selectedTickets.map((t) => t.id);
              const ticketsToReserve = newTicketIds.filter((id) => !existingTicketIds.includes(id));

              // If all tickets are already reserved, just update localStorage and go to checkout
              if (ticketsToReserve.length === 0) {
                // All tickets already reserved, just update the reservation data
                const reservationData: ReservationData = {
                  ticketIds: [...new Set([...existingTicketIds, ...newTicketIds])],
                  eventId: eventId,
                  reservedAt: Date.now(),
                };
                localStorage.setItem("tix_reservation", JSON.stringify(reservationData));
                router.push("/checkout");
                setIsReserving(false);
                return;
              }

              // Reserve only the new tickets
              const response = await reserveTickets(ticketsToReserve);

              if (response.success) {
                // Merge existing and newly reserved tickets
                const allReservedTicketIds = [...new Set([...existingTicketIds, ...response.ticket_ids])];
                
                // Store merged reservation in localStorage
                const reservationData: ReservationData = {
                  ticketIds: allReservedTicketIds,
                  eventId: eventId,
                  reservedAt: Date.now(),
                };
                localStorage.setItem("tix_reservation", JSON.stringify(reservationData));

                // Route to checkout
                router.push("/checkout");
              } else {
                // Show alert for failure and redirect back
                alert(response.message || "One or more seats are not available");
                router.push(`/events/${eventId}`);
              }
            } catch (err) {
              const errorMessage = err instanceof Error ? err.message : "Failed to reserve tickets";
              alert(errorMessage || "One or more seats are not available");
              router.push(`/events/${eventId}`);
            } finally {
              setIsReserving(false);
            }
          }}
          isReserving={isReserving}
        />

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
              selectedTicketIds={new Set(selectedTickets.map((t) => t.id))}
              onSeatSelect={(ticketId) => {
                const ticket = tickets.find((t) => t.id === ticketId);
                if (!ticket) return;

                // Check if ticket is already selected
                const isSelected = selectedTickets.some((t) => t.id === ticketId);
                if (isSelected) {
                  // Remove from selection
                  setSelectedTickets((prev) => prev.filter((t) => t.id !== ticketId));
                } else {
                  // Check if ticket is from the same event as current event
                  if (ticket.event_id !== eventId) {
                    alert("You cannot select tickets from different events");
                    return;
                  }
                  
                  // If there are already selected tickets, ensure they're all from the same event
                  if (selectedTickets.length > 0) {
                    const firstSelectedEventId = selectedTickets[0].event_id;
                    if (ticket.event_id !== firstSelectedEventId) {
                      alert("You cannot select tickets from different events. Please clear your current selection first.");
                      return;
                    }
                  }
                  
                  // Add to selection
                  setSelectedTickets((prev) => [...prev, ticket]);
                }
              }}
            />
          )}
        </div>
      </div>
    </div>
  );
}


