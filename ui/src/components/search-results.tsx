"use client";

import { SearchX } from "lucide-react";
import { ResultCard } from "./result-card";
import { ResultSkeleton } from "./result-skeleton";
import type { SearchResult } from "@/types/events";

interface SearchResultsProps {
  query: string;
  results: SearchResult[];
  loading: boolean;
}

export function SearchResults({ query, results, loading }: SearchResultsProps) {
  if (loading) {
    return (
      <div className="animate-pulse">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
          {Array.from({ length: 8 }).map((_, i) => (
            <ResultSkeleton key={i} />
          ))}
        </div>
      </div>
    );
  }

  if (results.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-16">
        <SearchX className="h-12 w-12 text-muted-foreground mb-4" />
        <h3 className="text-lg font-semibold mb-2">No results found</h3>
        <p className="text-sm text-muted-foreground text-center">
          Try searching for something else or check your spelling.
        </p>
      </div>
    );
  }

  return (
    <div>
      <div className="mb-4">
        <p className="text-sm text-muted-foreground">
          Found {results.length} {results.length === 1 ? "result" : "results"} for &quot;{query}&quot;
        </p>
      </div>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
        {results.map((result) => (
          <ResultCard
            key={result.id}
            id={result.id}
            title={result.title}
            location={result.location}
            date={result.date}
            price={result.price}
          />
        ))}
      </div>
    </div>
  );
}

