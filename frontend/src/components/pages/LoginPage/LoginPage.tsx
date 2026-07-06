import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import AuthLayout from '../../templates/AuthLayout/AuthLayout';
import LoginForm from '../../organisms/LoginForm/LoginForm';
import { useAuth } from '../../../contexts/AuthContext';
import './LoginPage.css';

export default function LoginPage() {
  const navigate = useNavigate();
  const { login, isAuthenticated } = useAuth();
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  if (isAuthenticated) {
    navigate('/home', { replace: true });
    return null;
  }

  const handleLogin = async (data: { email: string; password: string }) => {
    setError('');
    setLoading(true);
    try {
      await login(data.email, data.password);
      navigate('/home', { replace: true });
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Erro ao fazer login');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="login-page">
      <AuthLayout bannerSrc="/banner-login.png" bannerAlt="Code Connect - Login">
        <LoginForm onSubmit={handleLogin} error={error} loading={loading} />
      </AuthLayout>
    </div>
  );
}
