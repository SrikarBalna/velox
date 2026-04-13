'use client';

import React, { useEffect, useState } from 'react';
import Link from 'next/link';
import api from '@/lib/api';
import APIStatsCard from '@/components/dashboard/APIStatsCard';

interface APIKey {
  id: string;
  project: string;
  preview: string;
  createdAt: string;
}

export default function DashboardOverview() {
    const [profile, setProfile] = useState<{name: string, email: string} | null>(null);
    const [keys, setKeys] = useState<APIKey[]>([]);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        let isMounted = true;
        
        const fetchDashboardData = async () => {
            setIsLoading(true);
            try {
                // 1. Fetch Profile
                const profileRes = await api.get('/dashboard');
                if (isMounted && profileRes.data) {
                    setProfile({
                        name: profileRes.data.user_name || 'Admin',
                        email: profileRes.data.user_email || 'admin@velox.dev'
                    });
                }

                // 2. Fetch Keys
                const keysRes = await api.get('/auth/api-keys');
                if (isMounted && keysRes.data) {
                    const formattedKeys = keysRes.data.map((k: any) => ({
                        id: k.id,
                        project: k.name || 'Unnamed Project',
                        preview: k.display_hint || 'vlx_...',
                        createdAt: new Date(k.created_at).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
                    }));
                    setKeys(formattedKeys);
                }
            } catch (err) {
                console.error("Dashboard Data Fetch Error:", err);
            } finally {
                if (isMounted) setIsLoading(false);
            }
        };

        fetchDashboardData();

        return () => { isMounted = false; };
    }, []);

    return (
        <div className="flex flex-col gap-6 max-w-7xl mx-auto w-full pb-10">
            {/* Header Area */}
            <div className="flex flex-col md:flex-row md:items-end justify-between gap-6 mb-4">
                <div>
                    <h1 className="text-3xl font-bold text-white tracking-tight flex items-center gap-3">
                        <span className="opacity-50">Welcome back,</span> {profile?.name || '...'}
                    </h1>
                    <p className="text-foreground/50 text-sm mt-2">Here is the telemetry snapshot of your active API routes.</p>
                </div>
                
                <div className="flex gap-3">
                     <Link href="/dashboard/profile" className="px-5 py-2.5 rounded-lg bg-white/5 hover:bg-white/10 border border-white/10 text-white font-bold text-sm transition-all text-center flex items-center gap-2">
                        <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                           <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                           <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                        </svg>
                        Settings
                    </Link>
                    <Link href="/dashboard/api-keys" className="px-5 py-2.5 rounded-lg bg-primary hover:bg-primary-hover text-black font-bold text-sm transition-all shadow-[0_0_15px_rgba(255,90,0,0.3)] text-center flex items-center gap-2">
                        <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                           <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="3" d="M12 4v16m8-8H4" />
                        </svg>
                        API Key
                    </Link>
                </div>
            </div>

            {/* Content Area */}
            {isLoading ? (
                <div className="flex-1 flex flex-col items-center justify-center min-h-[400px]">
                    <div className="w-8 h-8 border-2 border-white/20 border-t-primary rounded-full animate-spin mb-4"></div>
                </div>
            ) : keys.length === 0 ? (
                <div className="dev-card p-12 bg-surface/40 flex flex-col items-center justify-center text-center border border-white/5 min-h-[400px]">
                     <div className="w-16 h-16 rounded-full bg-primary/10 border border-primary/20 flex items-center justify-center mb-6">
                        <svg className="w-8 h-8 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                           <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
                        </svg>
                    </div>
                    <h2 className="text-2xl font-bold text-white mb-2">No API Routes Configured</h2>
                    <p className="text-foreground/50 max-w-md mx-auto mb-8">Generate an API key to securely authenticate external submission nodes with Velox's execution engine.</p>
                    <Link href="/dashboard/api-keys" className="px-6 py-3 rounded-lg bg-white/10 hover:bg-white/20 border border-white/10 text-white font-bold transition-colors">
                        Configure First Route
                    </Link>
                </div>
            ) : (
                <div className="flex flex-col gap-4">
                    <h2 className="text-sm font-bold text-white/50 uppercase tracking-wider pl-1 font-mono">Live Telemetry Arrays</h2>
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                        {keys.map(key => (
                            <APIStatsCard key={key.id} apiKey={key} />
                        ))}
                    </div>
                </div>
            )}
        </div>
    );
}
