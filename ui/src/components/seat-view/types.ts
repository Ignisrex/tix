import type { Ticket } from "@/types/events";

export interface SeatViewProps {
  tickets: Ticket[];
  onSeatSelect?: (ticketId: string) => void;
}

export interface TicketTypeSection {
  typeId: string;
  typeName: string;
  tickets: Ticket[];
  colorConfig: ColorConfig;
}

export interface ColorConfig {
  color: string;
  bgColor: string;
  borderColor: string;
}

export interface SectionSeatsProps {
  section: TicketTypeSection;
  onSeatSelect?: (ticketId: string) => void;
}

export interface SeatButtonProps {
  ticket: Ticket;
  isAvailable: boolean;
  colorConfig: ColorConfig;
  onClick: () => void;
}

