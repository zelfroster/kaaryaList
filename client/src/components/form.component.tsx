import {
  ChangeEvent,
  Dispatch,
  FormEvent,
  SetStateAction,
  useState,
} from 'react';

import toast from 'react-hot-toast';

import { Task } from '@/app/page';

import Button from './button.component';

type FormPropTypes = {
  taskId?: number;
  task?: Task;
  tasks: Task[];
  setTasks: Dispatch<SetStateAction<Task[]>>;
  closeModal: () => void;
};

export default function Form(formProps: FormPropTypes) {
  const { taskId, task, tasks, setTasks, closeModal } = formProps;
  const [curTaskName, setCurTaskName] = useState(task?.name || '');

  function handleChange(event: ChangeEvent<HTMLInputElement>) {
    setCurTaskName(event.target.value);
  }

  function handleSubmit(event: FormEvent) {
    event.preventDefault();

    taskId ? updateTask() : createTask();
    closeModal();
  }

  function createTask() {
    const taskToCreate = { name: curTaskName };
    fetch(`http://localhost:9000/createTask`, {
      method: 'POST',
      headers: {
        'Content-type': 'application/json',
      },
      body: JSON.stringify(taskToCreate),
    })
      .then((res) => res.json())
      .then((newTask) => {
        setTasks(() => [...tasks, newTask]);
        toast.success('Task created successfully', {
          position: 'bottom-right',
        });
      });
  }

  function updateTask() {
    const taskToUpdate = { name: curTaskName, isComplete: task?.isComplete };
    fetch(`http://localhost:9000/updateTask/${task?.id}`, {
      method: 'PUT',
      headers: {
        'Content-type': 'application/json',
      },
      body: JSON.stringify(taskToUpdate),
    })
      .then((res) => res.json())
      .then((updatedTask) => {
        toast.success('Task updated successfully', {
          position: 'bottom-right',
        });
        setTasks(() =>
          tasks.map((task) => {
            return taskId === task.id ? updatedTask : task;
          }),
        );
      });
  }

  return (
    <form
      onSubmit={handleSubmit}
      className='flex flex-col gap-4 rounded-md border border-white/90 bg-black p-6'
    >
      <label htmlFor='taskName' className='text-xl font-bold'>
        Enter Task Name
      </label>
      <input
        type='text'
        id='taskName'
        value={curTaskName}
        className='rounded border border-white/60 bg-black/20 px-2 py-1 outline-none focus:border-white'
        onChange={handleChange}
      />
      <Button
        value={taskId ? 'Confirm Edit' : 'Create Task'}
        extraClassProps='justify-center'
      />
    </form>
  );
}
