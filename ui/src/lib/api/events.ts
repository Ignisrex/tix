import { request } from './client';
import type { Event, SearchEventResult, Ticket } from '@/types/events';

export async function searchEvents(
  query?: string,
  limit: number = 20,
  offset: number = 0
): Promise<SearchEventResult[]> {
  const params: Record<string, string | number> = {
    limit,
    offset,
  };

  if (query) {
    params.q = query;
  }

  return request<SearchEventResult[]>('/events', { params });
}

export async function getEvent(id: string): Promise<Event> {
  return request<Event>(`/events/${id}`);
}

/**
 * Get tickets for an event
 */
export async function getEventTickets(eventId: string): Promise<Ticket[]> {
  return request<Ticket[]>(`/events/${eventId}/tickets`);
}

