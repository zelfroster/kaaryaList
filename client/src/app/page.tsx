'use client';

import { useEffect, useState } from 'react';

import Card from '@/components/card.component';
import Button from '@/components/button.component';

import AddIcon from '../assets/add.svg';
import { createPortal } from 'react-dom';
import Modal from '@/components/modal.component';
import Form from '@/components/form.component';

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
      const res = await fetch('http://localhost:9000/getTasks');
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

  return (
    <>
      <header className='container mx-auto mb-6 flex h-20 justify-between border-b border-b-neutral-900 px-4 pb-16 pt-8'>
        <p className='text-4xl font-bold text-white'>कार्यList</p>
        <Button value='Login' />
      </header>
      <main className='container mx-auto mb-auto p-4'>
        <div className='flex justify-between'>
          <h2 className='mb-4 text-2xl font-bold'>Tasks</h2>
          <Button
            value='Add Task'
            icon={<AddIcon />}
            onClickHandler={openModal}
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
