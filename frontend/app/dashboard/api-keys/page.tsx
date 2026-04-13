'use client';

import React, { useState, useEffect } from 'react';
import api from '@/lib/api';
import Link from 'next/link';
import APIKeyModal from '@/components/dashboard/APIKeyModal';

interface MockKey {
  id: string;
  project: string;
  preview: string;
  createdAt: string;
}

const initialKeys: MockKey[] = [
  { id: '1', project: 'Velox Production', preview: 'vlx_live_...fc10', createdAt: 'Mar 15, 2026' },
  { id: '2', project: 'Default Workspace', preview: 'vlx_live_...092a', createdAt: 'Feb 28, 2026' },
];

export default function APIKeysPage() {
  const [keys, setKeys] = useState<MockKey[]>([]);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    fetchKeys();
  }, []);

  const fetchKeys = async () => {
    setIsLoading(true);
    try {
      const res = await api.get('/auth/api-keys');
      // Backend returns array of: { id, name, display_hint, created_at, ... }
      if (res.data) {
        const formattedKeys = res.data.map((k: any) => ({
          id: k.id,
          project: k.name || 'Unnamed Project',
          preview: k.display_hint || 'vlx_...',
          createdAt: new Date(k.created_at).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
        }));
        setKeys(formattedKeys);
      }
    } catch (err) {
      console.error("Failed to fetch SDK keys:", err);
    } finally {
      setIsLoading(false);
    }
  };

  // Extract unique projects to pass to the modal
  const existingProjects = Array.from(new Set(keys.map(k => k.project)));

  const handleKeyGenerated = () => {
    // Re-fetch the keys since the backend just generated one
    fetchKeys();
  };

  const handleDelete = async (id: string) => {
    try {
      await api.delete(`/auth/api-keys?id=${id}`);
      setKeys(keys.filter(k => k.id !== id));
    } catch (err) {
      console.error("Failed to delete key:", err);
    }
  };

  return (
    <div className="flex flex-col gap-8 max-w-5xl mx-auto w-full pt-4">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <div className="flex items-center gap-2 mb-2">
            <Link href="/dashboard" className="text-white/50 hover:text-white transition-colors text-sm font-bold flex items-center gap-1">
              <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M15 19l-7-7 7-7" />
              </svg>
              Overview
            </Link>
          </div>
          <h1 className="text-3xl font-bold text-white tracking-tight">API Keys</h1>
          <p className="text-foreground/60 text-sm mt-1">Manage your API keys to authenticate with the Velox Execution Engine.</p>
        </div>
        
        <button 
          onClick={() => setIsModalOpen(true)}
          className="px-5 py-2.5 rounded-lg bg-primary hover:bg-primary-hover text-black font-bold text-sm transition-all shadow-[0_0_15px_rgba(255,90,0,0.3)] shrink-0 flex items-center gap-2"
        >
          <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
             <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="3" d="M12 4v16m8-8H4" />
          </svg>
          Create API Key
        </button>
      </div>

      {/* Keys Table Container */}
      <div className="dev-card p-0 bg-surface/40 overflow-hidden border-white/5">
        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="border-b border-white/5 bg-black/20">
                <th className="px-6 py-4 text-xs font-bold text-white/50 uppercase tracking-wider">Project</th>
                <th className="px-6 py-4 text-xs font-bold text-white/50 uppercase tracking-wider">Key preview</th>
                <th className="px-6 py-4 text-xs font-bold text-white/50 uppercase tracking-wider">Created</th>
                <th className="px-6 py-4 text-xs font-bold text-white/50 uppercase tracking-wider text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-white/5">
              {isLoading ? (
                <tr>
                  <td colSpan={4} className="px-6 py-12 text-center text-white/50">
                    <div className="w-6 h-6 border-2 border-white/20 border-t-white rounded-full animate-spin mx-auto mb-2"></div>
                    <p className="text-sm font-bold">Loading Keys...</p>
                  </td>
                </tr>
              ) : keys.length === 0 ? (
                <tr>
                  <td colSpan={4} className="px-6 py-12 text-center">
                    <div className="w-12 h-12 rounded-full bg-white/5 flex items-center justify-center mx-auto mb-3">
                      <svg className="w-6 h-6 text-white/20" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                         <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
                      </svg>
                    </div>
                    <p className="text-sm font-bold text-white">No API keys yet</p>
                    <p className="text-xs text-white/40 mt-1">Create an API key to get started.</p>
                  </td>
                </tr>
              ) : (
                keys.map((key) => (
                  <tr key={key.id} className="hover:bg-white/[0.02] transition-colors group">
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="flex items-center gap-3">
                        <div className="w-8 h-8 rounded bg-white/5 border border-white/10 flex items-center justify-center shrink-0">
                          <span className="text-xs font-bold text-white">{key.project.charAt(0).toUpperCase()}</span>
                        </div>
                        <span className="text-sm font-bold text-white">{key.project}</span>
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-foreground/70 font-mono">
                      {key.preview}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-foreground/50">
                      {key.createdAt}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-right">
                      <button 
                        onClick={() => handleDelete(key.id)}
                        className="text-xs font-bold text-red-500/70 hover:text-red-500 transition-colors opacity-0 group-hover:opacity-100 p-2 rounded-md hover:bg-red-500/10"
                        title="Delete Key"
                      >
                         Delete
                      </button>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>

      <APIKeyModal 
        isOpen={isModalOpen} 
        onClose={() => setIsModalOpen(false)} 
        onKeyGenerated={handleKeyGenerated} 
        existingProjects={existingProjects}
      />
    </div>
  );
}
