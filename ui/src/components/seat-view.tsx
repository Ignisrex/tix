"use client";

import { useMemo } from "react";
import type { Ticket } from "@/types/events";
import type { SeatViewProps, TicketTypeSection } from "./seat-view/types";
import { TYPE_COLORS, SECTION_ORDER } from "./seat-view/constants";
import { SectionSeats } from "./seat-view/section-seats";

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
      const aIndex = SECTION_ORDER.indexOf(a.sortKey);
      const bIndex = SECTION_ORDER.indexOf(b.sortKey);
      const safeA = aIndex === -1 ? SECTION_ORDER.length : aIndex;
      const safeB = bIndex === -1 ? SECTION_ORDER.length : bIndex;
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
