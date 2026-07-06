import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import AuthLayout from '../../templates/AuthLayout/AuthLayout';
import RegisterForm from '../../organisms/RegisterForm/RegisterForm';
import { useAuth } from '../../../contexts/AuthContext';
import './RegisterPage.css';

export default function RegisterPage() {
  const navigate = useNavigate();
  const { register, isAuthenticated } = useAuth();
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  if (isAuthenticated) {
    navigate('/home', { replace: true });
    return null;
  }

  const handleRegister = async (data: { name: string; email: string; password: string; confirmPassword: string }) => {
    if (data.password !== data.confirmPassword) {
      setError('As senhas não conferem');
      return;
    }

    setError('');
    setLoading(true);
    try {
      await register(data.name, data.email, data.password);
      navigate('/', { replace: true });
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Erro ao cadastrar');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="register-page">
      <AuthLayout bannerSrc="/banner-cadastro.png" bannerAlt="Code Connect - Cadastro">
        <RegisterForm onSubmit={handleRegister} error={error} loading={loading} />
      </AuthLayout>
    </div>
  );
}
