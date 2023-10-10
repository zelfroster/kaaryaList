'use client';

import { useEffect, useState } from 'react';

import Card from '@/components/card.component';
import AddIcon from '../assets/add.svg';
import Button from '@/components/button.component';

export type Task = {
  id: number;
  name: string;
  isComplete: boolean;
};

export default function Home() {
  const [tasks, setTasks] = useState<Task[]>([]);

  useEffect(() => {
    async function getAllTasks() {
      const res = await fetch('http://localhost:9000/getTasks');
      const json = await res.json();
      setTasks(json);
    }
    getAllTasks();
  }, []);

  return (
    <>
      <header className='container mx-auto mb-6 flex h-20 justify-between border-b-white/10 px-4 py-8'>
        <p className='text-4xl font-bold text-white'>कार्यList</p>
        <Button value='Login' type='rect' />
      </header>
      <main className='container mx-auto mb-auto p-4'>
        <div className='flex justify-between'>
          <h2 className='mb-4 text-2xl font-bold'>Tasks</h2>
          <Button value='Add Task' type='rect' icon={<AddIcon />} />
        </div>
        <div className='flex flex-wrap gap-4'>
          {tasks &&
            tasks.map((task) => (
              <Card
                key={task.id}
                task={task}
                tasks={tasks}
                setTasks={setTasks}
              />
            ))}
        </div>
      </main>
      <footer className='p-4 text-center'>Made with &lt;3 by zelfroster</footer>
    </>
  );
}
