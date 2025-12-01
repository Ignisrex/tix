"use client";

import { useState } from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Search } from "lucide-react";

interface SearchBarProps {
  onFocus?: () => void;
  onSearch?: (query: string) => void;
  query?: string;
}

export function SearchBar({ onFocus, onSearch, query: controlledQuery }: SearchBarProps) {
  const [internalQuery, setInternalQuery] = useState("");
  const query = controlledQuery ?? internalQuery;

  function handleInputChange(e: React.ChangeEvent<HTMLInputElement>) {
    const value = e.target.value;
    if (controlledQuery === undefined) {
      setInternalQuery(value);
    }
  }

  function handleFocus() {
    onFocus?.();
  }

  function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (query.trim() && onSearch) {
      onSearch(query.trim());
    }
  }

  return (
    <form
      onSubmit={onSubmit}
      className="flex items-center gap-2 rounded-full border bg-background/80 px-2 py-1 shadow-md backdrop-blur-md max-w-3xl w-full"
    >
      <div className="flex flex-1 items-center px-4 py-2 hover:bg-muted/60 rounded-full transition-colors">
        <Input
          value={query}
          onChange={handleInputChange}
          onFocus={handleFocus}
          placeholder="Search events, venues, artists..."
          className="h-8 border-0 p-0 text-base shadow-none focus-visible:ring-0 focus-visible:ring-offset-0"
        />
      </div>

      <Button
        type="submit"
        className="ml-2 rounded-full bg-rose-500 hover:bg-rose-600 px-3 h-11 flex items-center gap-2"
      >
        <Search className="h-4 w-4" />
        <span className="hidden sm:inline text-sm font-semibold">Search</span>
      </Button>
    </form>
  );
}