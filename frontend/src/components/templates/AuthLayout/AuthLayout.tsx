import Banner from '../../organisms/Banner/Banner';
import './AuthLayout.css';

interface AuthLayoutProps {
  bannerSrc: string;
  bannerAlt: string;
  children: React.ReactNode;
}

export default function AuthLayout({ bannerSrc, bannerAlt, children }: AuthLayoutProps) {
  return (
    <div className="auth-layout">
      <main className="auth-layout__card">
        <Banner src={bannerSrc} alt={bannerAlt} />
        {children}
      </main>
    </div>
  );
}
