import Link from 'next/link';

export default function ProfilePage() {
  return (
    <div className="flex flex-col gap-6">
      <div className="flex items-center gap-4 mb-2">
        <Link href="/dashboard" className="text-white/50 hover:text-white transition-colors bg-white/5 hover:bg-white/10 p-2 rounded-lg">
            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
            </svg>
        </Link>
        <h1 className="text-3xl font-bold text-white">Profile Settings</h1>
      </div>

      <div className="grid grid-cols-1 xl:grid-cols-3 gap-6">
        
        {/* Main Settings Column */}
        <div className="xl:col-span-2 flex flex-col gap-6">
            
            {/* General Info Form */}
            <div className="dev-card p-6 bg-surface/40 flex flex-col">
                <h3 className="text-xl font-bold text-white mb-6">General Information</h3>
                
                <form className="flex flex-col gap-5">
                    <div className="flex flex-col gap-2">
                        <label className="text-sm font-bold text-white/50 uppercase tracking-wider">Display Name</label>
                        <input type="text" defaultValue="Admin User" className="w-full bg-black border border-white/10 rounded-lg px-4 py-3 text-white focus:outline-none focus:border-primary/50 focus:ring-1 focus:ring-primary/50 transition-all font-mono" />
                    </div>
                    
                    <div className="flex flex-col gap-2">
                        <label className="text-sm font-bold text-white/50 uppercase tracking-wider">Email Address</label>
                        <input type="email" defaultValue="admin@velox.dev" className="w-full bg-black border border-white/10 rounded-lg px-4 py-3 text-white focus:outline-none focus:border-primary/50 focus:ring-1 focus:ring-primary/50 transition-all font-mono" />
                    </div>
                    
                    <div className="flex justify-end mt-2">
                        <button type="button" className="px-6 py-2.5 bg-white text-black font-bold rounded-lg hover:bg-gray-200 transition-colors shadow-[0_0_15px_rgba(255,255,255,0.1)]">
                            Save Changes
                        </button>
                    </div>
                </form>
            </div>

            {/* Notification Preferences */}
            <div className="dev-card p-6 bg-surface/40 flex flex-col">
                <h3 className="text-xl font-bold text-white mb-6">Notification Preferences</h3>
                
                <div className="flex flex-col gap-4">
                    <div className="flex items-center justify-between p-4 bg-black/50 rounded-xl border border-white/5">
                        <div>
                            <h4 className="text-white font-bold">Execution Failed Alerts</h4>
                            <p className="text-sm text-white/40 mt-1">Receive an email immediately if a worker drops a job.</p>
                        </div>
                        <div className="w-12 h-6 bg-primary rounded-full relative cursor-pointer shadow-[0_0_10px_rgba(255,90,0,0.3)]">
                            <div className="w-4 h-4 bg-black rounded-full absolute right-1 top-1"></div>
                        </div>
                    </div>
                    
                    <div className="flex items-center justify-between p-4 bg-black/50 rounded-xl border border-white/5">
                        <div>
                            <h4 className="text-white font-bold">Weekly Analytics Report</h4>
                            <p className="text-sm text-white/40 mt-1">A summary of your latency and execution volume.</p>
                        </div>
                        <div className="w-12 h-6 bg-white/10 rounded-full relative cursor-pointer">
                            <div className="w-4 h-4 bg-white/50 rounded-full absolute left-1 top-1"></div>
                        </div>
                    </div>
                </div>
            </div>

            {/* Danger Zone */}
            <div className="dev-card p-6 bg-red-950/10 border-red-500/20 flex flex-col">
                <h3 className="text-xl font-bold text-red-500 mb-2">Danger Zone</h3>
                <p className="text-sm text-white/50 mb-6">Irreversible and destructive actions concerning your account.</p>
                
                <div className="flex items-center justify-between p-4 bg-black/50 rounded-xl border border-red-500/10">
                    <div>
                        <h4 className="text-white font-bold">Delete Account</h4>
                        <p className="text-sm text-white/40 mt-1">Permanently remove your account and all execution data.</p>
                    </div>
                    <button className="px-4 py-2 bg-red-500/10 hover:bg-red-500/20 text-red-500 font-bold rounded-lg border border-red-500/20 transition-colors">
                        Delete Account
                    </button>
                </div>
            </div>

        </div>

        {/* Side Column */}
        <div className="xl:col-span-1 flex flex-col gap-6">
            
            {/* Avatar block */}
            <div className="dev-card p-6 bg-surface/40 flex flex-col text-center">
                <h3 className="text-lg font-bold text-white mb-6 text-left">Avatar</h3>
                <div className="w-32 h-32 mx-auto rounded-full bg-white/5 flex items-center justify-center p-2 border border-white/10 shadow-inner mb-6 relative group cursor-pointer">
                    <img src="https://api.dicebear.com/7.x/avataaars/svg?seed=Admin" alt="Avatar" className="w-full h-full rounded-full group-hover:opacity-50 transition-opacity" />
                    <div className="absolute inset-0 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
                        <svg className="w-8 h-8 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
                        </svg>
                    </div>
                </div>
                <p className="text-xs text-white/40 font-mono">Click to upload new image</p>
                <p className="text-xs text-white/30 font-mono mt-1">JPEG or PNG, max 2MB.</p>
            </div>

            {/* Security Info */}
             <div className="dev-card p-6 bg-surface/40 flex flex-col">
                <h3 className="text-lg font-bold text-white mb-4 lg:mb-6 text-left">Security</h3>
                <div className="flex flex-col gap-4">
                    <button className="w-full py-3 bg-white/5 hover:bg-white/10 border border-white/10 text-white font-bold rounded-lg transition-colors text-sm">
                        Change Password
                    </button>
                    <button className="w-full py-3 bg-white/5 hover:bg-white/10 border border-white/10 text-white font-bold rounded-lg transition-colors text-sm flex items-center justify-center gap-2">
                        <svg className="w-4 h-4 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
                        </svg>
                        Enable 2FA Auth
                    </button>
                </div>
            </div>

        </div>

      </div>
    </div>
  );
}
