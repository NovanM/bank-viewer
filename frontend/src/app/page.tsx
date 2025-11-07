
"use client";

import React, { useEffect, useState } from "react";
import {
  useQuery,
  useMutation,
  useQueryClient,
  keepPreviousData, 
} from "@tanstack/react-query";

import type {
  APIResponse,
  BalanceData,
  IssuesData,
  IssuesQueryParams,
} from "@/types/api.types";

import { FileUploader } from "@/components/FileUploader";
import { BalanceView } from "@/components/BalanceView";
import { Datatable } from "@/components/Datatables"; 
import { Toast } from "@/components/Toast";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:9090";

const fetchBalance = async (): Promise<BalanceData> => {
  const res = await fetch(`${API_URL}/balance`);
  if (!res.ok) throw new Error("Network response was not ok");
  
  const json: APIResponse<BalanceData> = await res.json();
  if (!json.status || !json.data) throw new Error(json.message);
  
  return json.data;
};

const fetchIssues = async (params: IssuesQueryParams): Promise<IssuesData> => {
  const query = `?page=${params.page}&limit=${params.limit}&sort_by=${params.sortBy}&sort_dir=${params.sortDir}`;
  const res = await fetch(`${API_URL}/issues${query}`);
  if (!res.ok) throw new Error("Network response was not ok");

  const json: APIResponse<IssuesData> = await res.json();
  if (!json.status || !json.data) throw new Error(json.message);
  
  return json.data;
};

const uploadFile = async (file: File): Promise<string> => {
  const formData = new FormData();
  formData.append("file", file);

  const res = await fetch(`${API_URL}/upload`, {
    method: "POST",
    body: formData,
  });
  
  const json: APIResponse<null> = await res.json();
  if (!res.ok || !json.status) throw new Error(json.message);

  return json.message;
};


export default function Home() {
  const queryClient = useQueryClient();

  const [queryParams, setQueryParams] = useState<IssuesQueryParams>({
    page: 1,
    limit: 10,
    sortBy: "timestamp",
    sortDir: "desc",
  });

  const [notification, setNotification] = useState<{
    type: "success" | "error";
    message: string;
  } | null>(null);

  useEffect(() => {
    if (!notification) return;
    const timeout = window.setTimeout(() => setNotification(null), 4000);
    return () => window.clearTimeout(timeout);
  }, [notification]);


  const balanceQuery = useQuery({
    queryKey: ["balance"], 
    queryFn: fetchBalance,
  });

  const issuesQuery = useQuery({
    queryKey: ["issues", queryParams], 
    queryFn: () => fetchIssues(queryParams),
    placeholderData: keepPreviousData,
  });

  const uploadMutation = useMutation({
    mutationFn: uploadFile,
    onSuccess: (successMessage) => {
      setNotification({ type: "success", message: successMessage });
      setQueryParams(prev => ({ ...prev, page: 1 }));
      queryClient.invalidateQueries({ queryKey: ["balance"] });
      queryClient.invalidateQueries({ queryKey: ["issues"] });
    },
    onError: (err: Error) => {
      setNotification({ type: "error", message: `Upload failed: ${err.message}` });
    },
  });

  return (
    <main className="app-shell">
      <header className="page-header">
        <div>
          <h1>Bank Statement Viewer</h1>
          <p className="page-subtitle">Monitor balances and review flagged transactions in real time.</p>
        </div>
      </header>

      {notification && (
        <div className="toast-container">
          <Toast
            type={notification.type}
            message={notification.message}
            onClose={() => setNotification(null)}
          />
        </div>
      )}

      <div className="overview-grid">
        <BalanceView
          data={balanceQuery.data}
          isLoading={balanceQuery.isLoading}
          error={balanceQuery.error} 
        />

        <FileUploader
          isPending={uploadMutation.isPending}
          onUpload={(file) => uploadMutation.mutate(file)}
        />
      </div>

      <Datatable
        data={issuesQuery.data}
        isLoading={issuesQuery.isLoading}
        error={issuesQuery.error}
        queryParams={queryParams}
        setQueryParams={setQueryParams}
      />
    </main>
  );
}