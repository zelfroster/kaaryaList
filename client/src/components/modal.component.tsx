import Button from './button.component';
import CrossIcon from '../assets/cross.svg';

type ModalPropTypes = {
  closeModal: () => void;
  children: React.ReactNode;
};

export default function Modal({ closeModal, children }: ModalPropTypes) {
  return (
    <div className='absolute min-h-screen w-full bg-black/40 backdrop-blur-[2px]'>
      <div className='absolute inset-0 m-auto flex h-fit w-max flex-col'>
        {children}
        <div className='absolute -right-2 -top-2'>
          <Button
            shape='round'
            icon={<CrossIcon />}
            onClick={closeModal}
          />
        </div>
      </div>
    </div>
  );
}
