'use client';

import Card from '@/components/card.component';
import { useEffect, useState } from 'react';

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
      <header className='container mx-auto flex h-20 justify-between border-b-white/10 p-4'>
        <p className='text-4xl font-bold text-white'>कार्यList</p>
        <button className='h-fit rounded-[4px] border-2 border-white/60 bg-white/90 px-4 py-1 text-zinc-950'>
          Login
        </button>
      </header>
      <main className='container mx-auto mb-auto p-4'>
        <div className='flex justify-between'>
          <h2 className='mb-4 text-2xl font-bold'>Tasks</h2>
          <button className='p-4'>Add Task </button>
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
      <footer className='p-4 text-center'>Made by zelfroster</footer>
    </>
  );
}
