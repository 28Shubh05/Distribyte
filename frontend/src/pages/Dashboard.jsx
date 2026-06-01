import { useEffect, useState } from "react";
import api from "../services/api";
import UploadForm from "../components/UploadForm";

function Dashboard() {
  const [files, setFiles] = useState([]);

  const handleDelete = async (id) => {
        try {
            await api.delete(`/files/${id}`);

            alert("File deleted successfully");

            fetchFiles();
        } catch (error) {
            alert(
            error.response?.data?.error ||
            "Delete failed"
            );
        }
    };

  const fetchFiles = async () => {
    try {
      const response = await api.get("/files");
      setFiles(response.data.files);
    } catch (error) {
      console.error(error);
    }
  };

  useEffect(() => {
    fetchFiles();
  }, []);

  return (
    <div>

      <UploadForm
        onUploadSuccess={fetchFiles}
      />

      <h2>Stored Files</h2>

      {files.map((file) => (
        <div key={file.id}>
                <span>{file.original_name}</span>

                <button
                onClick={() =>
                    window.open(
                    `http://localhost:8080/download/${file.id}`,
                    "_blank"
                    )
                }
                >
                Download
                </button>

                <button
                onClick={() => handleDelete(file.id)}
                >
                Delete
                </button>
            </div>
        ))}

    </div>
  );
}

export default Dashboard;