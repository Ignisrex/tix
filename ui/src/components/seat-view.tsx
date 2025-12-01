"use client";

import { useState, useMemo } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import type { Ticket } from "@/types/events";

interface SeatViewProps {
  tickets: Ticket[];
  onSeatSelect?: (ticketId: string) => void;
}

interface TicketTypeSection {
  typeId: string;
  typeName: string;
  tickets: Ticket[];
  colorConfig: {
    color: string;
    bgColor: string;
    borderColor: string;
  };
}

// Color scheme for each ticket type section.
// We assume only these three ticket types exist for now: VIP, Front Row, GA.
const TYPE_COLORS: Record<string, { color: string; bgColor: string; borderColor: string }> = {
  vip: {
    color: "text-purple-700",
    bgColor: "bg-purple-100",
    borderColor: "border-purple-300",
  },
  front_row: {
    color: "text-blue-700",
    bgColor: "bg-blue-100",
    borderColor: "border-blue-300",
  },
  ga: {
    color: "text-green-700",
    bgColor: "bg-green-100",
    borderColor: "border-green-300",
  },
};

export function SeatView({ tickets, onSeatSelect }: SeatViewProps) {
  // Group tickets by ticket_type_id
  const ticketsByType = useMemo(() => {
    const grouped = tickets.reduce((acc, ticket) => {
      const key = ticket.ticket_type_id;
      if (!acc[key]) {
        acc[key] = [];
      }
      acc[key].push(ticket);
      return acc;
    }, {} as Record<string, Ticket[]>);

    // Sort tickets by ID within each type for consistent ordering
    Object.keys(grouped).forEach((typeId) => {
      grouped[typeId].sort((a, b) => a.id.localeCompare(b.id));
    });

    return grouped;
  }, [tickets]);

  // Create sections with color mapping
  const sections: TicketTypeSection[] = useMemo(() => {
    const ORDER: Array<"vip" | "front_row" | "ga"> = ["vip", "front_row", "ga"];

    const sectionsWithSortKey = Object.entries(ticketsByType).map(([typeId, tickets]) => {
      const firstTicket = tickets[0];
      const key = firstTicket.ticket_type_name as "vip" | "front_row" | "ga";

      const colorConfig = TYPE_COLORS[key] ?? TYPE_COLORS.ga;
      const displayName = firstTicket.ticket_type_display_name;

      return {
        section: {
          typeId,
          typeName: displayName,
          tickets,
          colorConfig,
        },
        sortKey: key,
      };
    });

    // Sort sections so that VIP comes first, then Front Row, then GA.
    sectionsWithSortKey.sort((a, b) => {
      const aIndex = ORDER.indexOf(a.sortKey);
      const bIndex = ORDER.indexOf(b.sortKey);
      const safeA = aIndex === -1 ? ORDER.length : aIndex;
      const safeB = bIndex === -1 ? ORDER.length : bIndex;
      return safeA - safeB;
    });

    return sectionsWithSortKey.map((item) => item.section);
  }, [ticketsByType]);

  return (
    <div className="space-y-8">
      {sections.map((section) => (
        <SectionSeats
          key={section.typeId}
          section={section}
          onSeatSelect={onSeatSelect}
        />
      ))}
    </div>
  );
}

interface SectionSeatsProps {
  section: TicketTypeSection;
  onSeatSelect?: (ticketId: string) => void;
}

function SectionSeats({ section, onSeatSelect }: SectionSeatsProps) {
  const [selectedSeat, setSelectedSeat] = useState<string | null>(null);

  const handleSeatClick = (ticket: Ticket) => {
    if (ticket.status === "available" && onSeatSelect) {
      setSelectedSeat(ticket.id);
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
            const isSelected = selectedSeat === ticket.id;

            return (
              <button
                key={ticket.id}
                onClick={() => handleSeatClick(ticket)}
                disabled={!isAvailable}
                className={`
                  aspect-square rounded-md border-2 transition-all
                  ${isAvailable 
                    ? `${section.colorConfig.bgColor} ${section.colorConfig.borderColor} hover:scale-105 hover:shadow-md cursor-pointer ${isSelected ? 'ring-2 ring-offset-2 ring-primary scale-105' : ''}` 
                    : 'bg-gray-300 border-gray-400 opacity-50 cursor-not-allowed'
                  }
                  flex items-center justify-center text-xs font-medium
                  ${isAvailable ? section.colorConfig.color : 'text-gray-500'}
                `}
                title={isAvailable ? `Seat ${ticket.id.slice(0, 8)} - Available` : `Seat ${ticket.id.slice(0, 8)} - Sold`}
              >
                {ticket.id.slice(0, 4)}
              </button>
            );
          })}
        </div>
      </CardContent>
    </Card>
  );
}

