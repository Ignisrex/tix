import type { Ticket } from "@/types/events";

export interface SeatViewProps {
  tickets: Ticket[];
  selectedTicketIds?: Set<string>;
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
  selectedTicketIds?: Set<string>;
  onSeatSelect?: (ticketId: string) => void;
}

export interface SeatButtonProps {
  ticket: Ticket;
  isAvailable: boolean;
  isSelected: boolean;
  colorConfig: ColorConfig;
  onClick: () => void;
}

