import { useState } from "react";

import Dashboard from "./pages/Dashboard";
import DeletedFiles from "./pages/DeletedFiles";

function App() {

  const [page, setPage] = useState("dashboard");

  return (
    <div>

      <h1>Distribyte</h1>

      <button
        onClick={() =>
          setPage("dashboard")
        }
      >
        Dashboard
      </button>

      <button
        onClick={() =>
          setPage("deleted")
        }
      >
        Deleted Files
      </button>

      <hr />

      {page === "dashboard" ? (
        <Dashboard />
      ) : (
        <DeletedFiles />
      )}

    </div>
  );
}

export default App;