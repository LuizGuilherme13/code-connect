import './Input.css';

interface InputProps {
  type?: string;
  placeholder?: string;
  value?: string;
  onChange?: (e: React.ChangeEvent<HTMLInputElement>) => void;
  name?: string;
  id?: string;
  required?: boolean;
}

export default function Input({
  type = 'text',
  placeholder,
  value,
  onChange,
  name,
  id,
  required = false,
}: InputProps) {
  return (
    <input
      className="input"
      type={type}
      placeholder={placeholder}
      value={value}
      onChange={onChange}
      name={name}
      id={id}
      required={required}
    />
  );
}
