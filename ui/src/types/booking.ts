// ui/src/types/booking.ts
export interface ReserveResponse {
  success: boolean;
  message: string;
  ticket_id: string;
}

export interface PurchaseResponse {
  success: boolean;
  message: string;
  ticket_id: string;
  total: number;
}