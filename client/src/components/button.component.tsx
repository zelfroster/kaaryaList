export default function Button({
  icon,
  value,
  type,
  onClickHandler,
  extraClassProps,
}: {
  icon?: React.ReactNode;
  value?: string;
  type?: 'square' | 'rect' | 'round';
  onClickHandler?: () => void;
  extraClassProps?: string;
}) {
  return (
    <button
      className={`flex h-fit gap-1 rounded-[4px] border-white/60 bg-[#dddddd] text-zinc-950 ${
        type == 'round'
          ? 'rounded-full'
          : type == 'square'
          ? 'p-1'
          : 'border-2 px-4 py-1'
      } ${extraClassProps}`}
      onClick={onClickHandler}
    >
      {value}
      {icon}
    </button>
  );
}
