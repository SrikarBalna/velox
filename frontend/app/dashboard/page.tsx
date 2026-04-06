'use client';

import React from 'react';
import StatCard from '@/components/dashboard/StatCard';
import { FocusProvider, useFocus } from '@/components/dashboard/FocusContext';
import InteractiveStat from '@/components/dashboard/InteractiveStat';

function FocusBackdrop() {
    const { focusedId } = useFocus();
    return (
        <div 
            className={`fixed inset-0 z-40 bg-black/70 backdrop-blur-md transition-opacity duration-500 pointer-events-none ${focusedId ? 'opacity-100' : 'opacity-0'}`}
        ></div>
    );
}

function DashboardContent() {
  const genericModalContent = (
      <div className="flex flex-col gap-6">
          <div className="grid grid-cols-2 gap-4">
              <div className="bg-black/50 p-6 rounded-xl border border-white/5">
                  <p className="text-white/50 text-sm mb-2">Trend (7d)</p>
                  <p className="text-4xl text-white font-black">+0%</p>
              </div>
              <div className="bg-black/50 p-6 rounded-xl border border-white/5">
                  <p className="text-white/50 text-sm mb-2">Variance</p>
                  <p className="text-4xl text-success font-black">0.00%</p>
              </div>
          </div>
          <div className="h-48 w-full bg-white/5 rounded-xl border border-white/5 flex items-center justify-center">
              <p className="text-white/30 font-mono text-sm">Insufficient data for historical plot</p>
          </div>
      </div>
  );

  return (
    <div className="flex flex-col gap-6">
      <FocusBackdrop />
      <div className="flex items-center justify-between mb-2">
        <h1 className="text-3xl font-bold text-white relative z-10">Command Center</h1>
      </div>

      <div className="grid grid-cols-1 xl:grid-cols-4 gap-6">
        {/* ========================================================= */}
        {/* ZONE A: MAIN ANALYTICS (Left Side)                        */}
        {/* ========================================================= */}
        <div className="xl:col-span-3 flex flex-col gap-6">
            
            {/* Metric Cards Row */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                <InteractiveStat id="stat-executions" title="Total Executions" modalContent={genericModalContent}>
                    <div className="h-full"><StatCard title="Total Executions" value="0" trend="0%" isPositive={true} icon="M13 10V3L4 14h7v7l9-11h-7z" /></div>
                </InteractiveStat>
                <InteractiveStat id="stat-latency" title="Avg Latency" modalContent={genericModalContent}>
                    <div className="h-full"><StatCard title="Avg Latency" value="--ms" trend="0%" isPositive={true} icon="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" /></div>
                </InteractiveStat>
                <InteractiveStat id="stat-success" title="Success Rate" modalContent={genericModalContent}>
                    <div className="h-full"><StatCard title="Success Rate" value="--%" trend="0%" isPositive={true} icon="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" /></div>
                </InteractiveStat>
                <InteractiveStat id="stat-workers" title="Active Workers" modalContent={genericModalContent}>
                    <div className="h-full"><StatCard title="Active Workers" value="0" trend="0" isPositive={false} icon="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 00-2-2h6a2 2 0 012 2v2M7 7h10" /></div>
                </InteractiveStat>
            </div>

            {/* Execution Volume Graph */}
            <InteractiveStat id="chart-volume" title="Execution Volume Details" modalContent={
                <div className="h-[400px] w-full bg-white/5 rounded-xl border border-white/5 flex items-center justify-center">
                    <p className="text-white/30 font-mono text-sm">Detailed Execution Chart Unavailable</p>
                </div>
            }>
                <div className="dev-card p-6 bg-surface/40 flex flex-col min-h-[350px]">
                    <div className="flex justify-between items-center mb-6">
                        <h3 className="text-lg font-bold text-white">Execution Volume</h3>
                        <select className="bg-black border border-white/10 text-xs text-white/70 px-3 py-1.5 rounded-md outline-none relative z-20 pointer-events-auto">
                            <option>Last 24 Hours</option>
                            <option>Last 7 Days</option>
                            <option>Last 30 Days</option>
                        </select>
                    </div>
                    
                    <div className="flex-1 relative mt-4">
                        <div className="absolute inset-0 flex flex-col justify-between pointer-events-none">
                            {[0, 1, 2, 3].map(i => (
                                <div key={i} className="w-full h-px bg-white/5"></div>
                            ))}
                        </div>
                        
                        <div className="absolute inset-x-0 bottom-0 h-full flex items-center justify-center z-10 pointer-events-none">
                            <div className="text-center">
                                <svg className="w-8 h-8 text-white/20 mx-auto mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="1" d="M13 10V3L4 14h7v7l9-11h-7z" />
                                </svg>
                                <p className="text-sm font-medium text-white/50">No execution data available</p>
                                <p className="text-xs text-white/30 mt-1 font-mono">Run your first job to see analytics</p>
                            </div>
                        </div>
                        
                        <div className="absolute inset-x-0 bottom-0 h-1/3 bg-gradient-to-b from-transparent to-surface pointer-events-none"></div>
                    </div>
                </div>
            </InteractiveStat>

            {/* Recent Computations List */}
            <div className="dev-card p-0 bg-surface/40 flex flex-col overflow-hidden min-h-[350px] relative z-10">
                <div className="p-6 border-b border-white/5 flex justify-between items-center">
                    <h3 className="text-lg font-bold text-white">Live Executions</h3>
                    <div className="flex gap-1.5">
                        <div className="w-1.5 h-1.5 rounded-full bg-success opacity-50 animate-pulse"></div>
                        <div className="w-1.5 h-1.5 rounded-full bg-success opacity-50 animate-pulse" style={{ animationDelay: '200ms' }}></div>
                        <div className="w-1.5 h-1.5 rounded-full bg-success opacity-50 animate-pulse" style={{ animationDelay: '400ms' }}></div>
                    </div>
                </div>
                
                <div className="flex-1 flex flex-col items-center justify-center p-6 text-center">
                    <svg className="w-8 h-8 text-white/20 mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="1" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                    </svg>
                    <p className="text-sm font-medium text-white/50">No computations yet</p>
                    <p className="text-xs text-white/30 mt-1 font-mono">Jobs will appear here as they run</p>
                </div>
            </div>

        </div>

        {/* ========================================================= */}
        {/* ZONE B: PROFILE & SETTINGS CONTEXT (Right Side)           */}
        {/* ========================================================= */}
        <div className="xl:col-span-1 flex flex-col gap-6 relative z-10">
            {/* Context cards remain non-interactive stats so they don't get focused */}
            
            {/* Profile Details */}
            <div className="dev-card p-6 bg-surface/40 flex flex-col">
                <h3 className="text-lg font-bold text-white mb-2">Workspace Profile</h3>
                <div className="flex flex-col items-center text-center mt-4 mb-4">
                    <div className="w-20 h-20 rounded-full bg-white/5 flex items-center justify-center p-1 border border-white/10 shadow-inner mb-4">
                        <img src="https://api.dicebear.com/7.x/avataaars/svg?seed=Admin" alt="Avatar" className="w-full h-full rounded-full" />
                    </div>
                    <h4 className="text-xl font-bold text-white">Admin User</h4>
                    <p className="text-sm text-foreground/50 font-mono mt-1">admin@velox.dev</p>
                </div>
                <button className="mt-auto px-4 py-2 bg-white/5 hover:bg-white/10 border border-white/10 rounded-lg text-xs font-bold text-white transition-colors cursor-pointer w-full">Edit Profile</button>
            </div>

            {/* Workspace Tier */}
            <div className="dev-card p-6 bg-surface/40 relative overflow-hidden flex flex-col border-primary/20">
                <div className="absolute -top-10 -right-10 w-32 h-32 bg-primary/20 blur-3xl rounded-full pointer-events-none"></div>
                
                <div className="flex justify-between items-start mb-4">
                    <h3 className="text-lg font-bold text-white relative z-10">Workspace Plan</h3>
                    <span className="text-xs font-bold font-mono py-1 px-3 rounded-full text-primary bg-primary/10 border border-primary/20 relative z-10">Active</span>
                </div>
                <div className="mt-2 mb-6 relative z-10">
                    <p className="text-4xl font-black text-white">Pro Tier</p>
                    <p className="text-sm text-foreground/50 mt-3">1.24M executions included per month. Priority worker access enabled.</p>
                </div>
                <button className="mt-auto w-full py-2.5 bg-primary hover:bg-primary-hover rounded-lg text-sm font-bold text-black transition-colors cursor-pointer shadow-[0_0_15px_rgba(255,90,0,0.3)] relative z-10">
                    Manage Billing
                </button>
            </div>

            {/* API Keys */}
            <div className="dev-card p-6 bg-surface/40 flex flex-col flex-1">
                <h3 className="text-lg font-bold text-white mb-2">API Keys</h3>
                <p className="text-sm text-foreground/50">Use these keys to authenticate your requests to the Velox Engine.</p>
                
                <div className="mt-6 mb-6">
                    <label className="text-xs font-bold text-white/50 uppercase tracking-wider mb-2 block">Production Key</label>
                    <div className="flex flex-col xl:flex-row gap-2">
                        <input type="password" value="vlx_live_98a72b4fc10" readOnly className="flex-1 w-full bg-black border border-white/10 rounded-lg px-3 py-2 text-sm text-white font-mono focus:outline-none" />
                        <button className="p-2 bg-white/5 hover:bg-white/10 border border-white/10 rounded-lg transition-colors group cursor-pointer flex justify-center items-center shrink-0">
                            <svg className="w-5 h-5 text-white/50 group-hover:text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                            </svg>
                        </button>
                    </div>
                </div>
                
                <button className="mt-auto text-sm font-bold text-primary hover:text-white transition-colors text-left cursor-pointer">Generate new key &rarr;</button>
            </div>
        </div>
      </div>
    </div>
  );
}

export default function DashboardOverview() {
    return (
        <FocusProvider>
            <DashboardContent />
        </FocusProvider>
    );
}
