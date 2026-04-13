'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import { useRouter, useSearchParams } from 'next/navigation';
import { Button } from '@/components/ui/Button';
import api from '@/lib/api';
import axios from 'axios';

export default function LoginPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [signupSuccess, setSignupSuccess] = useState(false);

  useEffect(() => {
    if (searchParams.get('signup') === 'success') {
      setSignupSuccess(true);
    }
  }, [searchParams]);

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    setSignupSuccess(false);

    try {
      const response = await api.post('/auth/login', {
        email,
        password,
      });

      const { token } = response.data.data;
      localStorage.setItem('velox_token', token);
      
      router.push('/dashboard');
    } catch (err: unknown) {
      if (axios.isAxiosError(err) && err.response) {
        setError(err.response.data.error || 'Login failed. Please check your credentials.');
      } else {
        setError('A connection error occurred. Please check your backend.');
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-background flex flex-col justify-center py-12 sm:px-6 lg:px-8 relative overflow-hidden">
      
      {/* Background Styling */}
      <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px] bg-primary/20 blur-[150px] rounded-full pointer-events-none z-0"></div>

      <div className="sm:mx-auto sm:w-full sm:max-w-md relative z-10 flex flex-col items-center">
        <Link href="/" className="flex items-center gap-2 group mb-6">
          <div className="w-10 h-10 rounded-xl bg-primary/10 flex items-center justify-center font-bold text-primary tracking-tighter border border-primary/20 shadow-[0_0_15px_rgba(255,90,0,0.1)] transition-all">
            <svg className="w-6 h-6 text-primary transform -skew-x-12" fill="currentColor" viewBox="0 0 24 24">
              <path d="M4.5 4L10.5 19h3.5L8 4H4.5z" />
              <path d="M12.5 4L18.5 19h3.5L16 4h-3.5z" className="opacity-40" />
            </svg>
          </div>
          <span className="font-extrabold text-2xl tracking-wide text-white ml-2">
            Velox
          </span>
        </Link>
        <h2 className="mt-2 text-center text-2xl font-bold leading-9 tracking-tight text-white">
          Sign in to your account
        </h2>
        <p className="mt-2 text-center text-sm text-white/50">
          Or{' '}
          <Link href="/signup" className="font-semibold text-primary hover:text-primary/80 transition-colors">
            create a new workspace for free
          </Link>
        </p>
      </div>

      <div className="mt-10 sm:mx-auto sm:w-full sm:max-w-[400px] relative z-10">
        <div className="bg-surface/80 backdrop-blur-xl py-10 px-6 shadow-[0_20px_60px_-15px_rgba(0,0,0,1)] border border-white/10 sm:rounded-2xl sm:px-10">
          
          {/* OAuth Buttons */}
          <div>
            <button className="w-full flex items-center justify-center gap-3 bg-white/5 hover:bg-white/10 border border-white/10 text-white font-semibold py-2.5 px-4 rounded-xl transition-all">
              <svg className="h-5 w-5" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                <path fillRule="evenodd" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z" clipRule="evenodd" />
              </svg>
              Continue with GitHub
            </button>
          </div>

          <div className="relative mt-8 mb-8">
            <div className="absolute inset-0 flex items-center" aria-hidden="true">
              <div className="w-full border-t border-white/10" />
            </div>
            <div className="relative flex justify-center text-sm font-medium leading-6">
              <span className="bg-surface/80 backdrop-blur-md px-4 text-white/40">Or continue with</span>
            </div>
          </div>

          {signupSuccess && (
            <div className="mb-4 p-3 bg-green-500/10 border border-green-500/20 rounded-xl text-green-500 text-sm">
              Account created successfully! Please sign in.
            </div>
          )}

          {error && (
            <div className="mb-4 p-3 bg-red-500/10 border border-red-500/20 rounded-xl text-red-500 text-sm">
              {error}
            </div>
          )}

          <form className="space-y-6" onSubmit={handleLogin}>
            <div>
              <label htmlFor="email" className="block text-sm font-medium leading-6 text-white/90">
                Email address
              </label>
              <div className="mt-2 text-white">
                <input
                  id="email"
                  name="email"
                  type="email"
                  autoComplete="email"
                  required
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  className="block w-full rounded-xl border-0 bg-white/5 py-2.5 text-white shadow-sm ring-1 ring-inset ring-white/10 focus:ring-2 focus:ring-inset focus:ring-primary sm:text-sm sm:leading-6 px-4 transition-all"
                  placeholder="example@velox.dev"
                />
              </div>
            </div>

            <div>
              <div className="flex items-center justify-between">
                <label htmlFor="password" className="block text-sm font-medium leading-6 text-white/90">
                  Password
                </label>
                <div className="text-sm leading-6">
                  <a href="#" className="font-semibold text-primary hover:text-primary/80 transition-colors">
                    Forgot password?
                  </a>
                </div>
              </div>
              <div className="mt-2">
                <input
                  id="password"
                  name="password"
                  type="password"
                  autoComplete="current-password"
                  required
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  className="block w-full rounded-xl border-0 bg-white/5 py-2.5 text-white shadow-sm ring-1 ring-inset ring-white/10 focus:ring-2 focus:ring-inset focus:ring-primary sm:text-sm sm:leading-6 px-4 transition-all"
                />
              </div>
            </div>

            <div className="pt-2">
              <Button 
                type="submit" 
                variant="primary" 
                className="w-full justify-center text-sm h-11"
                disabled={loading}
              >
                {loading ? 'Signing in...' : 'Sign in'}
              </Button>
            </div>
          </form>
        </div>
        
        {/* Helper bottom text */}
        <p className="text-center mt-6 text-xs text-white/30 font-mono">
          Protected by Enterprise Grade SSH Security.
        </p>
      </div>
    </div>
  );
}
