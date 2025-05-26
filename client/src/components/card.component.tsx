import { Dispatch, SetStateAction, useState } from 'react';
import { createPortal } from 'react-dom';

import { Task } from '@/app/page';

import Button from './button.component';
import Form from './form.component';
import Modal from './modal.component';

import Circle from '../assets/tick-circle.svg';
import CircleTick from '../assets/tick-circle-broken.svg';
import EditIcon from '../assets/pencil.svg';
import DeleteIcon from '../assets/dustbin.svg';
import toast from 'react-hot-toast';

type CardPropsType = {
  task: Task;
  tasks: Task[];
  setTasks: Dispatch<SetStateAction<Task[]>>;
};

export default function Card(cardProps: CardPropsType) {
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const { task, tasks, setTasks } = cardProps;
  const { name, isComplete } = task;

  function openModal() {
    setIsCreateModalOpen(true);
  }

  function closeModal() {
    setIsCreateModalOpen(false);
  }

  function deleteTask() {
    fetch(`http://localhost:9001/deleteTask/${task.id}`, {
      method: 'DELETE',
    })
      .then((res) => res.json())
      .then((res) => {
        if (res.retval) {
          const newTasks = tasks.filter((curTask) => curTask.id !== task.id);
          setTasks(newTasks);
          toast.success('Task deleted successfully', {
            position: 'bottom-right',
          });
        }
      });
  }

  function updateTask() {
    const taskToComplete = { name: name, isComplete: !isComplete };
    fetch(`http://localhost:9001/updateTask/${task.id}`, {
      method: 'PUT',
      headers: {
        'Content-type': 'application/json',
      },
      body: JSON.stringify(taskToComplete),
    })
      .then((res) => res.json())
      .then((updatedTask) => {
        setTasks(() =>
          tasks.map((curTask) => {
            return curTask.id === task.id ? updatedTask : curTask;
          }),
        );
      });
  }

  return (
    <div className='group relative flex gap-2 rounded-md border border-zinc-700 bg-neutral-900/90 py-4 pl-6 pr-10'> {/* Updated rounded-[4px] to rounded-md */}
      <label htmlFor={`isComplete-${task.id}`} onClick={updateTask} className="cursor-pointer"> {/* Unique htmlFor and cursor-pointer */}
        {isComplete ? <CircleTick /> : <Circle />}
      </label>
      <input
        className='hidden'
        id={`isComplete-${task.id}`} // Unique ID for input
        type='checkbox'
        checked={isComplete}
        readOnly
      />
      <p className={`${isComplete ? 'line-through text-zinc-500' : 'text-zinc-100'} flex-grow`}>{name}</p> {/* Added text colors and flex-grow */}
      {/* Updated span for action buttons */}
      <span className='absolute top-3 right-3 flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity duration-150 ease-in-out'>
        <Button variant='icon' icon={<EditIcon />} onClick={openModal} extraClassProps="rounded-full" /> {/* Use new variant, ensure icons are visible */}
        <Button
          variant='icon'
          icon={<DeleteIcon />}
          onClick={deleteTask}
          extraClassProps="rounded-full"
        />
      </span>
      {isCreateModalOpen &&
        createPortal(
          <Modal closeModal={closeModal}>
            <Form
              taskId={task.id}
              task={task}
              tasks={tasks}
              setTasks={setTasks}
              closeModal={closeModal}
            />
          </Modal>,
          document.body,
        )}
    </div>
  );
}
