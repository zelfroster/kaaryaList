import { Task } from '@/app/page';
import Button from './button.component';
import { ChangeEvent, Dispatch, SetStateAction, useState } from 'react';

export default function Form({
  task,
  tasks,
  setTasks,
}: {
  task: Task;
  tasks: Task[];
  setTasks: Dispatch<SetStateAction<Task[]>>;
}) {
  const [taskName, setTaskName] = useState(task.name);

  function handleChange(event: ChangeEvent<HTMLInputElement>) {
    setTaskName(event.target.value);
  }

  function updateTask(task: Task) {
    const updatedTask = fetch(`http://localhost:9000/updateTask/${task.id}`, {
      method: 'PUT',
      headers: {
        'Content-type': 'application/json',
      },
      body: JSON.stringify(task),
    }).then((res) => res.json());
    return updatedTask;
  }

  function handleSubmit(event: any) {
    event.preventDefault();

    const taskToUpdate = { ...task, ['name']: taskName };
    updateTask(taskToUpdate).then((newTask) => {
      setTasks(() => {
        const newTasks: Task[] = [];
        tasks.forEach((curTask) =>
          curTask.id === task.id
            ? newTasks.push(newTask)
            : newTasks.push(curTask),
        );
        return newTasks;
      });
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
        value={taskName}
        className='rounded border border-white/60 bg-black/20 px-2 py-1 outline-none focus:border-white'
        onChange={handleChange}
      />
      <Button value='Confirm Edit' extraClassProps='justify-center' />
    </form>
  );
}
