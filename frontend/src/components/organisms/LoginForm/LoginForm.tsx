import { useState } from 'react';
import FormField from '../../molecules/FormField/FormField';
import Input from '../../atoms/Input/Input';
import Checkbox from '../../atoms/Checkbox/Checkbox';
import Button from '../../atoms/Button/Button';
import Link from '../../atoms/Link/Link';
import SocialButton from '../../molecules/SocialButton/SocialButton';
import './LoginForm.css';

interface LoginFormProps {
  onSubmit?: (data: { email: string; password: string; rememberMe: boolean }) => void;
}

export default function LoginForm({ onSubmit }: LoginFormProps) {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [rememberMe, setRememberMe] = useState(false);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit?.({ email, password, rememberMe });
  };

  return (
    <div className="login-form">
      <div className="login-form__header">
        <h1 className="login-form__title">Login</h1>
        <p className="login-form__subtitle">Boas-vindas! Faça seu login.</p>
      </div>

      <form className="login-form__form" onSubmit={handleSubmit}>
        <FormField label="Email ou usuário" htmlFor="email">
          <Input
            type="text"
            placeholder="usuario123"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            name="email"
            id="email"
            required
          />
        </FormField>

        <FormField label="Senha" htmlFor="password">
          <Input
            type="password"
            placeholder="••••••"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            name="password"
            id="password"
            required
          />
        </FormField>

        <div className="login-form__options">
          <Checkbox
            checked={rememberMe}
            onChange={(e) => setRememberMe(e.target.checked)}
            label="Lembrar-me"
            id="rememberMe"
          />
          <Link href="#">Esqueci a senha</Link>
        </div>

        <Button type="submit">Login →</Button>
      </form>

      <div className="login-form__divider">
        <span>ou entre com outras contas</span>
      </div>

      <div className="login-form__social">
        <SocialButton iconSrc="/github.png" provider="Github" />
        <SocialButton iconSrc="/gmail.png" provider="Gmail" />
      </div>

      <div className="login-form__footer">
        <span>Ainda não tem conta?</span>
        <Link href="#">Crie seu cadastro!</Link>
      </div>
    </div>
  );
}
