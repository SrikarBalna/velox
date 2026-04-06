'use client';

import React, { useRef, useEffect } from 'react';
import { useFocus } from './FocusContext';

export default function InteractiveStat({ id, title, children, modalContent }: { id: string, title: string, children: React.ReactNode, modalContent: React.ReactNode }) {
    const { focusedId, setFocusedId, activeModalId, setActiveModalId } = useFocus();
    const timeoutRef = useRef<NodeJS.Timeout | null>(null);

    const isFocused = focusedId === id;
    const isModalOpen = activeModalId === id;

    const handleMouseEnter = () => {
        if (activeModalId) return; // Don't focus if a modal is open globally
        timeoutRef.current = setTimeout(() => {
            setFocusedId(id);
        }, 500); // 500ms delay to trigger focus
    };

    const handleMouseLeave = () => {
        if (timeoutRef.current) clearTimeout(timeoutRef.current);
        if (isFocused) {
            setFocusedId(null);
        }
    };

    const handleClick = () => {
        if (isFocused) {
            setFocusedId(null);
            setActiveModalId(id);
        }
    };

    // Close modal on escape
    useEffect(() => {
        const handleKeyDown = (e: KeyboardEvent) => {
            if (e.key === 'Escape' && isModalOpen) {
                setActiveModalId(null);
            }
        };
        document.addEventListener('keydown', handleKeyDown);
        return () => document.removeEventListener('keydown', handleKeyDown);
    }, [isModalOpen, setActiveModalId]);

    return (
        <>
            <div 
                className={`transition-all duration-500 ease-out ${isFocused ? 'scale-105 z-50 shadow-[0_0_40px_rgba(255,90,0,0.15)] ring-1 ring-primary/30' : 'z-10'} relative rounded-2xl`}
                onMouseEnter={handleMouseEnter}
                onMouseLeave={handleMouseLeave}
                onClick={handleClick}
                style={{ cursor: isFocused ? 'pointer' : 'default' }}
            >
                {/* Visual indicator when focused to let them know it's clickable */}
                {isFocused && (
                    <div className="absolute inset-x-0 -top-8 flex justify-center animate-bounce z-50 pointer-events-none">
                        <span className="bg-primary text-black text-[10px] font-black uppercase tracking-wider px-3 py-1 rounded-full shadow-lg">Click to Expand</span>
                    </div>
                )}
                {children}
            </div>

            {/* Modal Overlay */}
            {isModalOpen && (
                <div className="fixed inset-0 z-[100] flex items-center justify-center p-4">
                    {/* Backdrop */}
                    <div 
                        className="absolute inset-0 bg-black/80 backdrop-blur-md" 
                        onClick={() => setActiveModalId(null)}
                    ></div>
                    
                    {/* Modal Content Window */}
                    <div className="relative bg-surface border border-white/10 rounded-2xl shadow-2xl p-8 max-w-3xl w-full z-[101] overflow-hidden">
                        <div className="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-primary to-orange-400"></div>
                        
                        <button 
                            className="absolute top-4 right-4 text-white/40 hover:text-white transition-colors bg-white/5 hover:bg-white/10 p-2 rounded-lg"
                            onClick={() => setActiveModalId(null)}
                        >
                            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M6 18L18 6M6 6l12 12" /></svg>
                        </button>
                        
                        <div className="mb-8">
                            <h2 className="text-2xl font-bold text-white">{title}</h2>
                            <p className="text-sm text-white/50 font-mono mt-1">Detailed statistical view</p>
                        </div>
                        
                        <div className="animate-in fade-in slide-in-from-bottom-4 duration-500">
                            {modalContent}
                        </div>
                    </div>
                </div>
            )}
        </>
    );
}
