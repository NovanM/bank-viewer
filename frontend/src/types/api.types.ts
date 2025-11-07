
export interface Transaction {
  timestamp: string; 
  name: string;
  type: "DEBIT" | "CREDIT";
  amount: number; 
  status: "SUCCESS" | "FAILED" | "PENDING";
  description: string;
}

export interface APIResponse<T> {
  status: boolean;
  message: string;
  data?: T;
}

export interface BalanceData {
  total_balance: number;
}

export interface PaginationMetadata {
  current_page: number;
  page_size: number;
  total_items: number;
  total_pages: number;
}

export interface IssuesData {
  transactions: Transaction[];
  metadata: PaginationMetadata;
}

export interface IssuesQueryParams {
  page: number;
  limit: number;
  sortBy: string;
  sortDir: string;
}