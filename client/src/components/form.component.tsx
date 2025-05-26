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
    fetch(`http://localhost:9001/createTask`, {
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
    fetch(`http://localhost:9001/updateTask/${task?.id}`, {
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
      className='flex flex-col gap-4 rounded-md border border-zinc-700 bg-black p-6' // Updated border
    >
      <label htmlFor='taskName' className='text-base font-medium text-zinc-300 mb-1'> {/* Updated label */}
        Task Name
      </label>
      <input
        type='text'
        id='taskName'
        value={curTaskName}
        className='w-full bg-zinc-800 border border-zinc-700 text-zinc-100 px-3 py-2 rounded focus:border-sky-500 focus:ring-1 focus:ring-sky-500 outline-none' // Updated input
        onChange={handleChange}
        placeholder='Enter task name...' // Added placeholder
      />
      <Button
        type='submit' // Explicitly set type
        variant='primary' // Use new variant
        value={taskId ? 'Confirm Edit' : 'Create Task'}
        extraClassProps='justify-center w-full mt-2' // Added w-full and margin
      />
    </form>
  );
}
