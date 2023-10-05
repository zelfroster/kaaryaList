import { Task } from '@/app/page';
import { Dispatch, SetStateAction } from 'react';
import Circle from '../assets/tick-circle.svg';
import CircleTick from '../assets/tick-circle-broken.svg';

export default function Card({
  task,
  tasks,
  setTasks,
}: {
  task: Task;
  tasks: Task[];
  setTasks: Dispatch<SetStateAction<Task[]>>;
}) {
  const { id, name, isComplete } = task;
  return (
    <div
      id={id.toString()}
      className='flex gap-4 rounded-[4px] border border-zinc-600 bg-white/10 py-6 pl-8 pr-10'
    >
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
      />
      <p className='font-bold'>{name}</p>
    </div>
  );
}
