import AuthLayout from '../../templates/AuthLayout/AuthLayout';
import LoginForm from '../../organisms/LoginForm/LoginForm';
import './LoginPage.css';

export default function LoginPage() {
  const handleLogin = (data: { email: string; password: string; rememberMe: boolean }) => {
    console.log('Login:', data);
  };

  return (
    <div className="login-page">
      <AuthLayout bannerSrc="/banner-login.png" bannerAlt="Code Connect - Login">
        <LoginForm onSubmit={handleLogin} />
      </AuthLayout>
    </div>
  );
}
