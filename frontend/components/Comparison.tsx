'use client';

import { useState, useEffect } from 'react';

type Job = {
  id: number;
  startTime: number;
};

export default function Comparison() {
  const [tradJobs, setTradJobs] = useState<Job[]>([]);
  const [veloxJobs, setVeloxJobs] = useState<Job[]>([]);
  const [globalTime, setGlobalTime] = useState(Date.now());
  
  useEffect(() => {
    // Shared clock for real-time updates
    const clockInterval = setInterval(() => {
      setGlobalTime(Date.now());
    }, 100);

    // Initial state population
    const now = Date.now();
    setTradJobs([
      { id: 89213, startTime: now + 200 },
      { id: 89212, startTime: now + 100 },
      { id: 89211, startTime: now }
    ]);
    
    setVeloxJobs([{ id: 1, startTime: now }]);

    // Job spawner: 1 job every 500ms
    const spawner = setInterval(() => {
      const currentTime = Date.now();
      
      setTradJobs(prev => {
        const nextId = prev.length > 0 ? prev[0].id + 1 : 89211;
        const newJob = { id: nextId, startTime: currentTime };
        // Prepend to top, keep last 25 jobs
        return [newJob, ...prev].slice(0, 25); 
      });

      setVeloxJobs(prev => {
        const nextId = prev.length > 0 ? prev[0].id + 1 : 1;
        const newJob = { id: nextId, startTime: currentTime };
        return [newJob, ...prev].slice(0, 25); 
      });

    }, 500);

    return () => {
      clearInterval(clockInterval);
      clearInterval(spawner);
    };
  }, []);

  const formatVeloxId = (id: number) => id.toString().padStart(5, '0');

  // Calculate estimated total traditional backlog time
  const baseMinutes = 4;
  const simulatedExtraSeconds = Math.floor((globalTime - (tradJobs[tradJobs.length - 1]?.startTime || globalTime)) / 200);
  const totalSeconds = (baseMinutes * 60) + 12 + simulatedExtraSeconds;
  const estMins = Math.floor(totalSeconds / 60);
  const estSecs = (totalSeconds % 60).toString().padStart(2, '0');
  
  return (
    <section className="bg-background py-24 relative overflow-hidden" id="compare">
      {/* Background divider glow */}
      <div className="absolute top-0 inset-x-0 h-[1px] bg-gradient-to-r from-transparent via-white/10 to-transparent"></div>
      
      <div className="mx-auto max-w-7xl px-6 lg:px-8 relative z-10">
        <div className="mx-auto max-w-2xl lg:text-center mb-20 text-center">
          <h2 className="text-sm font-bold leading-7 text-primary tracking-widest uppercase mb-3">The Unfair Advantage</h2>
          <p className="text-4xl font-extrabold tracking-tight text-white mb-6">
            Why wait in queue?
          </p>
          <p className="text-lg leading-8 text-foreground/60 max-w-xl mx-auto">
            Traditional judging platforms buckle under tournament traffic, leaving users waiting minutes just to see a compilation error. We fixed that.
          </p>
        </div>

        <div className="grid md:grid-cols-2 gap-8 lg:gap-12 items-center relative">
          
          {/* Subtle connecting lines */}
          <div className="hidden md:block absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 w-32 h-[1px] bg-gradient-to-r from-white/5 via-primary/30 to-white/5 z-0"></div>

          {/* Traditional Platform Mockup */}
          <div className="dev-card p-8 bg-surface/40 backdrop-blur-sm z-10">
            <div className="flex items-center justify-between border-b border-white/5 pb-4 mb-6">
              <div className="flex items-center gap-3">
                <div className="w-2 h-2 rounded-full bg-red-400"></div>
                <h3 className="text-lg font-bold text-white/40 line-through">Traditional Platform</h3>
              </div>
              <span className="text-xs font-mono text-red-400/80 bg-red-400/10 px-2 py-1 rounded-md">Queue Delay</span>
            </div>
            
            <p className="mb-6 text-sm text-foreground/50">Your code is held up in massive, unscalable databases during contests.</p>
              
            <div 
              className="bg-black/40 rounded-xl p-4 border border-white/5 relative overflow-hidden flex flex-col h-[260px]"
            >
               <div className="absolute inset-0 bg-red-500/5 blur-xl pointer-events-none"></div>
              
               {tradJobs.map((job, idx) => {
                 // For Traditional, jobs are stuck in queue. 
                 const isExecuting = idx === tradJobs.length - 1; // Oldest job is executing
                 // Real-time elapsed counter for the executing/pending job
                 const elapsed = Math.max(0, (globalTime - job.startTime) / 1000).toFixed(1);
                 
                 return (
                   <div key={job.id} className="animate-feed-in mb-3 flex flex-col gap-2 p-3 rounded-lg bg-white/[0.02] border border-white/5 relative z-10 shrink-0">
                     <div className="flex justify-between text-xs text-white/50 font-mono">
                       <span>Job #{job.id}</span>
                       <span className="text-red-300/80">{isExecuting ? 'Executing...' : 'Pending Queue'}</span>
                     </div>
                     <div className="flex items-center justify-between">
                       <div className="flex items-center gap-2 opacity-60">
                         <div className="animate-spin h-3 w-3 border-2 border-red-400 border-t-transparent rounded-full"></div>
                         <span className="text-xs font-mono text-red-200/50">Wait: {elapsed}s</span>
                       </div>
                       <span className="text-xs font-mono text-red-400/80">{isExecuting ? `Est: ${estMins}m ${estSecs}s` : 'Unknown'}</span>
                     </div>
                   </div>
                 );
               })}
            </div>
          </div>

          {/* Velox Mockup */}
          <div className="dev-card p-8 bg-surface/80 border-primary/30 shadow-[0_0_50px_rgba(255,90,0,0.1)] relative z-20 md:scale-105">
            {/* Inner primary glow */}
            <div className="absolute inset-x-0 top-0 h-1/2 bg-gradient-to-b from-primary/10 to-transparent pointer-events-none rounded-t-2xl"></div>
            
            <div className="flex items-center justify-between border-b border-white/10 pb-4 mb-6 relative">
              <div className="flex items-center gap-3">
                <div className="w-2 h-2 rounded-full bg-primary animate-pulse shadow-[0_0_10px_rgba(255,90,0,0.8)]"></div>
                <h3 className="text-xl font-bold text-white tracking-wide">Velox Engine</h3>
              </div>
              <span className="text-xs font-mono text-primary bg-primary/10 px-3 py-1 rounded-md border border-primary/20 shadow-[0_0_10px_rgba(255,90,0,0.2)]">Instant Execution</span>
            </div>
            
            <p className="mb-6 text-sm text-foreground/80 leading-relaxed">Independent workers scale instantly to handle any traffic spike. Zero delays.</p>
              
            <div 
              className="bg-black/60 rounded-xl p-4 border border-white/10 relative overflow-hidden shadow-inner flex flex-col h-[260px]"
            >
               {veloxJobs.map((job) => {
                 return (
                   <div key={job.id} className="animate-feed-in mb-3 flex flex-col gap-2 p-3 rounded-lg bg-white/[0.04] border border-white/10 shrink-0">
                     <div className="flex justify-between text-xs text-white/80 font-mono">
                       <span>Job #{formatVeloxId(job.id)}</span>
                       <span className="text-success font-bold tracking-wide">Done</span>
                     </div>
                     <div className="flex items-center justify-between">
                       <div className="flex items-center gap-2">
                         <div className="h-4 w-4 bg-success/20 rounded-full flex items-center justify-center border border-success/30 shadow-[0_0_10px_rgba(34,197,94,0.3)]">
                           <svg className="w-2.5 h-2.5 text-success" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth="3">
                             <path strokeLinecap="round" strokeLinejoin="round" d="M5 13l4 4L19 7" />
                           </svg>
                         </div>
                         <span className="text-xs font-mono text-success/80">Accepted</span>
                       </div>
                       <span className="text-xs font-mono text-primary font-bold">12ms</span>
                     </div>
                   </div>
                 );
               })}
            </div>
          </div>

        </div>
      </div>
    </section>
  );
}
