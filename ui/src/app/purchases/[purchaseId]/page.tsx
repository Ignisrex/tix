"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { getPurchaseDetails } from "@/lib/api/booking";
import { getEvent } from "@/lib/api";
import type { Event } from "@/types/events";
import type { PurchaseDetailsResponse } from "@/types/booking";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { CheckCircle2, Home } from "lucide-react";
import Link from "next/link";
import { formatPrice, formatDateLong } from "@/lib/utils";

export default function PurchaseSuccessPage() {
  const params = useParams();
  const purchaseId = params.purchaseId as string;

  const [purchase, setPurchase] = useState<PurchaseDetailsResponse | null>(null);
  const [event, setEvent] = useState<Event | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchPurchaseDetails() {
      if (!purchaseId) {
        setError("Missing purchase ID");
        setLoading(false);
        return;
      }

      try {
        setLoading(true);
        const purchaseData = await getPurchaseDetails(purchaseId);
        setPurchase(purchaseData);

        // Fetch event details if we have tickets
        if (purchaseData.tickets && purchaseData.tickets.length > 0) {
          const eventId = purchaseData.tickets[0].event_id;
          try {
            const eventData = await getEvent(eventId);
            setEvent(eventData);
          } catch (err) {
            console.error("Failed to fetch event:", err);
            // Continue without event data
          }
        }
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to load purchase details");
      } finally {
        setLoading(false);
      }
    }

    fetchPurchaseDetails();
  }, [purchaseId]);

  if (loading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="container mx-auto px-4 py-8 max-w-2xl">
          <Skeleton className="h-64 w-full" />
        </div>
      </div>
    );
  }

  if (error || !purchase) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <Card className="max-w-md">
          <CardHeader>
            <CardTitle>Error</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-muted-foreground mb-4">
              {error || "Failed to load purchase details"}
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
            Your {purchase.tickets.length === 1 ? "ticket has" : "tickets have"} been confirmed
          </p>
        </div>

        {event && (
          <Card className="mb-6">
            <CardHeader>
              <CardTitle>Event Details</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <h2 className="text-2xl font-bold mb-2">{event.title}</h2>
                <p className="text-muted-foreground mb-4">{event.description}</p>
                <p className="text-sm text-muted-foreground">
                  {formatDateLong(event.start_date)}
                </p>
              </div>
            </CardContent>
          </Card>
        )}

        <Card className="mb-6">
          <CardHeader>
            <CardTitle>Purchase Details</CardTitle>
          </CardHeader>
          <CardContent className="space-y-6">
            <div className="space-y-3">
              <div className="flex justify-between">
                <span className="text-muted-foreground">Purchase ID:</span>
                <span className="font-mono text-sm">{purchase.purchase_id}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Purchase Date:</span>
                <span className="font-medium">
                  {new Date(purchase.purchase_created_at).toLocaleString()}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Total:</span>
                <span className="font-medium">{formatPrice(purchase.total_cents)}</span>
              </div>
            </div>

            <div className="border-t pt-4">
              <h3 className="font-semibold mb-4">
                Tickets ({purchase.tickets.length})
              </h3>
              <div className="space-y-3">
                {purchase.tickets.map((ticket) => (
                  <div
                    key={ticket.id}
                    className="flex items-center justify-between p-3 rounded-lg bg-muted/50"
                  >
                    <div className="flex-1">
                      <p className="font-medium text-sm">
                        {ticket.ticket_type_display_name || ticket.ticket_type_name}
                      </p>
                      <p className="text-xs text-muted-foreground mt-1 font-mono">
                        {ticket.id}
                      </p>
                    </div>
                    <div className="text-right">
                      <p className="font-semibold text-sm">
                        {formatPrice(ticket.ticket_type_price_cents)}
                      </p>
                    </div>
                  </div>
                ))}
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

