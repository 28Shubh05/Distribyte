import { useEffect, useState } from "react";
import api from "../services/api";
import UploadForm from "../components/UploadForm";
import toast from "react-hot-toast";
import { motion } from "framer-motion";

import {
  FileText,
  Image as ImageIcon,
  FileArchive,
  File,
  FileType,
  Download,
  Trash2,
  Search,
  HardDrive,
  Files,
  Activity,
  Loader2,
  Inbox,
} from "lucide-react";

function Dashboard() {
  const [files, setFiles] = useState([]);
  const [searchTerm, setSearchTerm] = useState("");
  const [loading, setLoading] = useState(true);


  const handleDownload = async (id, filename) => {
    try {
      const response = await api.get(
        `/download/${id}`,
        {
          responseType: "blob",
        }
      );

      const url = window.URL.createObjectURL(response.data);

      const link = document.createElement("a");
      link.href = url;
      link.download = filename;

      document.body.appendChild(link);
      link.click();

      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);

      toast.success("Download started");
    } catch (error) {
      console.error(error);

      toast.error(
        error.response?.data?.error ||
        "Download failed"
      );
    }
  };

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

      setFiles(response.data.files || []);
    } catch (error) {
      console.error(error);
    } finally {
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
    const ext = filename
      .split(".")
      .pop()
      .toLowerCase();

    const cls = "w-6 h-6";

    switch (ext) {
      case "pdf":
        return (
          <FileType
            className={`${cls} text-red-500`}
          />
        );

      case "png":
      case "jpg":
      case "jpeg":
        return (
          <ImageIcon
            className={`${cls} text-purple-500`}
          />
        );

      case "zip":
        return (
          <FileArchive
            className={`${cls} text-amber-500`}
          />
        );

      case "txt":
        return (
          <FileText
            className={`${cls} text-blue-500`}
          />
        );

      default:
        return (
          <File
            className={`${cls} text-slate-500`}
          />
        );
    }
  };

  const formatBytes = (bytes) => {
    if (bytes < 1024)
      return `${bytes} B`;

    if (bytes < 1024 * 1024)
      return `${(bytes / 1024).toFixed(2)} KB`;

    return `${(
      bytes /
      1024 /
      1024
    ).toFixed(2)} MB`;
  };

  const usagePercentage = Math.min(
    (totalStorage /
      (100 * 1024 * 1024)) *
      100,
    100
  );

  if (loading) {
    return (
      <div className="flex flex-col items-center justify-center py-20 text-slate-500 dark:text-slate-400">
        <Loader2 className="w-10 h-10 animate-spin mb-3" />
        <div className="text-lg">
          Loading...
        </div>
      </div>
    );
  }

  return (
    <motion.div
      className="space-y-8"
      initial={{
        opacity: 0,
        y: 15,
      }}
      animate={{
        opacity: 1,
        y: 0,
      }}
      transition={{
        duration: 0.3,
      }}
    >
      {/* Header */}

      <div>
        <h1 className="text-4xl font-bold text-slate-800 dark:text-slate-100">
          Dashboard
        </h1>

        <p className="text-slate-500 dark:text-slate-400 mt-1">
          Manage your files stored in
          Distribyte
        </p>
      </div>

      {/* Stats */}

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        <StatCard
          icon={<Files className="w-5 h-5" />}
          label="Total Files"
          value={files.length}
          accent="blue"
        />

        <StatCard
          icon={<HardDrive className="w-5 h-5" />}
          label="Storage Used"
          value={formatBytes(totalStorage)}
          accent="purple"
        />

        <StatCard
          icon={<Activity className="w-5 h-5" />}
          label="System Status"
          value="Healthy"
          accent="green"
        />
      </div>

      {/* Storage */}

      <div className="bg-white/90 dark:bg-slate-900/90 backdrop-blur-sm rounded-2xl shadow-sm border border-slate-200 dark:border-slate-800 p-5">
        <div className="flex justify-between mb-3">
          <span className="font-medium text-slate-700 dark:text-slate-200">
            Storage Usage
          </span>

          <span className="text-sm text-slate-500 dark:text-slate-400">
            {(
              totalStorage /
              1024 /
              1024
            ).toFixed(2)}{" "}
            MB of 100 MB
          </span>
        </div>

        <div className="w-full bg-slate-100 dark:bg-slate-800 rounded-full h-3 overflow-hidden">
          <div
            className="bg-gradient-to-r from-cyan-500 via-blue-500 to-indigo-600 h-3 rounded-full transition-all"
            style={{
              width: `${usagePercentage}%`,
            }}
          />
        </div>
      </div>

      {/* Upload */}

      <UploadForm
        onUploadSuccess={fetchFiles}
      />

      {/* Search */}

      <div className="bg-white/90 dark:bg-slate-900/90 backdrop-blur-sm rounded-2xl shadow-sm border border-slate-200 dark:border-slate-800 p-4">
        <div className="relative">
          <Search className="w-5 h-5 absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />

          <input
            type="text"
            placeholder="Search files..."
            value={searchTerm}
            onChange={(e) =>
              setSearchTerm(
                e.target.value
              )
            }
            className="w-full bg-transparent pl-10 py-2 outline-none text-slate-800 dark:text-slate-100"
          />
        </div>
      </div>

      {/* Files */}

      <div>
        <h2 className="text-xl font-semibold mb-4 text-slate-800 dark:text-slate-100">
          Stored Files
        </h2>

        {files.length === 0 ? (
          <EmptyState
            icon={
              <Inbox className="w-12 h-12" />
            }
            title="No files yet"
            description="Upload your first file to get started."
          />
        ) : filteredFiles.length === 0 ? (
          <EmptyState
            icon={
              <Search className="w-12 h-12" />
            }
            title="No matches found"
            description={`No files match "${searchTerm}"`}
          />
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
            {filteredFiles.map((file) => (
              <div
                key={file.id}
                className="bg-white/90 dark:bg-slate-900/90 backdrop-blur-sm rounded-2xl shadow-sm border border-slate-200 dark:border-slate-800 p-5 hover:shadow-xl hover:-translate-y-1 hover:border-blue-300 dark:hover:border-blue-700 transition-all flex flex-col"
              >
                <div className="flex gap-3 mb-4">
                  <div className="p-2 rounded-lg bg-slate-100 dark:bg-slate-800">
                    {getFileIcon(
                      file.original_name
                    )}
                  </div>

                  <div className="min-w-0">
                    <h3 className="font-semibold truncate text-slate-800 dark:text-slate-100">
                      {file.original_name}
                    </h3>

                    <p className="text-sm text-slate-500">
                      {formatBytes(
                        file.size
                      )}
                    </p>
                  </div>
                </div>

                <div className="flex gap-2 mt-auto">
                  <button
                    onClick={() =>
                      handleDownload(file.id, file.original_name)
                    }
                    className="flex-1 bg-blue-600 hover:bg-blue-700 text-white py-2 rounded-lg flex items-center justify-center gap-2 transition"
                  >
                    <Download className="w-4 h-4" />
                    Download
                  </button>

                  <button
                    onClick={() =>
                      handleDelete(file.id)
                    }
                    className
                    ="px-4 py-2 border border-red-200 dark:border-red-900 text-red-600 dark:text-red-400 rounded-lg hover:bg-red-50 dark:hover:bg-red-950/30 transition"
                  >
                    <Trash2 className="w-4 h-4" />
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </motion.div>
  );
}

function StatCard({
  icon,
  label,
  value,
  accent,
}) {
  const colors = {
    blue: "bg-blue-50 text-blue-600 dark:bg-blue-950/40 dark:text-blue-400",
    purple:
      "bg-purple-50 text-purple-600 dark:bg-purple-950/40 dark:text-purple-400",
    green:
      "bg-green-50 text-green-600 dark:bg-green-950/40 dark:text-green-400",
  };

  return (
    <div className="bg-white/90 dark:bg-slate-900/90 backdrop-blur-sm rounded-2xl shadow-sm border border-slate-200 dark:border-slate-800 p-5">
      <div className="flex justify-between items-center">
        <h3 className="text-sm text-slate-500 dark:text-slate-400">
          {label}
        </h3>

        <div
          className={`p-2 rounded-lg ${colors[accent]}`}
        >
          {icon}
        </div>
      </div>

      <p className="text-3xl font-extrabold tracking-tight mt-3 text-slate-800 dark:text-slate-100">
        {value}
      </p>
    </div>
  );
}

function EmptyState({
  icon,
  title,
  description,
}) {
  return (
    <div className="bg-white dark:bg-slate-900 rounded-2xl border border-dashed border-slate-300 dark:border-slate-700 p-10 text-center">
      <div className="flex justify-center text-slate-400 mb-3">
        {icon}
      </div>

      <h3 className="font-semibold text-slate-700 dark:text-slate-200">
        {title}
      </h3>

      <p className="text-sm text-slate-500 mt-1">
        {description}
      </p>
    </div>
  );
}

export default Dashboard;