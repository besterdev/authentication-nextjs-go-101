"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { DashboardLoadingSkeleton } from "@/components/layout/dashboard-loading-skeleton";
import { useAuthStore } from "@/stores/auth-store";

export const AuthGuard = ({ children }: { children: React.ReactNode }) => {
  const isLoading = useAuthStore((state) => state.isLoading);
  const isAuthenticated = useAuthStore((state) => state.user !== null);
  const router = useRouter();

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.replace("/login");
    }
  }, [isAuthenticated, isLoading, router]);

  if (isLoading) {
    return <DashboardLoadingSkeleton />;
  }

  if (!isAuthenticated) {
    return null;
  }

  return <>{children}</>;
};
