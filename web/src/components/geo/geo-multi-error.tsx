import { useState } from "react";
import type { BulkError } from "@/api/types.gen";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";

function GeoMultiErrorInner({ value, className }: { value: BulkError[]; className?: string }) {
  const [showErrors, setShowErrors] = useState(false);

  return (
    <div className={className}>
      {!showErrors ? (
        <Alert variant="destructive" className="py-2">
          <AlertDescription className="flex w-full items-center gap-2">
            <span className="pl-2">
              {value.length} {value.length > 1 ? "errors" : "error"} occurred
            </span>
            <Button
              type="button"
              variant="secondary"
              size="sm"
              className="ml-auto h-7"
              onClick={() => setShowErrors(true)}
            >
              show
            </Button>
          </AlertDescription>
        </Alert>
      ) : (
        <div className="space-y-2">
          {value.map((result) => (
            <Alert key={result.query} variant="destructive" className="py-2">
              <AlertDescription>
                <Badge variant="outline" className="mr-2 font-mono">
                  Q: {result.query}
                </Badge>
                error: {String(result.error)}
              </AlertDescription>
            </Alert>
          ))}
        </div>
      )}
    </div>
  );
}

export function GeoMultiError({ value, className }: { value: BulkError[]; className?: string }) {
  const resetKey = value.map((e) => `${e.query}:${String(e.error)}`).join("\0");
  return <GeoMultiErrorInner key={resetKey} value={value} className={className} />;
}
