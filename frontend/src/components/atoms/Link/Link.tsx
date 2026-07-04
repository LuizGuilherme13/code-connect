import './Link.css';

interface LinkProps {
  href?: string;
  children: React.ReactNode;
  onClick?: (e: React.MouseEvent<HTMLAnchorElement>) => void;
}

export default function Link({ href = '#', children, onClick }: LinkProps) {
  return (
    <a className="link" href={href} onClick={onClick}>
      {children}
    </a>
  );
}
