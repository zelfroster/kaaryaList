interface ButtonProps {
  type?: "submit" | "reset" | "button" | undefined;
  icon?: React.ReactNode;
  value?: string;
  shape?: 'square' | 'rect' | 'round';
  onClick?: () => void;
  extraClassProps?: string;
}

const Button: React.FC<ButtonProps> = ({
  type,
  icon,
  value,
  shape,
  onClick,
  extraClassProps,
}) => {
  return (
    <button
      type={type}
      className={`flex h-fit gap-1 rounded-[4px] border-white/60 bg-[#dddddd] text-zinc-950 ${shape == 'round'
        ? 'rounded-full'
        : shape == 'square'
          ? 'p-1'
          : 'border-2 px-4 py-1'
        } ${extraClassProps}`}
      onClick={onClick}
    >
      {value}
      {icon}
    </button>
  );
}

export default Button
