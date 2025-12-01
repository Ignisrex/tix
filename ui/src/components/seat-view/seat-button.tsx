"use client";

import { useState, useRef } from "react";
import type { SeatButtonProps } from "./types";

export function SeatButton({ ticket, isAvailable, isSelected, colorConfig, onClick }: SeatButtonProps) {
  const [ripples, setRipples] = useState<Array<{ id: number; x: number; y: number }>>([]);
  const buttonRef = useRef<HTMLButtonElement>(null);
  const rippleIdRef = useRef(0);

  const handleClick = (e: React.MouseEvent<HTMLButtonElement>) => {
    if (!isAvailable || !buttonRef.current) return;

    const button = buttonRef.current;
    const rect = button.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const y = e.clientY - rect.top;

    const newRipple = {
      id: rippleIdRef.current++,
      x,
      y,
    };

    setRipples((prev) => [...prev, newRipple]);

    // Remove ripple after animation completes
    setTimeout(() => {
      setRipples((prev) => prev.filter((r) => r.id !== newRipple.id));
    }, 600);

    onClick();
  };

  return (
    <button
      ref={buttonRef}
      onClick={handleClick}
      disabled={!isAvailable}
      className={`
        relative aspect-square rounded-md border-2 transition-all overflow-hidden
        ${isAvailable 
          ? `${colorConfig.bgColor} ${colorConfig.borderColor} hover:scale-105 hover:shadow-md cursor-pointer ${
              isSelected ? 'ring-2 ring-offset-2 ring-indigo-500 ring-offset-background' : ''
            }` 
          : 'bg-gray-300 border-gray-400 opacity-50 cursor-not-allowed'
        }
        flex items-center justify-center text-xs font-medium
        ${isAvailable ? colorConfig.color : 'text-gray-500'}
      `}
      title={isAvailable ? `Seat ${ticket.id.slice(0, 8)} - Available` : `Seat ${ticket.id.slice(0, 8)} - Sold`}
    >
      {ticket.id.slice(0, 4)}
      
      {/* Ripple effects */}
      {ripples.map((ripple) => (
        <span
          key={ripple.id}
          className="absolute rounded-full bg-white/40 pointer-events-none animate-ripple"
          style={{
            left: `${ripple.x}px`,
            top: `${ripple.y}px`,
            transform: 'translate(-50%, -50%)',
          }}
        />
      ))}
    </button>
  );
}

