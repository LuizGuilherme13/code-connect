import './Checkbox.css';

interface CheckboxProps {
  checked: boolean;
  onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  label: string;
  id: string;
}

export default function Checkbox({ checked, onChange, label, id }: CheckboxProps) {
  return (
    <label className="checkbox" htmlFor={id}>
      <input
        className="checkbox__input"
        type="checkbox"
        checked={checked}
        onChange={onChange}
        id={id}
      />
      <span className="checkbox__custom" />
      <span className="checkbox__label">{label}</span>
    </label>
  );
}
