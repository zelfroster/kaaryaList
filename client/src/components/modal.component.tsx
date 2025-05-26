import Button from './button.component';
import CrossIcon from '../assets/cross.svg';

type ModalPropTypes = {
  closeModal: () => void;
  children: React.ReactNode;
};

export default function Modal({ closeModal, children }: ModalPropTypes) {
  return (
    // Outer div for overlay
    <div className='fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4'> {/* Changed to fixed, added z-index, centering, padding for small screens */}
      {/* Modal content container - relative for positioning close button */}
      <div className='relative bg-black rounded-md shadow-xl w-full max-w-lg'> {/* Added bg, rounded, shadow, width constraints */}
        {/* The children (e.g. Form component) will provide its own padding and border */}
        {children}
        {/* Close button positioned relative to this container */}
        {/* The Form has p-6, so top-4 right-4 relative to the Form's boundary (created by this div) is good */}
        <div className='absolute top-4 right-4'> 
          <Button
            variant='icon' // Use new variant
            icon={<CrossIcon />}
            onClick={closeModal}
            extraClassProps="rounded-full" // Keep it round
          />
        </div>
      </div>
    </div>
  );
}
