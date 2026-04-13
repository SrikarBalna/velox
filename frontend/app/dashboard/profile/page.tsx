'use client';

import React, { useEffect, useState } from 'react';
import Link from 'next/link';
import api from '@/lib/api';

export default function ProfilePage() {
  const [profile, setProfile] = useState<{name: string, email: string} | null>(null);

  useEffect(() => {
    // Fetch real profile data from backend
    api.get('/dashboard')
      .then(res => {
        if (res.data) {
          setProfile({
            name: res.data.user_name,
            email: res.data.user_email,
          });
        }
      })
      .catch(err => {
        console.error("Failed to fetch dashboard data:", err);
      });
  }, []);

  return (
    <div className="flex flex-col gap-6 max-w-4xl mx-auto">
      <div className="flex items-center gap-4 mb-2">
        <Link href="/dashboard" className="text-white/50 hover:text-white transition-colors bg-white/5 hover:bg-white/10 p-2 rounded-lg">
            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
            </svg>
        </Link>
        <h1 className="text-3xl font-bold text-white relative z-10">User Identity</h1>
      </div>

      <div className="dev-card p-8 bg-surface/40 flex flex-col md:flex-row gap-10 items-center justify-center">
            {/* Avatar Definition block */}
            <div className="flex flex-col text-center">
                <div className="w-40 h-40 mx-auto rounded-full bg-white/5 flex items-center justify-center p-2 border border-white/10 shadow-inner mb-4">
                    <img src={`https://api.dicebear.com/7.x/avataaars/svg?seed=${profile?.name || 'Admin'}`} alt="Avatar" className="w-full h-full rounded-full" />
                </div>
                <h3 className="text-xl font-bold text-white tracking-tight">{profile?.name || '...'}</h3>
                <p className="text-xs text-white/40 font-mono mt-1">Globally Generated Avatar</p>
            </div>

            {/* Read-only Info Form */}
            <div className="flex flex-col flex-1 w-full gap-6">
                <div className="flex flex-col gap-2">
                    <label className="text-sm font-bold text-white/50 uppercase tracking-wider">Display Name</label>
                    <input 
                        type="text" 
                        readOnly 
                        value={profile?.name || 'Loading...'} 
                        className="w-full bg-black border border-white/10 rounded-lg px-4 py-3 text-white focus:outline-none transition-all font-mono opacity-80 cursor-default" 
                    />
                </div>
                
                <div className="flex flex-col gap-2">
                    <label className="text-sm font-bold text-white/50 uppercase tracking-wider">Email Address</label>
                    <input 
                        type="email" 
                        readOnly 
                        value={profile?.email || 'Loading...'} 
                        className="w-full bg-black border border-white/10 rounded-lg px-4 py-3 text-white focus:outline-none transition-all font-mono opacity-80 cursor-default" 
                    />
                </div>
            </div>
      </div>
    </div>
  );
}
