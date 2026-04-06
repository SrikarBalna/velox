import Link from 'next/link';
import LogoutButton from '@/components/dashboard/LogoutButton';

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="min-h-screen bg-background flex flex-col">
      <header className="h-20 border-b border-white/5 bg-surface/50 backdrop-blur-md flex items-center justify-between px-8 sticky top-0 z-20">
        <Link href="/" className="text-2xl font-black text-white hover:text-primary transition-colors flex items-center gap-2">
          <div className="w-8 h-8 rounded-lg bg-surface border border-white/10 flex flex-col gap-[2px] items-center justify-center relative overflow-hidden shadow-sm">
             <div className="absolute inset-0 bg-primary/20 animate-pulse"></div>
             <div className="flex gap-[3px]">
               <div className="w-[3px] h-[12px] bg-primary rounded-full transform -skew-x-12"></div>
               <div className="w-[3px] h-[12px] bg-white rounded-full transform -skew-x-12"></div>
             </div>
          </div>
          Velox
        </Link>
        
        <div className="flex items-center gap-6">
          <button className="px-5 py-2.5 bg-primary/10 hover:bg-primary/20 border border-primary/30 text-primary hover:text-primary-hover font-bold text-sm rounded-lg transition-all shadow-[0_0_15px_rgba(255,90,0,0.2)] hover:shadow-[0_0_25px_rgba(255,90,0,0.4)] flex items-center gap-2 cursor-pointer">
            <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth="2.5">
              <path strokeLinecap="round" strokeLinejoin="round" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
              <path strokeLinecap="round" strokeLinejoin="round" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            Open Sandbox
          </button>
          
          <div className="w-px h-6 bg-white/10"></div>
          
          <Link href="/dashboard/profile" className="flex items-center gap-3 cursor-pointer group">
            <div className="flex flex-col text-right hidden sm:flex">
                <span className="text-sm font-bold text-white">Admin</span>
                <span className="text-xs text-white/50">Pro Tier</span>
            </div>
            <div className="w-10 h-10 rounded-full bg-white/10 border border-white/20 flex items-center justify-center text-sm font-bold text-white group-hover:border-primary/50 transition-colors shadow-sm overflow-hidden relative">
              <img src="https://api.dicebear.com/7.x/avataaars/svg?seed=Admin" alt="Avatar" className="w-full h-full object-cover" />
            </div>
          </Link>
          <LogoutButton />
        </div>
      </header>

      <main className="flex-1 p-8 overflow-y-auto w-full">
        <div className="max-w-7xl mx-auto space-y-8">
          {children}
        </div>
      </main>
    </div>
  );
}
