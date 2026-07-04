import Button from '../../atoms/Button/Button';
import './SocialButton.css';

interface SocialButtonProps {
  iconSrc: string;
  provider: string;
  onClick?: () => void;
}

export default function SocialButton({ iconSrc, provider, onClick }: SocialButtonProps) {
  return (
    <Button variant="social" onClick={onClick}>
      <img className="social-button__icon" src={iconSrc} alt={`${provider} logo`} />
      <span>{provider}</span>
    </Button>
  );
}
