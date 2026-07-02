import { AuthGuard } from "@/components/layout/auth-guard";
import { DashboardNav } from "@/components/layout/dashboard-nav";

const DashboardGroupLayout = ({ children }: { children: React.ReactNode }) => (
  <AuthGuard>
    <div className="dashboard-surface min-h-screen">
      <DashboardNav />
      <main className="mx-auto max-w-5xl px-4 py-8 sm:px-6 sm:py-10">
        {children}
      </main>
    </div>
  </AuthGuard>
);

export default DashboardGroupLayout;
