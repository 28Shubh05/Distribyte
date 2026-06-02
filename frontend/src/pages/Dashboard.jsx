import { useEffect, useState } from "react";
import api from "../services/api";
import UploadForm from "../components/UploadForm";
import toast from "react-hot-toast";

function Dashboard() {
  const [files, setFiles] = useState([]);
  const [searchTerm, setSearchTerm] = useState("");
  const [loading, setLoading] = useState(true);

  const handleDelete = async (id) => {
    try {
      await api.delete(`/files/${id}`);

      toast.success("File deleted");

      fetchFiles();
    } catch (error) {
        toast.error(
          error.response?.data?.error ||
          "Delete failed"
        );
    }
  };

  const fetchFiles = async () => {
    setLoading(true);
    try {
      const response = await api.get("/files");
      setFiles(response.data.files);
    } catch (error) {
      console.error(error);
    }finally{
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchFiles();
  }, []);

  const totalStorage = files.reduce(
    (sum, file) => sum + file.size,
    0
  );

  const filteredFiles = files.filter((file) =>
    file.original_name
      .toLowerCase()
      .includes(searchTerm.toLowerCase())
  );

  const getFileIcon = (filename) => {

  const ext =
    filename.split(".").pop().toLowerCase();

    switch (ext) {

      case "pdf":
        return "📕";

      case "png":
      case "jpg":
      case "jpeg":
        return "🖼️";

      case "zip":
        return "🗜️";

      case "txt":
        return "📄";

      default:
        return "📁";
    }
  };

  if (loading) {
    return (
      <div className="text-center py-10">

        <div className="text-2xl">
          Loading...
        </div>

      </div>
      );
  }

  return (
    <div className="space-y-6">

      {/* Header */}
      <div>
        <h1 className="text-4xl font-bold text-slate-800">
          Dashboard
        </h1>

        <p className="text-gray-500 mt-1">
          Manage your files stored in Distribyte
        </p>
      </div>

      {/* Stats */}
      <div className="grid md:grid-cols-3 gap-4">

        <div className="bg-white rounded-xl shadow p-5">
          <h3 className="text-gray-500">
            Total Files
          </h3>

          <p className="text-3xl font-bold mt-2">
            {files.length}
          </p>
        </div>

        <div className="bg-white rounded-xl shadow p-5">
          <h3 className="text-gray-500">
            Storage Used
          </h3>

          <p className="text-3xl font-bold mt-2">
            {(totalStorage / 1024).toFixed(2)} KB
          </p>
        </div>

        <div className="bg-white rounded-xl shadow p-5">
          <h3 className="text-gray-500">
            System Status
          </h3>

          <p className="text-green-600 font-bold mt-2">
            Healthy
          </p>
        </div>

      </div>

      <div className="bg-white rounded-xl shadow p-5">

  <div className="flex justify-between mb-2">

    <span className="font-medium">
      Storage Usage
    </span>

    <span>
      {(totalStorage / 1024 / 1024).toFixed(2)} MB
    </span>

  </div>

  <div className="w-full bg-gray-200 rounded-full h-4">

    <div
      className="bg-blue-600 h-4 rounded-full"
      style={{
        width: `${Math.min(
          (totalStorage /
            (100 * 1024 * 1024)) *
            100,
          100
        )}%`,
      }}
    />

      </div>

    </div>

      {/* Upload */}
      <UploadForm onUploadSuccess={fetchFiles} />

      {/* Files Section */}
      <div>

        <div className="bg-white rounded-xl shadow p-4">

          <input
            type="text"
            placeholder="Search files..."
            value={searchTerm}
            onChange={(e) =>
              setSearchTerm(e.target.value)
            }
            className="w-full border rounded-lg p-3"
          />

      </div>

        <h2 className="text-2xl font-bold mb-4">
          Stored Files
        </h2>

        {files.length === 0 ? (
          <div className="bg-white rounded-xl shadow p-6 text-center">
            No files uploaded yet
          </div>
        ) : (
          <div className="grid gap-4">

            {filteredFiles.map((file) => (

              <div
                key={file.id}
                className="bg-white rounded-xl shadow p-5 flex justify-between items-center hover:shadow-lg transition"
              >

                <div>

                  <h3 className="font-semibold text-lg">
                    {getFileIcon(file.original_name)}
                    {" "}
                    {file.original_name}
                  </h3>

                  <p className="text-sm text-gray-500">
                    {(file.size / 1024).toFixed(2)} KB
                  </p>

                </div>

                <div className="flex gap-3">

                  <button
                    onClick={() =>
                      window.open(
                        `http://localhost:8080/download/${file.id}`,
                        "_blank"
                      )
                    }
                    className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700"
                  >
                    Download
                  </button>

                  <button
                    onClick={() =>
                      handleDelete(file.id)
                    }
                    className="bg-red-600 text-white px-4 py-2 rounded-lg hover:bg-red-700"
                  >
                    Delete
                  </button>

                </div>

              </div>

            ))}

          </div>
        )}

      </div>

    </div>
  );
}

export default Dashboard;