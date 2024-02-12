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
    <div className='group relative flex gap-2 rounded-[4px] border border-zinc-700 bg-neutral-900/90 py-4 pl-6 pr-10'>
      <label htmlFor='isComplete' onClick={updateTask}>
        {isComplete ? <CircleTick /> : <Circle />}
      </label>
      <input
        className='hidden'
        id='isComplete'
        type='checkbox'
        checked={isComplete}
        readOnly
      />
      <p className={isComplete ? 'line-through' : ''}>{name}</p>
      <span className='absolute -right-2 -top-2 hidden rounded-full bg-[#ddd] px-1 group-hover:flex'>
        <Button shape='round' icon={<EditIcon />} onClick={openModal} />
        <Button
          shape='round'
          icon={<DeleteIcon />}
          onClick={deleteTask}
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
