import { AuthLayout } from "@/components/layout/auth-layout";
import { RegisterForm } from "@/components/auth/register-form";

const RegisterPage = () => (
  <AuthLayout
    title="Create an account"
    description="Register with your email and password"
  >
    <RegisterForm />
  </AuthLayout>
);

export default RegisterPage;
