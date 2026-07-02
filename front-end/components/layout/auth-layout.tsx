import { ShieldCheck } from "lucide-react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

export const AuthLayout = ({
  title,
  description,
  children,
}: {
  title: string;
  description: string;
  children: React.ReactNode;
}) => (
  <div className="auth-surface relative min-h-screen overflow-hidden">
    <div className="pointer-events-none absolute inset-0 bg-[linear-gradient(to_right,oklch(0.92_0.01_240/0.35)_1px,transparent_1px),linear-gradient(to_bottom,oklch(0.92_0.01_240/0.35)_1px,transparent_1px)] bg-[size:3rem_3rem] [mask-image:radial-gradient(ellipse_at_center,black,transparent_78%)]" />

    <div className="relative mx-auto flex min-h-screen w-full max-w-6xl items-center justify-center px-4 py-10 sm:px-6">
      <div className="grid w-full items-center gap-10 lg:grid-cols-[1fr_26rem] lg:gap-16">
        <div className="hidden space-y-6 lg:block">
          <div className="inline-flex items-center gap-2 rounded-full border border-primary/15 bg-white/80 px-3 py-1 text-xs font-medium text-primary shadow-sm backdrop-blur">
            <ShieldCheck className="size-3.5" />
            Secure authentication
          </div>
          <div className="space-y-3">
            <h1 className="text-4xl font-semibold tracking-tight text-foreground">
              Welcome back to your workspace
            </h1>
            <p className="max-w-md text-base leading-relaxed text-muted-foreground">
              Sign in to manage your account, review session details, and keep
              your access tokens up to date.
            </p>
          </div>
          <div className="grid max-w-md gap-3 sm:grid-cols-2">
            <div className="rounded-2xl border border-white/70 bg-white/70 p-4 shadow-sm backdrop-blur">
              <p className="text-sm font-medium text-foreground">
                JWT session
              </p>
              <p className="mt-1 text-xs text-muted-foreground">
                Auto refresh before expiry
              </p>
            </div>
            <div className="rounded-2xl border border-white/70 bg-white/70 p-4 shadow-sm backdrop-blur">
              <p className="text-sm font-medium text-foreground">
                Protected routes
              </p>
              <p className="mt-1 text-xs text-muted-foreground">
                Dashboard guarded client-side
              </p>
            </div>
          </div>
        </div>

        <Card className="w-full border-white/80 bg-white/90 shadow-xl shadow-primary/5 backdrop-blur-xl">
          <CardHeader className="space-y-2 pb-2">
            <div className="mb-1 flex size-10 items-center justify-center rounded-xl bg-primary/10 text-primary lg:hidden">
              <ShieldCheck className="size-5" />
            </div>
            <CardTitle className="text-2xl font-semibold tracking-tight">
              {title}
            </CardTitle>
            <CardDescription className="text-sm leading-relaxed">
              {description}
            </CardDescription>
          </CardHeader>
          <CardContent>{children}</CardContent>
        </Card>
      </div>
    </div>
  </div>
);
