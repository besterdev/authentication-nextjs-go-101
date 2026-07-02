"use client";

import { createContext, useContext } from "react";

const AuthLoadingContext = createContext(false);

export const AuthLoadingProvider = ({
  isLoading,
  children,
}: {
  isLoading: boolean;
  children: React.ReactNode;
}) => (
  <AuthLoadingContext.Provider value={isLoading}>
    {children}
  </AuthLoadingContext.Provider>
);

export const useAuthLoading = () => useContext(AuthLoadingContext);
