"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { isTicketAvailable } from "@/lib/utils";
import type { Ticket } from "@/types/events";
import type { SectionSeatsProps } from "./types";
import { SeatButton } from "./seat-button";

export function SectionSeats({ section, selectedTicketIds, onSeatSelect }: SectionSeatsProps) {
  const handleSeatClick = (ticket: Ticket) => {
    if (isTicketAvailable(ticket) && onSeatSelect) {
      onSeatSelect(ticket.id);
    }
  };

  const availableCount = section.tickets.filter(isTicketAvailable).length;
  const soldCount = section.tickets.filter((t) => t.status === "sold").length;
  const reservedCount = section.tickets.filter((t) => t.is_reserved).length;

  return (
    <Card className="border">
      <CardHeader>
        <CardTitle>{section.typeName}</CardTitle>
        <p className="text-sm text-muted-foreground">
          {availableCount} available • {reservedCount} reserved • {soldCount} sold
        </p>
      </CardHeader>
      <CardContent className="pt-6">
        <div className="grid grid-cols-10 gap-2">
          {section.tickets.map((ticket) => (
            <SeatButton
              key={ticket.id}
              ticket={ticket}
              isAvailable={isTicketAvailable(ticket)}
              isSelected={selectedTicketIds?.has(ticket.id) || false}
              colorConfig={section.colorConfig}
              onClick={() => handleSeatClick(ticket)}
            />
          ))}
        </div>
      </CardContent>
    </Card>
  );
}