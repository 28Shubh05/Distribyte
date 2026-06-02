import { useState } from "react";
import api from "../services/api";
import toast from "react-hot-toast";

function UploadForm({ onUploadSuccess }) {

  const [selectedFile, setSelectedFile] = useState(null);

  const handleUpload = async () => {

    if (!selectedFile) {
      toast.error("Select a file first");
      return;
    }

    const formData = new FormData();

    formData.append("file", selectedFile);

    try {

      await api.post(
        "/upload",
        formData,
        {
          headers: {
            "Content-Type":
              "multipart/form-data",
          },
        }
      );

      toast.success("Upload Successful");

      setSelectedFile(null);

      onUploadSuccess();

    } catch (error) {

      toast.error(
        error.response?.data?.error ||
        "Upload Failed"
      );
    }
  };

  return (
    <div className="bg-white rounded-xl shadow p-6">

      <h2 className="text-2xl font-bold mb-4">
        Upload File
      </h2>

      <div className="flex flex-col md:flex-row gap-4">

        <input
          type="file"
          onChange={(e) =>
            setSelectedFile(e.target.files[0])
          }
          className="border rounded-lg p-3 flex-1"
        />

        <button
          onClick={handleUpload}
          className="bg-blue-600 text-white px-6 py-3 rounded-lg hover:bg-blue-700"
        >
          Upload
        </button>

      </div>

      {selectedFile && (
        <p className="mt-3 text-gray-600">
          Selected: {selectedFile.name}
        </p>
      )}

    </div>
  );
}

export default UploadForm;