"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import type { Ticket } from "@/types/events";
import type { SectionSeatsProps } from "./types";
import { SeatButton } from "./seat-button";

export function SectionSeats({ section, onSeatSelect }: SectionSeatsProps) {
  const handleSeatClick = (ticket: Ticket) => {
    if (ticket.status === "available" && onSeatSelect) {
      onSeatSelect(ticket.id);
    }
  };

  const availableCount = section.tickets.filter((t) => t.status === "available").length;
  const soldCount = section.tickets.filter((t) => t.status === "sold").length;

  return (
    <Card className="border">
      <CardHeader>
        <CardTitle>{section.typeName}</CardTitle>
        <p className="text-sm text-muted-foreground">
          {availableCount} available • {soldCount} sold
        </p>
      </CardHeader>
      <CardContent className="pt-6">
        <div className="grid grid-cols-10 gap-2">
          {section.tickets.map((ticket) => {
            const isAvailable = ticket.status === "available";

            return (
              <SeatButton
                key={ticket.id}
                ticket={ticket}
                isAvailable={isAvailable}
                colorConfig={section.colorConfig}
                onClick={() => handleSeatClick(ticket)}
              />
            );
          })}
        </div>
      </CardContent>
    </Card>
  );
}

