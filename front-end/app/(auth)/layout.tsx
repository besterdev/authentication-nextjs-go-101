import { GuestGuard } from "@/components/layout/guest-guard";

const AuthGroupLayout = ({ children }: { children: React.ReactNode }) => (
  <GuestGuard>{children}</GuestGuard>
);

export default AuthGroupLayout;
