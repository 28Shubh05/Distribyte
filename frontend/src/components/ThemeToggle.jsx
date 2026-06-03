import { Moon, Sun } from "lucide-react";

function ThemeToggle({ dark, setDark }) {
  return (
    <button
      onClick={() => setDark(!dark)}
      className="
        relative
        flex
        items-center
        w-16
        h-9
        rounded-full
        bg-slate-300
        dark:bg-slate-700
        transition-all
        duration-300
        p-1
      "
    >
      <div
        className={`
          absolute
          w-7
          h-7
          rounded-full
          bg-white
          shadow-md
          flex
          items-center
          justify-center
          transition-all
          duration-300
          ${dark ? "translate-x-7" : "translate-x-0"}
        `}
      >
        {dark ? (
          <Moon
            size={15}
            className="text-slate-700"
          />
        ) : (
          <Sun
            size={15}
            className="text-yellow-500"
          />
        )}
      </div>
    </button>
  );
}

export default ThemeToggle;