"use client";

import { Clock3, Mail, Shield, UserRound } from "lucide-react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { useAuth } from "@/hooks/use-auth";
import { useTokenCountdown } from "@/hooks/use-token-countdown";

const formatDate = (value: string) => new Date(value).toLocaleString();

const profileFields = (
  user: { id: string; email: string; created_at: string },
  countdownLabel: string,
) => [
  {
    label: "User ID",
    value: user.id,
    icon: UserRound,
    mono: true,
  },
  {
    label: "Email",
    value: user.email,
    icon: Mail,
    mono: false,
  },
  {
    label: "Member since",
    value: formatDate(user.created_at),
    icon: Shield,
    mono: false,
  },
  {
    label: "Token expires in",
    value: countdownLabel,
    icon: Clock3,
    mono: true,
  },
];

const DashboardPage = () => {
  const { user } = useAuth();
  const { countdownLabel } = useTokenCountdown();

  if (!user) {
    return null;
  }

  return (
    <div className="space-y-8">
      <div className="space-y-2">
        <p className="text-sm font-medium text-primary">Overview</p>
        <h1 className="text-3xl font-semibold tracking-tight text-foreground">
          Welcome back
        </h1>
        <p className="max-w-2xl text-sm leading-relaxed text-muted-foreground">
          Your session is active. Review your account details and monitor access
          token expiry below.
        </p>
      </div>

      <Card className="overflow-hidden border-white/80 bg-white/90 shadow-lg shadow-primary/5">
        <CardHeader className="border-b border-border/60 bg-muted/20">
          <CardTitle>Profile</CardTitle>
          <CardDescription>
            Account details from the authentication API
          </CardDescription>
        </CardHeader>
        <CardContent className="grid gap-4 p-6 sm:grid-cols-2">
          {profileFields(user, countdownLabel).map(
            ({ label, value, icon: Icon, mono }) => (
              <div
                key={label}
                className="rounded-2xl border border-border/70 bg-background/80 p-4 shadow-sm"
              >
                <div className="mb-3 flex items-center gap-2 text-muted-foreground">
                  <div className="flex size-8 items-center justify-center rounded-lg bg-primary/10 text-primary">
                    <Icon className="size-4" />
                  </div>
                  <span className="text-xs font-medium uppercase tracking-wide">
                    {label}
                  </span>
                </div>
                <p
                  className={
                    mono
                      ? "font-mono text-sm text-foreground"
                      : "text-sm font-medium text-foreground"
                  }
                >
                  {value}
                </p>
              </div>
            ),
          )}
        </CardContent>
      </Card>
    </div>
  );
};

export default DashboardPage;
