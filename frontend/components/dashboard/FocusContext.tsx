'use client';
import { createContext, useContext, useState, ReactNode } from 'react';

type FocusState = {
  focusedId: string | null;
  setFocusedId: (id: string | null) => void;
  activeModalId: string | null;
  setActiveModalId: (id: string | null) => void;
};

export const FocusContext = createContext<FocusState | null>(null);

export function FocusProvider({ children }: { children: ReactNode }) {
  const [focusedId, setFocusedId] = useState<string | null>(null);
  const [activeModalId, setActiveModalId] = useState<string | null>(null);

  return (
    <FocusContext.Provider value={{ focusedId, setFocusedId, activeModalId, setActiveModalId }}>
      {children}
    </FocusContext.Provider>
  );
}

export const useFocus = () => {
  const ctx = useContext(FocusContext);
  if (!ctx) throw new Error("useFocus must be used within FocusProvider");
  return ctx;
}
