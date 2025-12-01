/**
 * Event types matching backend API responses
 */

// Backend Event type (from core/types/types.go)
export interface Event {
  id: string;
  title: string;
  description: string;
  start_date: string; // ISO date string
  venue_id: string;
  created_at: string; // ISO date string
}

// Enriched search result with venue information (from SearchEventResult)
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

// Search results wrapper (if backend returns paginated results)
export interface SearchEventResults {
  results: SearchEventResult[];
  total: number;
}

// UI display model for search results
export interface SearchResult {
  id: string;
  title: string;
  location: string;
  date?: string;
  price?: number;
}

