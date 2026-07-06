import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../../contexts/AuthContext';
import './HomePage.css';

export default function HomePage() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/', { replace: true });
  };

  return (
    <div className="home-page">
      <div className="home-page__card">
        <h1 className="home-page__title">Bem-vindo, {user?.name}!</h1>
        <p className="home-page__email">{user?.email}</p>
        <button className="home-page__logout" onClick={handleLogout}>
          Sair
        </button>
      </div>
    </div>
  );
}
