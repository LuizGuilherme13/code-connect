import { useState } from 'react';
import FormField from '../../molecules/FormField/FormField';
import Input from '../../atoms/Input/Input';
import Button from '../../atoms/Button/Button';
import Link from '../../atoms/Link/Link';
import SocialButton from '../../molecules/SocialButton/SocialButton';
import './RegisterForm.css';

interface RegisterFormProps {
  onSubmit?: (data: { name: string; email: string; password: string; confirmPassword: string }) => void;
  error?: string;
  loading?: boolean;
}

export default function RegisterForm({ onSubmit, error, loading }: RegisterFormProps) {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit?.({ name, email, password, confirmPassword });
  };

  return (
    <div className="register-form">
      <div className="register-form__header">
        <h1 className="register-form__title">Cadastro</h1>
        <p className="register-form__subtitle">Olá! Preencha seus dados.</p>
      </div>

      <form className="register-form__form" onSubmit={handleSubmit}>
        <FormField label="Nome" htmlFor="name">
          <Input type="text" placeholder="Seu nome"
            value={name} onChange={(e) => setName(e.target.value)}
            name="name" id="name" required />
        </FormField>

        <FormField label="Email" htmlFor="email">
          <Input type="email" placeholder="seu@email.com"
            value={email} onChange={(e) => setEmail(e.target.value)}
            name="email" id="email" required />
        </FormField>

        <FormField label="Senha" htmlFor="password">
          <Input type="password" placeholder="••••••"
            value={password} onChange={(e) => setPassword(e.target.value)}
            name="password" id="password" required />
        </FormField>

        <FormField label="Confirmar Senha" htmlFor="confirmPassword">
          <Input type="password" placeholder="••••••"
            value={confirmPassword} onChange={(e) => setConfirmPassword(e.target.value)}
            name="confirmPassword" id="confirmPassword" required />
        </FormField>

        {error && <p className="form-error">{error}</p>}

        <Button type="submit" disabled={loading}>Cadastrar →</Button>
      </form>

      <div className="register-form__divider">
        <span>ou entre com outras contas</span>
      </div>

      <div className="register-form__social">
        <SocialButton iconSrc="/github.png" provider="Github" />
        <SocialButton iconSrc="/gmail.png" provider="Gmail" />
      </div>

      <div className="register-form__footer">
        <span>Já tem conta?</span>
        <Link to="/">Faça seu login!</Link>
      </div>
    </div>
  );
}
