import { useEffect, useState } from "react";
import api from "../services/api";
import toast from "react-hot-toast";
import { motion } from "framer-motion";

import {
  Trash2,
  RotateCcw,
  Inbox,
  HardDrive,
  Clock,
  Search,
  Loader2,
  Files,
} from "lucide-react";

function DeletedFiles() {
  const [files, setFiles] = useState([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState("");

  const fetchDeletedFiles = async () => {
    setLoading(true);

    try {
      const response = await api.get("/deleted-files");

      setFiles(response.data.files || []);
    } catch (error) {
      console.error(error);
    } finally {
      setLoading(false);
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

  const totalDeletedStorage = files.reduce(
    (sum, file) => sum + file.size,
    0
  );

  const filteredFiles = files.filter((file) =>
    file.original_name
      .toLowerCase()
      .includes(searchTerm.toLowerCase())
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
          Deleted Files
        </h1>

        <p className="text-slate-500 dark:text-slate-400 mt-1">
          Restore previously deleted files
        </p>
      </div>

      {/* Stats */}

      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <StatCard
          icon={<Trash2 className="w-5 h-5" />}
          label="Deleted Files"
          value={files.length}
          accent="red"
        />

        <StatCard
          icon={<HardDrive className="w-5 h-5" />}
          label="Recoverable Storage"
          value={formatBytes(totalDeletedStorage)}
          accent="amber"
        />
      </div>

      {/* Search */}

      <div className="bg-white/90 dark:bg-slate-900/90 backdrop-blur-sm rounded-2xl shadow-sm border border-slate-200 dark:border-slate-800 p-4">
        <div className="relative">
          <Search className="w-5 h-5 absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />

          <input
            type="text"
            placeholder="Search deleted files..."
            value={searchTerm}
            onChange={(e) =>
              setSearchTerm(e.target.value)
            }
            className="w-full bg-transparent pl-10 py-2 outline-none text-slate-800 dark:text-slate-100"
          />
        </div>
      </div>

      {/* Deleted Files */}

      <div>
        <h2 className="text-xl font-semibold mb-4 text-slate-800 dark:text-slate-100">
          Deleted Items
        </h2>

        {files.length === 0 ? (
          <EmptyState
            icon={
              <Inbox className="w-12 h-12" />
            }
            title="No deleted files"
            description="Deleted files will appear here."
          />
        ) : filteredFiles.length === 0 ? (
          <EmptyState
            icon={
              <Search className="w-12 h-12" />
            }
            title="No matches found"
            description={`No deleted files match "${searchTerm}"`}
          />
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
            {filteredFiles.map((file) => (
              <div
                key={file.id}
                className="
                  bg-white/90
                  dark:bg-slate-900/90
                  backdrop-blur-sm
                  rounded-2xl
                  shadow-sm
                  border
                  border-slate-200
                  dark:border-slate-800
                  p-5
                  hover:shadow-xl
                  hover:-translate-y-1
                  hover:border-red-300
                  dark:hover:border-red-700
                  transition-all
                  flex
                  flex-col
                "
              >
                <div className="flex items-start gap-3 mb-4">
                  <div className="p-2 rounded-lg bg-red-50 dark:bg-red-950/40">
                    <Trash2 className="w-5 h-5 text-red-500" />
                  </div>

                  <div className="min-w-0 flex-1">
                    <h3
                      className="font-semibold truncate text-slate-800 dark:text-slate-100"
                      title={file.original_name}
                    >
                      {file.original_name}
                    </h3>

                    <p className="text-sm text-slate-500">
                      {formatBytes(file.size)}
                    </p>

                    {file.deleted_at && (
                      <div className="flex items-center gap-1 text-xs text-red-500 mt-2">
                        <Clock className="w-3 h-3" />

                        {new Date(
                          file.deleted_at
                        ).toLocaleString()}
                      </div>
                    )}
                  </div>
                </div>

                <button
                  onClick={() =>
                    handleRestore(file.id)
                  }
                  className="
                    mt-auto
                    bg-green-600
                    hover:bg-green-700
                    text-white
                    py-2
                    rounded-lg
                    flex
                    items-center
                    justify-center
                    gap-2
                    transition
                  "
                >
                  <RotateCcw className="w-4 h-4" />

                  Restore
                </button>
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
    red:
      "bg-red-50 text-red-600 dark:bg-red-950/40 dark:text-red-400",

    amber:
      "bg-amber-50 text-amber-600 dark:bg-amber-950/40 dark:text-amber-400",
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

export default DeletedFiles;