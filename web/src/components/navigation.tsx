import { Link, useRouterState } from "@tanstack/react-router";
import { BookOpen, Database, FolderGit2, MapPin, Search } from "lucide-react";

import { Separator } from "@/components/ui/separator";
import { cn } from "@/lib/utils";

const items = [
  { to: "/", label: "Lookup Address", icon: Search, desktopOnly: false },
  { to: "/lookup/bulk", label: "Bulk Lookup", icon: Database, desktopOnly: false },
  { to: "/lookup/docs", label: "API Documentation", icon: BookOpen, desktopOnly: true },
] as const;

export function Navigation({ slim, hideSource }: { slim?: boolean; hideSource?: boolean }) {
  const pathname = useRouterState({ select: (s) => s.location.pathname });

  return (
    <aside aria-label="Navigation">
      <div className={cn("overflow-y-auto rounded-md", !slim && "px-3 py-0 lg:py-4")}>
        <span className="flex items-center pl-2">
          <MapPin className="text-primary size-9" aria-hidden />
          <span className="ml-2 self-center whitespace-nowrap text-3xl font-bold tracking-tight text-primary">
            GeoIP
          </span>
        </span>

        <Separator className="my-2" />

        <ul className="m-0 list-none space-y-2 p-0">
          {items.map(({ to, label, icon: Icon, desktopOnly }) => {
            const active = pathname === to || (to !== "/" && pathname.startsWith(to));
            return (
              <li key={to} className={cn(desktopOnly && "hidden lg:list-item")}>
                <Link
                  to={to}
                  className={cn(
                    "route flex items-center rounded px-3 py-2 text-base font-semibold no-underline transition-all duration-250 ease-in-out",
                    active ? "bg-accent text-primary" : "text-foreground hover:bg-accent/80",
                  )}
                >
                  <Icon className="size-5 shrink-0" aria-hidden />
                  <span className="ml-3 bg-linear-to-r from-emerald-400 to-sky-400 bg-clip-text text-transparent">
                    {label}
                  </span>
                </Link>
              </li>
            );
          })}
        </ul>

        {!hideSource && (
          <>
            <Separator className="my-2 hidden lg:block" />
            <ul className="m-0 hidden list-none space-y-2 p-0 lg:block">
              <li>
                <a
                  href="https://github.com/lrstanley/geoip"
                  className="route flex items-center rounded px-3 py-2 text-base font-semibold no-underline transition-all duration-250 ease-in-out hover:bg-accent/80"
                  target="_blank"
                  rel="noreferrer"
                >
                  <FolderGit2 className="size-5 shrink-0" aria-hidden />
                  <span className="ml-3 bg-linear-to-r from-emerald-400 to-sky-400 bg-clip-text text-transparent">
                    Github Project
                  </span>
                </a>
              </li>
            </ul>
          </>
        )}
      </div>
    </aside>
  );
}
