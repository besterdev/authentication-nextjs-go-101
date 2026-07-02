export const DashboardLoadingSkeleton = () => (
  <div className="dashboard-surface min-h-screen">
    <header className="sticky top-0 z-20 border-b border-border/70 bg-white/80 backdrop-blur-xl">
      <div className="mx-auto flex h-16 max-w-5xl items-center justify-between px-4 sm:px-6">
        <div className="flex items-center gap-3">
          <div className="size-9 animate-pulse rounded-xl bg-muted" />
          <div className="space-y-1.5">
            <div className="h-3.5 w-28 animate-pulse rounded bg-muted" />
            <div className="h-3 w-20 animate-pulse rounded bg-muted" />
          </div>
        </div>
        <div className="size-9 animate-pulse rounded-full bg-muted" />
      </div>
    </header>
    <main className="mx-auto max-w-5xl space-y-8 px-4 py-8 sm:px-6 sm:py-10">
      <div className="space-y-2">
        <div className="h-4 w-16 animate-pulse rounded bg-muted" />
        <div className="h-8 w-48 animate-pulse rounded bg-muted" />
        <div className="h-4 w-72 max-w-full animate-pulse rounded bg-muted" />
      </div>
      <div className="overflow-hidden rounded-xl border border-white/80 bg-white/90 shadow-lg shadow-primary/5">
        <div className="border-b border-border/60 bg-muted/20 p-6">
          <div className="h-5 w-24 animate-pulse rounded bg-muted" />
          <div className="mt-2 h-4 w-56 animate-pulse rounded bg-muted" />
        </div>
        <div className="grid gap-4 p-6 sm:grid-cols-2">
          <div className="h-28 animate-pulse rounded-2xl bg-muted/70" />
          <div className="h-28 animate-pulse rounded-2xl bg-muted/70" />
          <div className="h-28 animate-pulse rounded-2xl bg-muted/70" />
          <div className="h-28 animate-pulse rounded-2xl bg-muted/70" />
        </div>
      </div>
    </main>
  </div>
);
