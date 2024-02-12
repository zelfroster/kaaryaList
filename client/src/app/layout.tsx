import Script from 'next/script';
import type { Metadata } from 'next';
import './globals.css';

export const metadata: Metadata = {
  title: 'kaaryaList',
  description:
    'A Task manager application to help you be more productive.',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang='en' className='scroll-smooth'>
      <head>
        <Script
          src='https://cdn.jsdelivr.net/npm/particles.js@2.0.0/particles.min.js'
          strategy='beforeInteractive'
          defer
        />
        <Script
          id='particlesJsScript'
          dangerouslySetInnerHTML={{
            __html: `
          particlesJS.load("particles-js", "/particles.json")
        `,
          }}
          strategy='afterInteractive'
          defer
        />
      </head>
      <body className='relative flex min-h-screen w-full flex-col bg-black text-zinc-300'>
        {children}
        <div id='particles-js' className='fixed inset-0 -z-10'></div>
      </body>
    </html>
  );
}
