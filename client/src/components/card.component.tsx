import { Dispatch, SetStateAction, useState } from 'react';
import { createPortal } from 'react-dom';

import { Task } from '@/app/page';

import Button from './button.component';
import Form from './form.component';
import Modal from './modal.component';

import Circle from '../assets/tick-circle.svg';
import CircleTick from '../assets/tick-circle-broken.svg';
import EditIcon from '../assets/pencil.svg';

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

  function completeTask() {
    const taskToComplete = { name: name, isComplete: !isComplete };
    console.log(taskToComplete, task.id);
    fetch(`http://localhost:9000/updateTask/${task.id}`, {
      method: 'PUT',
      headers: {
        'Content-type': 'application/json',
      },
      body: JSON.stringify(taskToComplete),
    })
      .then((res) => res.json())
      .then((updatedTask) => {
        console.log(updatedTask);
        setTasks(() =>
          tasks.map((curTask) => {
            return curTask.id === task.id ? updatedTask : curTask;
          }),
        );
        console.log(task);
      });
  }

  return (
    <div className='group relative flex gap-2 rounded-[4px] border border-zinc-700 bg-neutral-900/90 py-4 pl-6 pr-10'>
      <label htmlFor='isComplete' onClick={completeTask}>
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
      <Button
        type='round'
        icon={<EditIcon />}
        extraClassProps='absolute -right-2 -top-2 hidden group-hover:block'
        onClickHandler={openModal}
      />
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
