'use client';

import Button from '@/components/button.component';
import React, { useState } from 'react'
import toast from 'react-hot-toast';

interface FormData {
  username?: string | undefined;
  email?: string | undefined;
  password?: string | undefined;
}

const initialFormData = {
  username: "",
  email: "",
  password: "",
}

const Auth = () => {
  const [formData, setFormData] = useState<FormData>(initialFormData)

  function handleInput(e: any) {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: value });
  }

  function handleSignUp(e: any) {
    e.preventDefault();
    fetch(`http://localhost:9001/registerUser`, {
      method: 'POST',
      headers: {
        'Content-type': 'application/json',
      },
      body: JSON.stringify(formData),
    })
      .then((res) => res.json())
      .then((result) => {
        console.log(result)
        setFormData(initialFormData)
        toast.success('Task created successfully', {
          position: 'bottom-right',
        });
      })
  }

  return (
    <div className='flex items-center justify-center min-h-screen w-full'>
      <div className='p-8 flex flex-col gap-8 items-center justify-center border border-neutral-500 backdrop-blur-sm rounded-md'>
        <h2 className='text-4xl font-bold'>Register</h2>
        <form className='flex flex-col justify-center gap-4' onSubmit={handleSignUp}>
          <div className='flex flex-col gap-2 text-white'>
            <input
              className='w-[400px] bg-transparent text-white border border-solid border-neutral-500 px-4 py-1 rounded-md outline-transparent active:outline-0'
              placeholder='Enter username'
              type='text'
              name='username'
              value={formData?.username}
              onChange={handleInput}
            />
            <input
              className='w-[400px] bg-transparent text-white border border-solid border-neutral-500 px-4 py-1 rounded-md outline-transparent active:outline-0'
              placeholder='Enter email'
              type='email'
              name='email'
              value={formData?.email}
              onChange={handleInput}
            />
            <input
              className='w-[400px] bg-transparent text-white border border-solid border-neutral-500 px-4 py-1 rounded-md outline-transparent active:outline-0'
              placeholder='Enter password'
              type='password'
              name='password'
              value={formData?.password}
              onChange={handleInput}
            />
          </div>
          <Button type='submit' shape='rect' value='Register' extraClassProps='w-max' />
        </form>
      </div>
    </div>
  )
}

export default Auth
