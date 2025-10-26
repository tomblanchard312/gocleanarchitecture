import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "http://localhost:8080"
export const WS_URL = import.meta.env.VITE_WS_URL || "ws://localhost:8080/ws"

