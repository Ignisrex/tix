import type { ColorConfig } from "./types";

// Color scheme for each ticket type section.
// We assume only these three ticket types exist for now: VIP, Front Row, GA.
export const TYPE_COLORS: Record<string, ColorConfig> = {
  vip: {
    color: "text-purple-700",
    bgColor: "bg-purple-100",
    borderColor: "border-purple-300",
  },
  front_row: {
    color: "text-blue-700",
    bgColor: "bg-blue-100",
    borderColor: "border-blue-300",
  },
  ga: {
    color: "text-green-700",
    bgColor: "bg-green-100",
    borderColor: "border-green-300",
  },
};

export const SECTION_ORDER: Array<"vip" | "front_row" | "ga"> = ["vip", "front_row", "ga"];

