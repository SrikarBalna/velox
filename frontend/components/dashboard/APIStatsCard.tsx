import React, { useEffect, useState } from 'react';
import api from '@/lib/api';

interface APIStatsCardProps {
    apiKey: {
        id: string;
        project: string;
        preview: string;
        createdAt: string;
    };
}

export default function APIStatsCard({ apiKey }: APIStatsCardProps) {
    const [stats, setStats] = useState<any>(null);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        let isMounted = true;
        
        const fetchStats = async () => {
            setIsLoading(true);
            try {
                const res = await api.get(`/auth/api-keys/stats?id=${apiKey.id}`);
                if (isMounted && res.data) {
                    setStats(res.data);
                }
            } catch (err) {
                console.error(`Failed to load stats for key ${apiKey.id}:`, err);
                // Fallback to empty if it fails so UI doesn't crash completely.
                if (isMounted) {
                    setStats({ total_requests: 0, success_rate: 0, peak_rpm: 0 });
                }
            } finally {
                if (isMounted) setIsLoading(false);
            }
        };

        fetchStats();

        return () => {
            isMounted = false;
        };
    }, [apiKey.id]);

    return (
        <div className="dev-card p-6 bg-surface/40 flex flex-col border border-white/5 relative overflow-hidden group">
            {/* Soft backdrop glow to maintain premium feeling */}
            <div className="absolute -top-10 -right-10 w-32 h-32 bg-primary/5 blur-3xl rounded-full pointer-events-none group-hover:bg-primary/10 transition-colors"></div>
            
            <div className="flex justify-between items-start mb-6 border-b border-white/5 pb-4">
                <div>
                    <div className="flex items-center gap-2 mb-1 relative z-10">
                        <div className="w-2 h-2 rounded-full bg-primary animate-pulse"></div>
                        <h3 className="text-lg font-bold text-white">{apiKey.project}</h3>
                    </div>
                    <p className="text-xs text-foreground/50 font-mono tracking-wide relative z-10">{apiKey.preview}</p>
                </div>
            </div>

            <div className="flex-1">
                {isLoading ? (
                    <div className="flex flex-col items-center justify-center h-full py-8 text-white/30">
                        <div className="w-5 h-5 border-2 border-white/10 border-t-white rounded-full animate-spin mb-3"></div>
                        <span className="text-xs font-bold font-mono tracking-widest uppercase">Syncing Telemetry</span>
                    </div>
                ) : (
                    <div className="grid grid-cols-2 gap-4">
                        <div className="bg-black/40 p-4 rounded-xl border border-white/5 flex flex-col justify-center">
                            <p className="text-white/40 text-xs font-bold uppercase tracking-wider mb-2">Total Requests</p>
                            <p className="text-2xl text-white font-black font-mono">
                                {stats?.total_requests?.toLocaleString() || 0}
                            </p>
                        </div>
                        
                        <div className="bg-black/40 p-4 rounded-xl border border-white/5 flex flex-col justify-center">
                            <p className="text-white/40 text-xs font-bold uppercase tracking-wider mb-2">Peak RPM</p>
                            <p className="text-2xl text-white font-black font-mono">
                                {stats?.peak_rpm?.toLocaleString() || 0} <span className="text-sm text-white/30 font-medium">/min</span>
                            </p>
                        </div>

                        <div className="bg-black/40 p-4 rounded-xl border border-white/5 flex flex-col justify-center col-span-2">
                             <div className="flex justify-between items-end mb-2">
                                <p className="text-white/40 text-xs font-bold uppercase tracking-wider">Success Rate</p>
                                <span className={`text-sm font-bold ${stats?.success_rate >= 90 ? 'text-success' : 'text-primary'}`}>
                                    {stats?.success_rate ? stats.success_rate.toFixed(1) : '0.0'}%
                                </span>
                             </div>
                             
                             <div className="w-full h-1.5 bg-white/10 rounded-full overflow-hidden">
                                <div 
                                    className={`h-full rounded-full transition-all duration-1000 ${stats?.success_rate >= 90 ? 'bg-success' : 'bg-primary'}`} 
                                    style={{ width: `${stats?.success_rate || 0}%` }}
                                ></div>
                             </div>
                        </div>
                    </div>
                )}
            </div>
            
            <div className="mt-6 text-xs text-white/30 font-mono flex items-center gap-1.5">
                <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                Registered: {apiKey.createdAt}
            </div>
        </div>
    );
}
