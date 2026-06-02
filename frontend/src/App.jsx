import { useState } from "react";

import Dashboard from "./pages/Dashboard";
import DeletedFiles from "./pages/DeletedFiles";

function App() {
  const [page, setPage] = useState("dashboard");

  return (
    <div className="min-h-screen bg-slate-100">

      {/* Navbar */}
      <nav className="bg-white shadow-md">

        <div className="max-w-6xl mx-auto px-6 py-4 flex justify-between items-center">

            <div>
              <h1 className="text-3xl font-bold text-blue-600">
                Distribyte
              </h1>

              <p className="text-sm text-gray-500">
                Distributed Object Storage
              </p>
          </div>

          <div className="flex gap-4">

            <button
              onClick={() => setPage("dashboard")}
              className={`px-4 py-2 rounded-lg ${
                page === "dashboard"
                  ? "bg-blue-600 text-white"
                  : "bg-gray-200"
              }`}
            >
              Dashboard
            </button>

            <button
              onClick={() => setPage("deleted")}
              className={`px-4 py-2 rounded-lg ${
                page === "deleted"
                  ? "bg-red-600 text-white"
                  : "bg-gray-200"
              }`}
            >
              Deleted Files
            </button>

          </div>

        </div>

      </nav>

      <main className="max-w-6xl mx-auto p-6">

        {page === "dashboard"
          ? <Dashboard />
          : <DeletedFiles />
        }

      </main>

    </div>
  );
}

export default App;