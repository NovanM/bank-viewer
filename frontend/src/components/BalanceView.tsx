
"use client";
import React from "react";
import { BalanceData } from "@/types/api.types";

interface Props {
  data: BalanceData | undefined;
  isLoading: boolean;
  error: Error | null;
}

const formatCurrency = (amount: number) => {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(amount);
};


export function BalanceView({ data, isLoading, error }: Props) {
  return (
    <section className="component-section">
      <h2>End Balance</h2>
      <p className="component-description">Latest reconciled closing balance for the selected statement period.</p>
      {isLoading && <p className="loading-text">Loading balance...</p>}
      {error && <p className="error-text">Error: {error.message}</p>}
      {data && (
        <div 
          className={`balance-view ${data.total_balance < 0 ? 'negative' : ''}`}
        >
          {formatCurrency(data.total_balance)}
        </div>
      )}
    </section>
  );
}