import { useState } from "react";
import api from "../services/api";

function UploadForm({ onUploadSuccess }) {
  const [selectedFile, setSelectedFile] = useState(null);

  const handleUpload = async () => {
    if (!selectedFile) {
      alert("Select a file first");
      return;
    }

    const formData = new FormData();

    formData.append("file", selectedFile);

    try {
      await api.post("/upload", formData, {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      });

      alert("Upload Successful");

      setSelectedFile(null);

      onUploadSuccess();

    } catch (error) {
      alert(
        error.response?.data?.error ||
        "Upload Failed"
      );
    }
  };

  return (
    <div>
      <h2>Upload File</h2>

      <input
        type="file"
        onChange={(e) =>
          setSelectedFile(e.target.files[0])
        }
      />

      <button onClick={handleUpload}>
        Upload
      </button>
    </div>
  );
}

export default UploadForm;