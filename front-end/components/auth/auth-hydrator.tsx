"use client";

import { useEffect } from "react";
import { useAuthStore } from "@/stores/auth-store";

export const AuthHydrator = ({ children }: { children: React.ReactNode }) => {
  const hydrate = useAuthStore((state) => state.hydrate);
  const cancelHydrate = useAuthStore((state) => state.cancelHydrate);

  useEffect(() => {
    void hydrate();
    return () => cancelHydrate();
  }, [hydrate, cancelHydrate]);

  return <>{children}</>;
};
