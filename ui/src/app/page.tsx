"use client";

import { useState } from "react";
import { SearchBar } from "@/components/search-bar";
import { SearchResults } from "@/components/search-results";
import { searchEvents } from "@/lib/api";
import type { Event, SearchResult } from "@/types/events";

export default function Home() {
  const [isFocused, setIsFocused] = useState(false);
  const [query, setQuery] = useState("");
  const [results, setResults] = useState<SearchResult[]>([]);
  const [loading, setLoading] = useState(false);

  function handleFocus() {
    setIsFocused(true);
  }

  async function handleSearch(searchQuery: string) {
    setQuery(searchQuery);
    setLoading(true);
    setResults([]);

    try {
      const events = await searchEvents(searchQuery);
      
      // Map backend Event format to UI SearchResult format
      const mappedResults: SearchResult[] = events.map((event: Event) => ({
        id: event.id,
        title: event.title,
        location: "Location TBD", // TODO: Fetch venue info or use SearchEventResult if available
        date: event.start_date,
        price: undefined, // TODO: Fetch ticket prices if needed
      }));

      setResults(mappedResults);
    } catch (error) {
      console.error("Failed to search events:", error);
      setResults([]);
    } finally {
      setLoading(false);
    }
  }

  return (
    <main className="min-h-screen flex flex-col">
      {/* Background gradient - only shows when not focused */}
      {!isFocused && (
        <div className="absolute inset-0 bg-gradient-to-b from-background to-muted -z-10" />
      )}

      {/* Search bar container - moves up slightly when hero fades */}
      <div
        className={`w-full flex justify-center px-4 sm:px-8 transition-all duration-1000 ease-[cubic-bezier(0.4,0,0.2,1)] ${
          isFocused ? "pt-4 -translate-y-16" : "pt-24 translate-y-0"
        }`}
      >
        <div className="w-full max-w-3xl space-y-10">
          {/* Hero section - fades out smoothly when focused */}
          <div
            className={`text-center mb-10 transition-opacity duration-1000 ease-[cubic-bezier(0.4,0,0.2,1)] ${
              isFocused ? "opacity-0 pointer-events-none" : "opacity-100"
            }`}
          >
            <h1 className="text-3xl sm:text-4xl font-semibold">
              Find your next event
            </h1>
            <p className="mt-2 text-sm text-muted-foreground">
              Search by city, venue, or event name.
            </p>
          </div>

          <div>
            <SearchBar onFocus={handleFocus} onSearch={handleSearch} />
          </div>
        </div>
      </div>

      {/* Results area - appears when searching */}
      {isFocused && (loading || results.length > 0) && (
        <div className="w-full px-4 sm:px-8 mt-8">
          <div className="w-full max-w-7xl mx-auto">
            <SearchResults query={query} results={results} loading={loading} />
          </div>
        </div>
      )}
    </main>
  );
}
