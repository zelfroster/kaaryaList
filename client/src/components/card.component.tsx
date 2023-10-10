import { Dispatch, SetStateAction, useState } from 'react';
import { createPortal } from 'react-dom';

import { Task } from '@/app/page';

import Button from './button.component';

import Circle from '../assets/tick-circle.svg';
import CircleTick from '../assets/tick-circle-broken.svg';
import EditIcon from '../assets/pencil.svg';
import Form from './form.component';
import Modal from './modal.component';

type CardPropsType = {
  task: Task;
  tasks: Task[];
  setTasks: Dispatch<SetStateAction<Task[]>>;
};

export default function Card(cardProps: CardPropsType) {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const { task, tasks, setTasks } = cardProps;
  const { name, isComplete } = task;

  function editTask() {
    setIsModalOpen(true);
  }

  return (
    <div className='group relative flex gap-2 rounded-[4px] border border-zinc-700 bg-zinc-50/10 py-4 pl-6 pr-10'>
      <label
        htmlFor='isComplete'
        onClick={() =>
          setTasks(() => {
            const newTask: Task[] = [];
            tasks.forEach((taskk) =>
              taskk.id === task.id
                ? newTask.push({ ...taskk, ['isComplete']: !isComplete })
                : newTask.push(taskk),
            );
            return newTask;
          })
        }
      >
        {isComplete ? <CircleTick /> : <Circle />}
      </label>
      <input
        className='hidden'
        id='isComplete'
        type='checkbox'
        checked={isComplete}
        onChange={() => null}
      />
      <p>{name}</p>
      <Button
        type='round'
        icon={<EditIcon />}
        extraClassProps='absolute -right-2 -top-2 hidden group-hover:block'
        onClickHandler={editTask}
      />
      {isModalOpen &&
        createPortal(
          <Modal toggleModal={() => setIsModalOpen(false)}>
            <Form task={task} tasks={tasks} setTasks={setTasks} />
          </Modal>,
          document.body,
        )}
    </div>
  );
}
