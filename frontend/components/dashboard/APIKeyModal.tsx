'use client';

import React, { useState } from 'react';
import api from '@/lib/api';

interface APIKeyModalProps {
  isOpen: boolean;
  onClose: () => void;
  onKeyGenerated: (project: string, keyPreview: string) => void;
  existingProjects: string[];
}

export default function APIKeyModal({ isOpen, onClose, onKeyGenerated, existingProjects }: APIKeyModalProps) {
  const [step, setStep] = useState<1 | 2>(1);
  const [selectedProject, setSelectedProject] = useState<string>('');
  const [isCreatingNew, setIsCreatingNew] = useState(false);
  const [newProjectName, setNewProjectName] = useState('');
  const [generatedKey, setGeneratedKey] = useState<string>('');
  const [copied, setCopied] = useState(false);

  const [isGenerating, setIsGenerating] = useState(false);
  const [error, setError] = useState<string | null>(null);

  if (!isOpen) return null;

  const handleGenerate = async () => {
    const projectToUse = isCreatingNew ? newProjectName : selectedProject;
    if (!projectToUse) return;

    setIsGenerating(true);
    setError(null);
    try {
      // Hit the real backend to generate an API key
      // Let the 'project name' act as the internal API key name for our mapping
      const res = await api.post('/auth/api-keys', {
        name: projectToUse,
        scopes: ['submit', 'status']
      });
      // The backend returns { key, id, display_hint, ... }
      setGeneratedKey(res.data.key);
      setStep(2);
    } catch (err: any) {
      console.error("Failed to generate key:", err);
      setError(err.response?.data?.error || "Failed to generate API Key.");
    } finally {
      setIsGenerating(false);
    }
  };

  const handleCopy = () => {
    navigator.clipboard.writeText(generatedKey);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const handleClose = () => {
    if (step === 2) {
      // Return the generated key info back to the parent to append to the table
      const preview = `vlx_live_${generatedKey.slice(-4)}`;
      onKeyGenerated(isCreatingNew ? newProjectName : selectedProject, preview);
    }
    // Reset state
    setStep(1);
    setSelectedProject('');
    setIsCreatingNew(false);
    setNewProjectName('');
    setGeneratedKey('');
    setCopied(false);
    onClose();
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div className="absolute inset-0 bg-black/80 backdrop-blur-sm" onClick={handleClose}></div>
      
      <div className="relative w-full max-w-md bg-[#111] border border-white/10 rounded-2xl shadow-2xl flex flex-col overflow-hidden animate-in fade-in zoom-in-95 duration-200">
        
        {step === 1 && (
          <>
            <div className="p-6 border-b border-white/5">
              <h2 className="text-xl font-bold text-white tracking-tight">Create API key</h2>
              <p className="text-sm text-foreground/60 mt-2">API keys allow you to authenticate your code execution requests on the Velox engine.</p>
            </div>

            <div className="p-6 flex flex-col gap-6">
              {!isCreatingNew ? (
                <div className="flex flex-col gap-3">
                  <label className="text-sm font-bold text-white">Select a project</label>
                  <select 
                    value={selectedProject}
                    onChange={(e) => setSelectedProject(e.target.value)}
                    className="w-full bg-black/50 border border-white/10 rounded-lg px-4 py-3 text-sm text-white focus:outline-none focus:border-primary/50 transition-colors"
                  >
                    <option value="" disabled>Select a project...</option>
                    {existingProjects.map(p => (
                      <option key={p} value={p}>{p}</option>
                    ))}
                  </select>
                  
                  <div className="relative flex items-center py-2">
                    <div className="flex-grow border-t border-white/5"></div>
                    <span className="flex-shrink-0 mx-4 text-xs font-bold text-white/30 uppercase tracking-widest">OR</span>
                    <div className="flex-grow border-t border-white/5"></div>
                  </div>

                  <button 
                    onClick={() => setIsCreatingNew(true)}
                    className="text-sm font-bold text-primary hover:text-primary-hover text-left transition-colors"
                  >
                    Create API key in a new project
                  </button>
                </div>
              ) : (
                <div className="flex flex-col gap-3">
                  <label className="text-sm font-bold text-white">Project name</label>
                  <input 
                    type="text" 
                    placeholder="e.g., Sentiment Analysis Tool"
                    value={newProjectName}
                    onChange={(e) => setNewProjectName(e.target.value)}
                    className="w-full bg-black/50 border border-white/10 rounded-lg px-4 py-3 text-sm text-white focus:outline-none focus:border-primary/50 transition-colors"
                    autoFocus
                  />
                  
                  <button 
                    onClick={() => {
                        setIsCreatingNew(false);
                        setNewProjectName('');
                        setError(null);
                    }}
                    className="text-sm font-bold text-white/50 hover:text-white transition-colors text-left mt-2"
                  >
                    &larr; Back to select project
                  </button>
                </div>
              )}
            </div>

            {error && (
              <div className="px-6 pb-2 text-sm text-red-500 font-bold">
                {error}
              </div>
            )}

            <div className="p-6 bg-white/[0.02] border-t border-white/5 flex justify-end gap-3">
              <button 
                onClick={handleClose}
                className="px-5 py-2.5 rounded-lg text-sm font-bold text-white hover:bg-white/5 transition-colors"
              >
                Cancel
              </button>
              <button 
                onClick={handleGenerate}
                disabled={isCreatingNew ? !newProjectName : !selectedProject || isGenerating}
                className="flex items-center justify-center min-w-[140px] px-5 py-2.5 rounded-lg bg-primary hover:bg-primary-hover text-black text-sm font-bold transition-all disabled:opacity-50 disabled:cursor-not-allowed shadow-[0_0_15px_rgba(255,90,0,0.3)] disabled:shadow-none"
              >
                {isGenerating ? (
                   <div className="w-4 h-4 border-2 border-black/20 border-t-black rounded-full animate-spin"></div>
                ) : (
                   'Create API key'
                )}
              </button>
            </div>
          </>
        )}

        {step === 2 && (
          <>
            <div className="p-6 border-b border-white/5">
              <div className="w-10 h-10 rounded-full bg-success/20 flex items-center justify-center mb-4">
                 <svg className="w-5 h-5 text-success" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="4" d="M5 13l4 4L19 7" />
                 </svg>
              </div>
              <h2 className="text-xl font-bold text-white tracking-tight">API key generated</h2>
              <p className="text-sm text-foreground/60 mt-2">
                Save this key somewhere safe and accessible. For security reasons, <strong>you won't be able to view it again</strong>.
              </p>
            </div>

            <div className="p-6 flex flex-col gap-6">
               <div className="flex flex-col gap-2 relative">
                  <input 
                    type="text" 
                    value={generatedKey}
                    readOnly
                    className="w-full bg-black/80 border border-white/20 rounded-lg px-4 py-4 text-sm text-white font-mono focus:outline-none pr-28 selection:bg-primary/30"
                  />
                  <div className="absolute right-2 top-2 bottom-2 flex items-center">
                    <button 
                      onClick={handleCopy}
                      className={`px-3 py-1.5 rounded-md text-xs font-bold transition-all flex items-center gap-2 ${copied ? 'bg-success/20 text-success' : 'bg-white/10 hover:bg-white/20 text-white'}`}
                    >
                      {copied ? (
                        <>
                          <svg className="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                             <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="3" d="M5 13l4 4L19 7" />
                          </svg>
                          Copied
                        </>
                      ) : (
                        <>
                          <svg className="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                             <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-1M8 5a2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2" />
                          </svg>
                          Copy
                        </>
                      )}
                    </button>
                  </div>
               </div>
               
               <div className="p-4 bg-yellow-500/10 border border-yellow-500/20 rounded-xl flex items-start gap-3">
                  <svg className="w-5 h-5 text-yellow-500 shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                  </svg>
                  <p className="text-xs text-yellow-500 leading-relaxed">
                    Do not share your API key in publicly accessible areas such as GitHub, client-side code, and so forth.
                  </p>
               </div>
            </div>

            <div className="p-6 bg-white/[0.02] border-t border-white/5 flex justify-end">
              <button 
                onClick={handleClose}
                className="px-6 py-2.5 rounded-lg bg-white hover:bg-gray-200 text-black text-sm font-bold transition-colors"
              >
                Close
              </button>
            </div>
          </>
        )}
      </div>
    </div>
  );
}
