import { Card, CardContent } from "@/components/ui/card";
import Link from "next/link";
import { formatDateShort } from "@/lib/utils";

interface ResultCardProps {
  id: string;
  title: string;
  location: string;
  date?: string;
}

export function ResultCard({ id, title, location, date }: ResultCardProps) {
  return (
    <Link href={`/events/${id}`}>
      <Card className="overflow-hidden hover:shadow-lg transition-shadow cursor-pointer h-full gap-0 p-0">
        <div className="aspect-video relative bg-muted">
          {/* Placeholder for image */}
          <div className="w-full h-full flex items-center justify-center text-muted-foreground">
            Event Image
          </div>
        </div>
        <CardContent className="p-4">
          <h3 className="font-semibold text-lg mb-2 line-clamp-2">{title}</h3>
          <p className="text-sm text-muted-foreground mb-3">{location}</p>
          {date && (
            <span className="text-sm font-medium">{formatDateShort(date)}</span>
          )}
        </CardContent>
      </Card>
    </Link>
  );
}

