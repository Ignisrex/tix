"use client";

import { useState } from "react";
import { SearchBar } from "@/components/search-bar";
import { SearchResults } from "@/components/search-results";
import { searchEvents } from "@/lib/api";
import type { SearchEventResult, SearchResult } from "@/types/events";
import Lottie from "lottie-react";
import ticketsAnimation from "@/assets/animations/ticket.json";

export default function Home() {
  const [isFocused, setIsFocused] = useState(false);
  const [query, setQuery] = useState("");
  const [results, setResults] = useState<SearchResult[]>([]);
  const [loading, setLoading] = useState(false);
  const [hasSearched, setHasSearched] = useState(false);
  const [showAnimation, setShowAnimation] = useState(true);

  function handleFocus() {
    setIsFocused(true);
  }

  async function handleSearch(searchQuery: string) {
    setQuery(searchQuery);
    setLoading(true);
    setResults([]);
    setHasSearched(true);
    setShowAnimation(false);

    try {
      const events = await searchEvents(searchQuery);
      
      // Map SearchEventResult to UI SearchResult format
      const mappedResults: SearchResult[] = events.map((event: SearchEventResult) => ({
        id: event.id,
        title: event.title,
        location: event.venue_location || `${event.venue_name || "Unknown venue"}`,
        date: event.start_date,
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
      <div className="absolute inset-0 bg-gradient-to-b from-background to-muted -z-10" />

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
      {/* animation just below search bar */}
      {showAnimation && (
        <div className="flex justify-center mt-4 transition-opacity duration-500 opacity-100">
          <Lottie
            animationData={ticketsAnimation}
            loop
            autoplay
            style={{ width: 300, height: 300 }}
          />
        </div>
      )}

      {/* Results area - appears when searching */}
      {isFocused && (loading || hasSearched) && (
        <div className="w-full px-4 sm:px-8 mt-8">
          <div
            className={`w-full max-w-7xl mx-auto transition-all duration-500 ease-out ${
              loading ? "opacity-0 translate-y-2" : "opacity-100 translate-y-0"
            }`}
          >
            <SearchResults query={query} results={results} loading={loading} />
          </div>
        </div>
      )}
    </main>
  );
}
