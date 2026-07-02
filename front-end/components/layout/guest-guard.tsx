"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { AuthLoadingProvider } from "@/hooks/use-auth-loading";
import { useAuthStore } from "@/stores/auth-store";

export const GuestGuard = ({ children }: { children: React.ReactNode }) => {
  const isLoading = useAuthStore((state) => state.isLoading);
  const isAuthenticated = useAuthStore((state) => state.user !== null);
  const router = useRouter();

  useEffect(() => {
    if (!isLoading && isAuthenticated) {
      router.replace("/dashboard");
    }
  }, [isAuthenticated, isLoading, router]);

  if (isAuthenticated) {
    return null;
  }

  return (
    <AuthLoadingProvider isLoading={isLoading}>{children}</AuthLoadingProvider>
  );
};
