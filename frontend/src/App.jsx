import { useEffect, useState } from "react";

import Dashboard from "./pages/Dashboard";
import DeletedFiles from "./pages/DeletedFiles";
import ThemeToggle from "./components/ThemeToggle";

import Login from "./pages/Login";
import Register from "./pages/Register";

function App() {
  const [page, setPage] = useState("dashboard");

  const [showRegister, setShowRegister] =
    useState(false);

  const [loggedIn, setLoggedIn] =
    useState(
      !!localStorage.getItem("token")
    );

  const [dark, setDark] = useState(() => {

    const savedTheme =
      localStorage.getItem("theme");

    return savedTheme === "dark";
  });

  useEffect(() => {

    document.documentElement.classList.toggle(
      "dark",
      dark
    );

    localStorage.setItem(
      "theme",
      dark ? "dark" : "light"
    );

  }, [dark]);

  const handleLogout = () => {

    localStorage.removeItem("token");

    setLoggedIn(false);

    setPage("dashboard");
  };

  if (!loggedIn) {

    return (
      <div className="min-h-screen bg-slate-100 dark:bg-slate-950 transition-colors duration-300">

        <div className="flex justify-end p-4">
          <ThemeToggle
            dark={dark}
            setDark={setDark}
          />
        </div>

        <div className="max-w-md mx-auto">

          {showRegister ? (
            <>
              <Register />

              <p className="text-center mt-4 text-slate-600 dark:text-slate-400">
                Already have an account?{" "}
                <button
                  onClick={() =>
                    setShowRegister(false)
                  }
                  className="text-blue-600 font-semibold"
                >
                  Login
                </button>
              </p>
            </>
          ) : (
            <>
              <Login
                onLogin={() =>
                  setLoggedIn(true)
                }
              />

              <p className="text-center mt-4 text-slate-600 dark:text-slate-400">
                Don't have an account?{" "}
                <button
                  onClick={() =>
                    setShowRegister(true)
                  }
                  className="text-blue-600 font-semibold"
                >
                  Register
                </button>
              </p>
            </>
          )}

        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-slate-100 dark:bg-slate-950 transition-colors duration-300">

      {/* Navbar */}
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
              onClick={() =>
                setPage("dashboard")
              }
              className={`px-4 py-2 rounded-xl transition-all ${
                page === "dashboard"
                  ? "bg-blue-600 text-white shadow-lg"
                  : "bg-slate-200 dark:bg-slate-800 dark:text-slate-300"
              }`}
            >
              Dashboard
            </button>

            <button
              onClick={() =>
                setPage("deleted")
              }
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

            <button
              onClick={handleLogout}
              className="px-4 py-2 rounded-xl bg-red-600 hover:bg-red-700 text-white transition"
            >
              Logout
            </button>

          </div>

        </div>

      </nav>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto p-6">

        {page === "dashboard" ? (
          <Dashboard />
        ) : (
          <DeletedFiles />
        )}

      </main>

      {/* Footer */}
      <footer className="text-center py-6 text-slate-500 dark:text-slate-400 border-t border-slate-200 dark:border-slate-800 mt-10">
        Distribyte • Built with Go • PostgreSQL • Redis • React
      </footer>

    </div>
  );
}

export default App;