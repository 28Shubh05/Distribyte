import { useEffect, useState } from "react";
import api from "../services/api";
import toast from "react-hot-toast";

function DeletedFiles() {
  const [files, setFiles] = useState([]);

  const fetchDeletedFiles = async () => {
    try {
      const response = await api.get("/deleted-files");
      setFiles(response.data.files);
    } catch (error) {
      console.error(error);
    }
  };

  useEffect(() => {
    fetchDeletedFiles();
  }, []);

  const handleRestore = async (id) => {
    try {
      await api.post(`/restore/${id}`);

      toast.success("File restored");

      fetchDeletedFiles();
    } catch (error) {
      toast.error(
        error.response?.data?.error ||
        "Restore failed"
      );
    }
  };

  const totalDeletedStorage = files.reduce(
    (sum, file) => sum + file.size,
    0
  );

  return (
    <div className="space-y-6">

      {/* Header */}
      <div>
        <h1 className="text-4xl font-bold text-slate-800">
          Deleted Files
        </h1>

        <p className="text-gray-500 mt-1">
          Restore previously deleted files
        </p>
      </div>

      {/* Stats */}
      <div className="grid md:grid-cols-2 gap-4">

        <div className="bg-white rounded-xl shadow p-5">
          <h3 className="text-gray-500">
            Deleted Files
          </h3>

          <p className="text-3xl font-bold mt-2">
            {files.length}
          </p>
        </div>

        <div className="bg-white rounded-xl shadow p-5">
          <h3 className="text-gray-500">
            Recoverable Storage
          </h3>

          <p className="text-3xl font-bold mt-2">
            {(totalDeletedStorage / 1024).toFixed(2)} KB
          </p>
        </div>

      </div>

      {/* Deleted Files List */}
      <div>

        <h2 className="text-2xl font-bold mb-4">
          Deleted Items
        </h2>

        {files.length === 0 ? (
          <div className="bg-white rounded-xl shadow p-6 text-center">
            No deleted files
          </div>
        ) : (
          <div className="grid gap-4">

            {files.map((file) => (

              <div
                key={file.id}
                className="bg-white rounded-xl shadow p-5 flex justify-between items-center hover:shadow-lg transition"
              >

                <div>

                  <h3 className="font-semibold text-lg">
                    🗑️ {file.original_name}
                  </h3>

                  <p className="text-sm text-gray-500">
                    {(file.size / 1024).toFixed(2)} KB
                  </p>

                  {file.deleted_at && (
                    <p className="text-xs text-red-500 mt-1">
                      Deleted: {new Date(file.deleted_at).toLocaleString()}
                    </p>
                  )}

                </div>

                <button
                  onClick={() => handleRestore(file.id)}
                  className="bg-green-600 text-white px-4 py-2 rounded-lg hover:bg-green-700"
                >
                  Restore
                </button>

              </div>

            ))}

          </div>
        )}

      </div>

    </div>
  );
}

export default DeletedFiles;