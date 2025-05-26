import React from 'react';

interface ButtonProps {
  type?: "submit" | "reset" | "button" | undefined;
  icon?: React.ReactNode;
  value?: string;
  children?: React.ReactNode; // Added to allow children for text content
  variant?: 'primary' | 'secondary' | 'icon';
  size?: 'medium'; // For now, only 'medium', can be expanded
  onClick?: () => void;
  extraClassProps?: string;
  disabled?: boolean; // Added disabled prop
}

const Button: React.FC<ButtonProps> = ({
  type = "button",
  icon,
  value,
  children,
  variant = 'primary',
  size = 'medium',
  onClick,
  extraClassProps = '',
  disabled = false,
}) => {
  const baseStyles =
    'font-medium rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-sky-500 focus:ring-offset-black transition-all duration-150 ease-in-out flex items-center justify-center gap-2'; // Added flex, items-center, justify-center, gap-2 for icon+text alignment

  const sizeStyles = {
    medium: {
      primary: 'px-4 py-2',
      secondary: 'px-4 py-2',
      icon: 'p-2 min-w-10 min-h-10', // Added min-w/h for tap target
    },
  };

  const variantStyles = {
    primary: 'bg-sky-500 hover:bg-sky-600 text-white',
    secondary: 'bg-zinc-700 hover:bg-zinc-600 text-zinc-100',
    icon: 'bg-transparent hover:bg-zinc-700 text-zinc-400 hover:text-zinc-100',
  };
  
  const disabledStyles = 'opacity-50 cursor-not-allowed';

  const currentSizeStyles = sizeStyles[size];
  const paddingStyles = currentSizeStyles[variant as keyof typeof currentSizeStyles]; // Type assertion

  return (
    <button
      type={type}
      className={`
        ${baseStyles}
        ${paddingStyles}
        ${variantStyles[variant]}
        ${disabled ? disabledStyles : ''}
        ${extraClassProps}
      `}
      onClick={onClick}
      disabled={disabled}
    >
      {icon}
      {value}
      {children} {/* Children can be used for text, allowing more flexibility */}
    </button>
  );
};

export default Button;
