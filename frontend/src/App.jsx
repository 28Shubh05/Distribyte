import { useEffect, useState } from "react";

import Dashboard from "./pages/Dashboard";
import DeletedFiles from "./pages/DeletedFiles";
import ThemeToggle from "./components/ThemeToggle";

function App() {
  const [page, setPage] = useState("dashboard");

  const [dark, setDark] = useState(() => {
    return localStorage.getItem("theme") === "dark";
  });

  useEffect(() => {
    document.documentElement.classList.toggle("dark", dark);
    localStorage.setItem("theme", dark ? "dark" : "light");
  }, [dark]);

  return (
    <div className="min-h-screen bg-slate-100 dark:bg-slate-950 transition-colors duration-300">

      <nav className="bg-white dark:bg-slate-900 shadow-md border-b border-slate-200 dark:border-slate-800">

        <div className="max-w-7xl mx-auto px-6 py-4 flex justify-between items-center">

          <div>
            <h1 className="text-3xl font-bold bg-gradient-to-r from-blue-600 to-cyan-500 bg-clip-text text-transparent">
              Distribyte
            </h1>

            <p className="text-sm text-slate-500 dark:text-slate-400">
              Distributed Object Storage
            </p>
          </div>

          <div className="flex items-center gap-3">

            <button
              onClick={() => setPage("dashboard")}
              className={`px-4 py-2 rounded-xl transition-all ${
                page === "dashboard"
                  ? "bg-blue-600 text-white shadow-lg"
                  : "bg-slate-200 dark:bg-slate-800 dark:text-slate-300"
              }`}
            >
              Dashboard
            </button>

            <button
              onClick={() => setPage("deleted")}
              className={`px-4 py-2 rounded-xl transition-all ${
                page === "deleted"
                  ? "bg-red-600 text-white shadow-lg"
                  : "bg-slate-200 dark:bg-slate-800 dark:text-slate-300"
              }`}
            >
              Deleted Files
            </button>

            <ThemeToggle
              dark={dark}
              setDark={setDark}
            />

          </div>

        </div>

      </nav>

      <main className="max-w-7xl mx-auto p-6">
        {page === "dashboard"
          ? <Dashboard />
          : <DeletedFiles />}
      </main>

      <footer className="text-center py-6 text-slate-500 dark:text-slate-400">
        Distribyte • Built with Go • PostgreSQL • React
      </footer>

    </div>
  );
}

export default App;