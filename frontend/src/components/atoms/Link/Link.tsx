import { Link as RouterLink } from 'react-router-dom';
import './Link.css';

interface LinkProps {
  href?: string;
  to?: string;
  children: React.ReactNode;
  onClick?: (e: React.MouseEvent<HTMLAnchorElement>) => void;
}

export default function Link({ href, to, children, onClick }: LinkProps) {
  if (to) {
    return (
      <RouterLink className="link" to={to} onClick={onClick}>
        {children}
      </RouterLink>
    );
  }

  return (
    <a className="link" href={href || '#'} onClick={onClick}>
      {children}
    </a>
  );
}
