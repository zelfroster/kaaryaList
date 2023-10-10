import Button from './button.component';
import CrossIcon from '../assets/cross.svg';

export default function Modal({
  toggleModal,
  children,
}: {
  toggleModal: () => void;
  children: React.ReactNode;
}) {
  return (
    <div className='absolute min-h-screen w-full bg-black/40 backdrop-blur-[2px]'>
      <div className='absolute inset-0 m-auto flex h-fit w-max flex-col'>
        {children}
        <div className='absolute -right-2 -top-2'>
          <Button
            type='round'
            icon={<CrossIcon />}
            onClickHandler={toggleModal}
          />
        </div>
      </div>
    </div>
  );
}
