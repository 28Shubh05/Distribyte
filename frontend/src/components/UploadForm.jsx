import { useState, useCallback } from "react";
import api from "../services/api";
import toast from "react-hot-toast";
import { UploadCloud, File as FileIcon, X, Loader2 } from "lucide-react";

function UploadForm({ onUploadSuccess }) {
  const [selectedFile, setSelectedFile] = useState(null);
  const [isDragging, setIsDragging] = useState(false);
  const [uploading, setUploading] = useState(false);

  const handleUpload = async () => {
    if (!selectedFile) {
      toast.error("Select a file first");
      return;
    }

    const formData = new FormData();
    formData.append("file", selectedFile);

    setUploading(true);
    try {
      await api.post("/upload", formData, {
        headers: { "Content-Type": "multipart/form-data" },
      });
      toast.success("Upload Successful");
      setSelectedFile(null);
      onUploadSuccess();
    } catch (error) {
      toast.error(error.response?.data?.error || "Upload Failed");
    } finally {
      setUploading(false);
    }
  };

  const handleDrop = useCallback((e) => {
    e.preventDefault();
    setIsDragging(false);
    const file = e.dataTransfer.files?.[0];
    if (file) setSelectedFile(file);
  }, []);

  const formatSize = (bytes) => {
    if (bytes < 1024) return `${bytes} B`;
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
    return `${(bytes / 1024 / 1024).toFixed(2)} MB`;
  };

  return (
    <div className="bg-white dark:bg-slate-900 rounded-2xl shadow-sm border border-slate-200 dark:border-slate-800 p-6">
      <h2 className="text-xl font-semibold mb-4 text-slate-800 dark:text-slate-100">
        Upload File
      </h2>

      <label
        onDragOver={(e) => { e.preventDefault(); setIsDragging(true); }}
        onDragLeave={() => setIsDragging(false)}
        onDrop={handleDrop}
        className={`flex flex-col items-center justify-center gap-3 border-2 border-dashed rounded-xl p-8 cursor-pointer transition-all
          ${isDragging
            ? "border-blue-500 bg-blue-50 dark:bg-blue-950/30"
            : "border-slate-300 dark:border-slate-700 hover:border-blue-400 hover:bg-slate-50 dark:hover:bg-slate-800/50"
          }`}
      >
        <UploadCloud className="w-10 h-10 text-blue-600 dark:text-blue-400" />
        <div className="text-center">
          <p className="font-medium text-slate-700 dark:text-slate-200">
            Drop your file here, or <span className="text-blue-600 dark:text-blue-400">browse</span>
          </p>
          <p className="text-sm text-slate-500 dark:text-slate-400 mt-1">
            Any file type supported
          </p>
        </div>
        <input
          type="file"
          onChange={(e) => setSelectedFile(e.target.files[0])}
          className="hidden"
        />
      </label>

      {selectedFile && (
        <div className="mt-4 flex items-center justify-between bg-slate-50 dark:bg-slate-800 rounded-lg p-3 border border-slate-200 dark:border-slate-700">
          <div className="flex items-center gap-3 min-w-0">
            <FileIcon className="w-5 h-5 text-slate-500 dark:text-slate-400 shrink-0" />
            <div className="min-w-0">
              <p className="text-sm font-medium text-slate-800 dark:text-slate-100 truncate">
                {selectedFile.name}
              </p>
              <p className="text-xs text-slate-500 dark:text-slate-400">
                {formatSize(selectedFile.size)}
              </p>
            </div>
          </div>
          <button
            onClick={() => setSelectedFile(null)}
            className="p-1 rounded hover:bg-slate-200 dark:hover:bg-slate-700 text-slate-500"
          >
            <X className="w-4 h-4" />
          </button>
        </div>
      )}

      <button
        onClick={handleUpload}
        disabled={uploading || !selectedFile}
        className="mt-4 w-full bg-blue-600 hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed text-white font-medium px-6 py-3 rounded-lg flex items-center justify-center gap-2 transition"
      >
        {uploading ? (
          <><Loader2 className="w-4 h-4 animate-spin" /> Uploading...</>
        ) : (
          <><UploadCloud className="w-4 h-4" /> Upload</>
        )}
      </button>
    </div>
  );
}

export default UploadForm;
