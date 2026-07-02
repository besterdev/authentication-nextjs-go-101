import { AuthLayout } from "@/components/layout/auth-layout";
import { LoginForm } from "@/components/auth/login-form";

const LoginPage = () => (
  <AuthLayout
    title="Sign in"
    description="Enter your credentials to access your account"
  >
    <LoginForm />
  </AuthLayout>
);

export default LoginPage;
