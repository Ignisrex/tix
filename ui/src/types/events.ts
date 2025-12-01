/**
 * Event types matching backend API responses
 */


export interface Event {
  id: string;
  title: string;
  description: string;
  start_date: string; // ISO date string
  venue_id: string;
  created_at: string; // ISO date string
}

export interface SearchEventResult {
  id: string;
  title: string;
  description: string;
  start_date: string;
  venue_id: string;
  venue_name: string;
  venue_location: string;
  created_at: string;
}


export interface SearchEventResults {
  results: SearchEventResult[];
  total: number;
}


export interface SearchResult {
  id: string;
  title: string;
  location: string;
  date?: string;
  price?: number;
}


export type TicketStatus = "available" | "sold";

export interface Ticket {
  id: string;
  event_id: string;
  ticket_type_id: string;
  status: TicketStatus;
  ticket_type_name: string;
  ticket_type_display_name: string;
  ticket_type_price_cents: number;
}

export interface TicketType {
  id: string;
  name: string;
  display_name: string;
  price_cents: number;
}

export interface TicketWithType extends Ticket {
  type_name?: string;
  price_cents?: number;
}

