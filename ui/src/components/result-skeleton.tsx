import { Skeleton } from "@/components/ui/skeleton";
import { Card } from "@/components/ui/card";

export function ResultSkeleton() {
  return (
    <Card className="overflow-hidden">
      <div className="aspect-video relative">
        <Skeleton className="w-full h-full" />
      </div>
      <div className="p-4 space-y-3">
        <Skeleton className="h-6 w-3/4" />
        <Skeleton className="h-4 w-1/2" />
        <div className="flex justify-between items-center">
          <Skeleton className="h-4 w-24" />
          <Skeleton className="h-5 w-16" />
        </div>
      </div>
    </Card>
  );
}

