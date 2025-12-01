import { useEffect, useState, useRef } from 'react';
import { createTicketStream } from '@/lib/api/events';
import type { Ticket } from '@/types/events';

interface UseTicketStreamResult {
  tickets: Ticket[];
  connected: boolean;
  error: Error | null;
}

/**
 * Hook to stream real-time ticket updates via Server-Sent Events
 * Automatically reconnects on disconnect with exponential backoff
 */
export function useTicketStream(eventId: string | null): UseTicketStreamResult {
  const [tickets, setTickets] = useState<Ticket[]>([]);
  const [connected, setConnected] = useState(false);
  const [error, setError] = useState<Error | null>(null);
  
  const eventSourceRef = useRef<EventSource | null>(null);
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const reconnectAttemptsRef = useRef(0);
  const maxReconnectAttempts = 5;
  const baseReconnectDelay = 1000; // 1 second

  useEffect(() => {
    if (!eventId) {
      return;
    }

    let isMounted = true;

    const connect = () => {
      // Clean up any existing connection
      if (eventSourceRef.current) {
        eventSourceRef.current.close();
      }

      // Clear any pending reconnect
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
        reconnectTimeoutRef.current = null;
      }

      try {
        const eventSource = createTicketStream(eventId);
        eventSourceRef.current = eventSource;

        eventSource.onopen = () => {
          if (isMounted) {
            setConnected(true);
            setError(null);
            reconnectAttemptsRef.current = 0; // Reset on successful connection
          }
        };

        eventSource.onmessage = (event) => {
          if (isMounted) {
            try {
              const data = JSON.parse(event.data) as Ticket[];
              setTickets(data);
              setError(null);
            } catch (err) {
              console.error('Error parsing SSE data:', err);
              setError(err instanceof Error ? err : new Error('Failed to parse ticket data'));
            }
          }
        };

        eventSource.onerror = (err) => {
          if (isMounted) {
            setConnected(false);
            
            // Only attempt reconnect if connection was open (not initial connection failure)
            if (eventSource.readyState === EventSource.CLOSED) {
              if (reconnectAttemptsRef.current < maxReconnectAttempts) {
                const delay = baseReconnectDelay * Math.pow(2, reconnectAttemptsRef.current);
                reconnectAttemptsRef.current++;
                
                reconnectTimeoutRef.current = setTimeout(() => {
                  if (isMounted) {
                    connect();
                  }
                }, delay);
              } else {
                setError(new Error('Failed to connect after multiple attempts'));
              }
            } else {
              // Connection error but still trying
              setError(new Error('Connection error'));
            }
          }
        };
      } catch (err) {
        if (isMounted) {
          setError(err instanceof Error ? err : new Error('Failed to create event stream'));
          setConnected(false);
        }
      }
    };

    // Initial connection
    connect();

    // Cleanup on unmount
    return () => {
      isMounted = false;
      
      if (eventSourceRef.current) {
        eventSourceRef.current.close();
        eventSourceRef.current = null;
      }
      
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
        reconnectTimeoutRef.current = null;
      }
    };
  }, [eventId]);

  return { tickets, connected, error };
}

