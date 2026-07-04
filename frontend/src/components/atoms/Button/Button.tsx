import './Button.css';

interface ButtonProps {
  type?: 'button' | 'submit' | 'reset';
  children: React.ReactNode;
  variant?: 'primary' | 'social';
  onClick?: () => void;
  disabled?: boolean;
}

export default function Button({
  type = 'button',
  children,
  variant = 'primary',
  onClick,
  disabled = false,
}: ButtonProps) {
  return (
    <button
      className={`button button--${variant}`}
      type={type}
      onClick={onClick}
      disabled={disabled}
    >
      {children}
    </button>
  );
}
