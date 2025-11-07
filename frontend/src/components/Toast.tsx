"use client";

import React from "react";

type ToastType = "success" | "error" | "info";

interface ToastProps {
  message: string;
  type?: ToastType;
  onClose: () => void;
}

export function Toast({ message, type = "info", onClose }: ToastProps) {
  return (
    <div className={`toast toast-${type}`} role="status">
      <span className="toast-message">{message}</span>
      <button
        type="button"
        className="toast-dismiss"
        onClick={onClose}
        aria-label="Close notification"
      >
        x
      </button>
    </div>
  );
}
