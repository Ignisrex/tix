"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { getEvent, getTicket } from "@/lib/api";
import type { Event, Ticket } from "@/types/events";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { CheckCircle2, Home } from "lucide-react";
import Link from "next/link";
import { formatPrice, formatDateLong } from "@/lib/utils";

export default function SuccessPage() {
  const params = useParams();
  const eventId = params.eventId as string;
  const ticketId = params.ticketId as string;

  const [event, setEvent] = useState<Event | null>(null);
  const [ticket, setTicket] = useState<Ticket | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchTicketDetails() {
      if (!eventId || !ticketId) {
        setError("Missing event or ticket ID");
        setLoading(false);
        return;
      }

      try {
        setLoading(true);
        const [eventData, ticketData] = await Promise.all([
          getEvent(eventId),
          getTicket(eventId, ticketId),
        ]);

        setEvent(eventData);
        setTicket(ticketData);
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to load ticket details");
      } finally {
        setLoading(false);
      }
    }

    fetchTicketDetails();
  }, [eventId, ticketId]);


  if (loading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="container mx-auto px-4 py-8 max-w-2xl">
          <Skeleton className="h-64 w-full" />
        </div>
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
              <Button variant="outline">
                <Home className="mr-2 h-4 w-4" />
                Go Home
              </Button>
            </Link>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background">
      <div className="container mx-auto px-4 py-8 max-w-2xl">
        <div className="flex flex-col items-center text-center mb-8">
          <div className="rounded-full bg-green-100 dark:bg-green-900/20 p-4 mb-6">
            <CheckCircle2 className="h-16 w-16 text-green-600 dark:text-green-400" />
          </div>
          <h1 className="text-4xl font-bold mb-2">Payment Successful!</h1>
          <p className="text-lg text-muted-foreground">
            Your ticket has been confirmed
          </p>
        </div>

        <Card className="mb-6">
          <CardHeader>
            <CardTitle>Ticket Details</CardTitle>
          </CardHeader>
          <CardContent className="space-y-6">
            <div>
              <h2 className="text-2xl font-bold mb-2">{event.title}</h2>
              <p className="text-muted-foreground mb-4">{event.description}</p>
              <p className="text-sm text-muted-foreground">
                {formatDateLong(event.start_date)}
              </p>
            </div>

            <div className="border-t pt-4">
              <div className="space-y-3">
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Ticket Type:</span>
                  <span className="font-medium">{ticket.ticket_type_display_name}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Ticket ID:</span>
                  <span className="font-mono text-sm">{ticket.id}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Price:</span>
                  <span className="font-medium">{formatPrice(ticket.ticket_type_price_cents)}</span>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        <div className="flex flex-col sm:flex-row gap-4 justify-center">
          <Link href="/" className="w-full sm:w-auto">
            <Button className="w-full bg-indigo-500 hover:bg-indigo-600 text-white">
              <Home className="mr-2 h-4 w-4" />
              Back to Home
            </Button>
          </Link>
        </div>
      </div>
    </div>
  );
}

