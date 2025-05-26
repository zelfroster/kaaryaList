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
    <div className='flex items-center justify-center min-h-screen w-full px-4'> {/* Added px-4 for smaller screens */}
      <div className='p-8 flex flex-col gap-8 items-center justify-center border border-zinc-700 backdrop-blur-sm rounded-md max-w-md w-full'> {/* Updated border, max-width */}
        <h2 className='text-4xl font-bold text-zinc-100'>Register</h2> {/* Ensure text color consistency */}
        <form className='flex flex-col justify-center gap-4 w-full' onSubmit={handleSignUp}> {/* Added w-full to form */}
          <div className='flex flex-col gap-2 text-white'> {/* Gap changed to 4 below for consistency */}
            <input
              className='w-full bg-zinc-800 border border-zinc-700 text-zinc-100 px-3 py-2 rounded focus:border-sky-500 focus:ring-1 focus:ring-sky-500 outline-none'
              placeholder='Enter username'
              type='text'
              name='username'
              value={formData?.username}
              onChange={handleInput}
            />
            <input
              className='w-full bg-zinc-800 border border-zinc-700 text-zinc-100 px-3 py-2 rounded focus:border-sky-500 focus:ring-1 focus:ring-sky-500 outline-none'
              placeholder='Enter email'
              type='email'
              name='email'
              value={formData?.email}
              onChange={handleInput}
            />
            <input
              className='w-full bg-zinc-800 border border-zinc-700 text-zinc-100 px-3 py-2 rounded focus:border-sky-500 focus:ring-1 focus:ring-sky-500 outline-none'
              placeholder='Enter password'
              type='password'
              name='password'
              value={formData?.password}
              onChange={handleInput}
            />
          </div>
          <Button type='submit' variant='primary' value='Register' extraClassProps='w-full' /> {/* Updated Button */}
        </form>
      </div>
    </div>
  )
}

export default Auth
