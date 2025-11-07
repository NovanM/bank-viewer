
"use client";

import React from "react";
import { IssuesData, IssuesQueryParams, Transaction } from "@/types/api.types";

interface Props {
  data: IssuesData | undefined;
  isLoading: boolean;
  error: Error | null;
  queryParams: IssuesQueryParams;
  setQueryParams: React.Dispatch<React.SetStateAction<IssuesQueryParams>>;
}

export function Datatable({ 
  data,
  isLoading,
  error,
  queryParams,
  setQueryParams,
}: Props) {

  const getStatusClass = (status: Transaction["status"]) => {
    if (status === "FAILED") return "status-failed";
    if (status === "PENDING") return "status-pending";
    return "";
  };
  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString("id-ID");
  }
  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(amount);
  };
  const handleSort = (newSortBy: string) => {
    const isAsc = queryParams.sortBy === newSortBy && queryParams.sortDir === "asc";
    setQueryParams({
      ...queryParams,
      sortBy: newSortBy,
      sortDir: isAsc ? "desc" : "asc",
      page: 1, 
    });
  };
  const handlePageChange = (newPage: number) => {
    setQueryParams({ ...queryParams, page: newPage });
  };
  
  const metadata = data?.metadata;
  const renderContent = () => {
    if (isLoading) {
      return <p className="loading-text">Loading issues...</p>;
    }
    
    if (error) {
      return <p className="error-text">Error: {error.message}</p>;
    }
    
    if (!data || data.metadata.total_items === 0) {
      return (
        <div className="no-data-message">
          <p>No issues found. All transactions are successful.</p>
        </div>
      );
    }
    
    return (
      <>
        <div className="table-wrapper">
          <table className="issues-table">
            <thead>
              <tr>
                <th onClick={() => handleSort("timestamp")}>Timestamp</th>
                <th onClick={() => handleSort("name")}>Name</th>
                <th onClick={() => handleSort("amount")}>Amount</th>
                <th>Status</th>
                <th>Description</th>
              </tr>
            </thead>
            <tbody>
              {data.transactions.map((tx) => (
                <tr key={tx.timestamp + tx.name}>
                  <td>{formatDate(tx.timestamp)}</td>
                  <td>{tx.name}</td>
                  <td>{formatCurrency(tx.amount)} ({tx.type})</td>
                  <td>
                    <span className={`status ${getStatusClass(tx.status)}`}>
                      {tx.status}
                    </span>
                  </td>
                  <td>{tx.description}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
        
        {metadata && (
          <div className="paginator">
            <button 
              className="button button--ghost"
              onClick={() => handlePageChange(metadata.current_page - 1)}
              disabled={metadata.current_page <= 1}
            >
              Previous
            </button>
            <span>
              Page {metadata.current_page} of {metadata.total_pages} 
              (Total: {metadata.total_items} issues)
            </span>
            <button
              className="button button--ghost"
              onClick={() => handlePageChange(metadata.current_page + 1)}
              disabled={metadata.current_page >= metadata.total_pages}
            >
              Next
            </button>
          </div>
        )}
      </>
    );
  }

  return (
    <section className="component-section">
      <h2>Non-Successful Transactions</h2>
      <p className="component-description">Review transactions that failed or are pending and drill into their details.</p>
      {renderContent()}
    </section>
  );
}