// ui/src/types/booking.ts
export interface ReserveRequest {
  ticket_ids: string[];
}

export interface ReserveResponse {
  success: boolean;
  message: string;
  ticket_ids: string[];
}

export interface ReservationData {
  ticketIds: string[];
  eventId: string;
  reservedAt: number; // timestamp in milliseconds
}

export interface PurchaseRequest {
  ticket_ids: string[];
}

export interface PurchaseResponse {
  success: boolean;
  message: string;
  ticket_ids: string[];
  total: number;
  purchase_id: string;
}

export interface PurchaseTicketDetail {
  id: string;
  event_id: string;
  ticket_type_id: string;
  status: string;
  ticket_type_name: string;
  ticket_type_display_name: string;
  ticket_type_price_cents: number;
}

export interface PurchaseDetailsResponse {
  purchase_id: string;
  total_cents: number;
  purchase_created_at: string;
  tickets: PurchaseTicketDetail[];
}