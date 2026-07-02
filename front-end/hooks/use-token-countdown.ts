"use client";

import { useEffect, useState } from "react";
import { getTokenExpiresAt } from "@/lib/auth-storage";

const formatCountdown = (totalSeconds: number) => {
  if (totalSeconds <= 0) {
    return "Refreshing...";
  }

  const minutes = Math.floor(totalSeconds / 60);
  const seconds = totalSeconds % 60;
  return `${minutes}:${String(seconds).padStart(2, "0")}`;
};

export const useTokenCountdown = () => {
  const [remainingSeconds, setRemainingSeconds] = useState<number | null>(null);

  useEffect(() => {
    const update = () => {
      const expiresAt = getTokenExpiresAt();
      if (!expiresAt) {
        setRemainingSeconds(null);
        return;
      }

      setRemainingSeconds(
        Math.max(0, Math.floor((expiresAt - Date.now()) / 1000)),
      );
    };

    update();
    const intervalId = setInterval(update, 1000);
    return () => clearInterval(intervalId);
  }, []);

  return {
    remainingSeconds,
    countdownLabel:
      remainingSeconds === null ? "—" : formatCountdown(remainingSeconds),
  };
};
