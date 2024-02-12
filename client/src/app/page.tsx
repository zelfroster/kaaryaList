'use client';

import { useEffect, useState } from 'react';
import { createPortal } from 'react-dom';
import { Toaster } from 'react-hot-toast';

import Card from '@/components/card.component';
import Button from '@/components/button.component';
import Modal from '@/components/modal.component';
import Form from '@/components/form.component';

import AddIcon from '../assets/add.svg';
import { useRouter } from 'next/navigation';

export type Task = {
  id: number;
  name: string;
  isComplete: boolean;
};

export default function Home() {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);

  useEffect(() => {
    async function getAllTasks() {
      const res = await fetch('http://localhost:9001/getTasks');
      const json = await res.json();
      setTasks(json);
    }
    getAllTasks();
  }, []);

  function openModal() {
    setIsEditModalOpen(true);
  }

  function closeModal() {
    setIsEditModalOpen(false);
  }

  console.log(tasks);

  const router = useRouter()

  return (
    <>
      <header className='container mx-auto mb-6 flex h-20 justify-between border-b border-b-neutral-900 px-4 pb-16 pt-8'>
        <p className='text-4xl font-bold text-white'>kaaryaList</p>
        <Button onClick={() => router.push('/auth')} value='Login' />
      </header>
      <main className='container mx-auto mb-auto p-4'>
        <div className='flex justify-between'>
          <h2 className='mb-4 text-2xl font-bold'>Tasks</h2>
          <Button
            value='Add Task'
            icon={<AddIcon />}
            onClick={openModal}
          />
          {isEditModalOpen &&
            createPortal(
              <Modal closeModal={closeModal}>
                <Form
                  tasks={tasks}
                  setTasks={setTasks}
                  closeModal={closeModal}
                />
              </Modal>,
              document.body,
            )}
        </div>
        {tasks && tasks.length > 0
          ? <div className='flex flex-wrap gap-4'>
            {tasks.map((task) => (
              <Card
                key={task.id}
                task={task}
                tasks={tasks}
                setTasks={setTasks}
              />
            ))}
          </div>
          : <div className='flex flex-col justify-center items-center p-12 text-neutral-500'>
            <p>You have no tasks yet. </p>
            <p>Create a task to show here.</p>
          </div>}
        <Toaster />
      </main>
      <footer className='p-4 text-center'>Made with &lt;3 by zelfroster</footer>
    </>
  );
}
