import { useEffect, useState } from "react";
import api from "../services/api";

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

      alert("File restored successfully");

      fetchDeletedFiles();

    } catch (error) {

      alert(
        error.response?.data?.error ||
        "Restore failed"
      );
    }
  };

  return (
    <div>

      <h2>Deleted Files</h2>

      {files.length === 0 ? (
        <p>No deleted files</p>
      ) : (
        files.map((file) => (
          <div key={file.id}>

            <span>{file.original_name}</span>

            <button
              onClick={() =>
                handleRestore(file.id)
              }
            >
              Restore
            </button>

          </div>
        ))
      )}

    </div>
  );
}

export default DeletedFiles;