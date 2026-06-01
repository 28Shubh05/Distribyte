import { useEffect, useState } from "react";
import api from "../services/api";

function Dashboard() {
  const [files, setFiles] = useState([]);

  useEffect(() => {
    fetchFiles();
  }, []);

  const fetchFiles = async () => {
    try {
      const response = await api.get("/files");
      setFiles(response.data.files);
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <div>
      <h2>Stored Files</h2>

      {files.map((file) => (
        <div key={file.id}>
          <p>{file.original_name}</p>
        </div>
      ))}
    </div>
  );
}

export default Dashboard;