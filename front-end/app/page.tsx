"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { ShieldCheck } from "lucide-react";
import { useAuth } from "@/hooks/use-auth";

const HomePage = () => {
  const { isAuthenticated, isLoading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (isLoading) return;
    router.replace(isAuthenticated ? "/dashboard" : "/login");
  }, [isAuthenticated, isLoading, router]);

  return (
    <div className="auth-surface flex min-h-screen items-center justify-center">
      <div className="flex flex-col items-center gap-4 rounded-2xl border border-white/80 bg-white/90 px-8 py-7 shadow-xl shadow-primary/5 backdrop-blur">
        <div className="flex size-12 items-center justify-center rounded-2xl bg-primary/10 text-primary">
          <ShieldCheck className="size-6 animate-pulse" />
        </div>
        <p className="text-sm font-medium text-foreground">Loading workspace</p>
        <p className="text-xs text-muted-foreground">Preparing your session...</p>
      </div>
    </div>
  );
};

export default HomePage;
