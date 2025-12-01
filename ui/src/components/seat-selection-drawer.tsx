"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import {CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Drawer } from "@/components/ui/drawer";
import { ShoppingCart, X, ChevronRight } from "lucide-react";
import type { Ticket } from "@/types/events";
import { formatPrice } from "@/lib/utils";
import { cn } from "@/lib/utils";

interface SeatSelectionDrawerProps {
  selectedTickets: Ticket[];
  onRemoveTicket: (ticketId: string) => void;
  onProceedToCheckout: () => void;
  isReserving: boolean;
}

export function SeatSelectionDrawer({
  selectedTickets,
  onRemoveTicket,
  onProceedToCheckout,
  isReserving,
}: SeatSelectionDrawerProps) {
  const [isOpen, setIsOpen] = useState(false);

  if (selectedTickets.length === 0) {
    return null;
  }

  const totalCents = selectedTickets.reduce(
    (sum, ticket) => sum + (ticket.ticket_type_price_cents || 0),
    0
  );

  return (
    <>
      {/* Tab that shows when items are selected */}
      <button
        onClick={() => setIsOpen(true)}
        className={cn(
          "fixed right-0 top-1/2 -translate-y-1/2 z-30",
          "bg-indigo-500 hover:bg-indigo-600 text-white",
          "rounded-l-lg shadow-lg transition-all duration-200",
          "px-3 py-4 flex items-center gap-2",
          "hover:px-4",
          isOpen && "opacity-0 pointer-events-none"
        )}
        aria-label="Open cart"
      >
        <ShoppingCart className="h-5 w-5" />
        <span className="font-semibold text-sm">{selectedTickets.length}</span>
        <ChevronRight className="h-4 w-4" />
      </button>

      {/* Drawer */}
      <Drawer isOpen={isOpen} onClose={() => setIsOpen(false)}>
        <div className="h-full flex flex-col">
          {/* Header */}
          <CardHeader className="border-b pb-4">
            <div className="flex items-center justify-between">
              <CardTitle className="text-lg">
                Selected Seats ({selectedTickets.length})
              </CardTitle>
              <Button
                variant="ghost"
                size="sm"
                onClick={() => setIsOpen(false)}
                className="h-8 w-8 p-0"
              >
                <X className="h-4 w-4" />
              </Button>
            </div>
          </CardHeader>

          {/* Content */}
          <CardContent className="flex-1 overflow-y-auto pt-6">
            <div className="space-y-3">
              {selectedTickets.map((ticket) => (
                <div
                  key={ticket.id}
                  className="flex items-center justify-between p-3 rounded-lg bg-muted/50"
                >
                  <div className="flex-1">
                    <p className="font-medium text-sm">
                      {ticket.ticket_type_display_name || ticket.ticket_type_name}
                    </p>
                  </div>
                  <div className="flex items-center gap-3">
                    <p className="font-semibold text-sm">
                      {formatPrice(ticket.ticket_type_price_cents || 0)}
                    </p>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => onRemoveTicket(ticket.id)}
                      className="h-8 w-8 p-0 text-muted-foreground hover:text-destructive"
                    >
                      <X className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>

          {/* Footer */}
          <div className="border-t p-6 space-y-4">
            <div className="flex items-center justify-between">
              <span className="text-lg font-semibold">Total</span>
              <span className="text-lg font-bold">{formatPrice(totalCents)}</span>
            </div>
            <Button
              onClick={onProceedToCheckout}
              disabled={isReserving || selectedTickets.length === 0}
              className="w-full bg-indigo-500 hover:bg-indigo-600 text-white"
              size="lg"
            >
              {isReserving ? "Reserving..." : "Proceed to Checkout"}
            </Button>
          </div>
        </div>
      </Drawer>
    </>
  );
}

