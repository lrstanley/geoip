import { createRootRoute, Link, Outlet } from "@tanstack/react-router";
import { Suspense } from "react";
import { Button } from "@/components/ui/button";
import { Toaster } from "@/components/ui/sonner";
import { TooltipProvider } from "@/components/ui/tooltip";

export const Route = createRootRoute({
  component: RootLayout,
  notFoundComponent: NotFoundPage,
});

function RootLayout() {
  return (
    <TooltipProvider delayDuration={200}>
      <Suspense
        fallback={
          <div className="flex min-h-[40vh] items-center justify-center p-20 text-muted-foreground">loading...</div>
        }
      >
        <Outlet />
      </Suspense>
      <Toaster richColors position="top-center" />
    </TooltipProvider>
  );
}

function NotFoundPage() {
  return (
    <div className="p-10 text-center">
      <h1 className="mb-2 text-2xl font-semibold">page not found</h1>
      <p className="text-muted-foreground mb-6">you know life is always ridiculous</p>
      <div className="flex justify-center gap-2">
        <Button type="button" variant="outline" onClick={() => window.history.back()}>
          go back
        </Button>
        <Button type="button" asChild>
          <Link to="/" search={{ q: undefined }}>
            lookup
          </Link>
        </Button>
      </div>
    </div>
  );
}
