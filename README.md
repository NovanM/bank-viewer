
# Bank Statement Viewer

**Solution results video link:**

This is a full-stack application that allows users to upload a bank statement CSV file, view their calculated end balance, and inspect a list of problematic transactions.

This project was built as a take-home test to demonstrate best practices in backend (Go) and frontend (Next.js) architecture, including *concurrency*, *streaming*, *state management*, and CI.

## Key Features

  * **CSV Upload:** Accepts CSV file uploads  with the format `timestamp, name, type, amount, status, description`  and parses them using efficient streaming.
  * **Balance Calculation:** Displays the final total balance, calculated only from "SUCCESS" transactions (total credits minus total debits).
  * **Issue Table:** Displays a list of "PENDING" and "FAILED" transactions in a table.
  * **Pagination & Sorting:** The issue table supports server-side pagination and sorting (e.g., `?page=2&sort_by=amount`).
  * **Status Styling:** Provides clear visual styling for `PENDING` (warning/yellow)  and `FAILED` (red)  statuses.

## Tech Stack

  * **Backend:** Go (Golang) (REST API, pure `net/http`)
  * **Frontend:** Next.js (App Router) & React Query (TanStack Query)
  * **DevOps:** Docker, Docker Compose, GitHub Actions (CI)
  * **Other:** TypeScript, Pure CSS (No UI Libraries)

-----

## Project Structure (Monorepo)

This project is organized as a *monorepo* to maintain a clear separation of concerns:

  * **`/backend`**: Contains the standalone Go (API) application.
  * **`/frontend`**: Contains the standalone Next.js (UI) application.
  * **`/docker-compose.yml`**: Orchestrates both services for local development.
  * **`/.github/workflows/`**: Contains separate CI pipelines for the backend and frontend.

-----

## Setup & Installation Instructions 

There are two methods to run this project. The Docker Compose method is highly recommended.

### Method 1: Docker Compose (Recommended)

This is the easiest way to run both the backend and frontend services simultaneously.

1.  Ensure you have Docker and Docker Compose installed.

2.  Clone this repository.

3.  From the root directory, run the following command:

    ```bash
    # Build and start both containers (backend & frontend)
    docker-compose up --build
    ```

4.  Access the applications:

      * **Frontend (UI):** `http://localhost:3000`
      * **Backend (API):** `http://localhost:9090` (or your configured port)

### Method 2: Manual (Local)

You will need to run two separate terminals.

#### 1\. Backend (Go)

```bash
# Navigate to the backend directory
cd backend

# Run the server
# (The server will run on the port defined in main.go, e.g., :9090)
go run main.go
```

#### 2\. Frontend (Next.js)

```bash
# Open a NEW terminal
# Navigate to the frontend directory
cd frontend

# Install dependencies
npm install

# Run the Next.js development server
npm run dev
```

The frontend application will be available at `http://localhost:3000`.

-----

## Architecture Decisions

This section explains the technical decisions made, focusing on a clean, scalable, and non-"vibe coding" architecture.

### Backend (Go)

  * **Clean Architecture :** We implemented a `handler` -\> \`service\` -\> \`repository\` separation. This makes the code highly testable (business logic in the \`service\` is isolated) and maintainable.
  * **Streaming Upload & Validation :** To handle large CSV files without consuming excessive memory, the handler parses the request as a stream. We also implemented a "Gatekeeper" (`http.MaxBytesReader` at 20MB) to reject requests that are too large *before* memory is consumed, as a DoS protection.
  * **"Free Rollback" Error Handling:** Our service design parses the *entire* file *first*. Only if the parsing is 100% successful is the new data `Store`-d in the repository . This prevents our in-memory data from being left in a corrupted or partial state if parsing fails midway.
  * **Concurrency & Data Consistency:** We handle concurrent uploads using a "Last Writer Wins" strategy. The `repository.Store` method is protected by a `sync.RWMutex`, ensuring only one write operation can occur at a time, preventing data corruption.
  * **Backend-Driven Pagination :** Instead of sending thousands of issues to the frontend, we implemented *pagination* and *sorting* on the server-side (`GET /issues?page=...`). This is scalable and keeps the frontend lightweight.

### Frontend (Next.js)

  * **React Query (TanStack Query):** We deliberately avoided manual `useState` + `useEffect` for API data. We use React Query to manage "Server State." This gives us caching, automatic refetching, and mutation handling (e.g., a successful `POST /upload` automatically triggers a refetch of `GET /balance` and `GET /issues`).
  * **Reusable Components :** The UI is broken down into declarative components (`FileUploader`, \`BalanceView\`, \`Datatable\`) that receive *props*, aligning with a clean architecture philosophy.
  * \*\*Pure CSS :** As per the constraints, all styling (including `FAILED`/`PENDING` statuses ) was achieved using Pure CSS in `globals.css` without any external UI libraries.

### DevOps (Extra Features) 

  * **Dockerfile Multi-stage :** We use multi-stage builds for both the `backend/Dockerfile` and \`frontend/Dockerfile\` to ensure the final production images are minimal and secure.
  * **Monorepo CI :** We created two separate GitHub Actions workflows. Both are triggered by `paths` so that a change in \`frontend/\` does not run the \`backend/\` CI, saving resources.
      * **`backend.yml`:** Runs a Go linter (`golangci-lint`), vulnerability scanner (`govulncheck`), and unit tests with code coverage (`go test -race -coverprofile=...`).
      * **`frontend.yml`:** Runs the ESLint linter and `npm run build`.