import { useState } from "react";
import api from "../services/api";
import toast from "react-hot-toast";

function Register() {

  const [name,setName] = useState("");
  const [email,setEmail] = useState("");
  const [password,setPassword] = useState("");

  const handleSubmit = async (e) => {
    e.preventDefault();

    try {

      await api.post("/register", {
        name,
        email,
        password,
      });

      toast.success(
        "Registration successful"
      );

    } catch (err) {

      toast.error(
        err.response?.data?.error ||
        "Registration failed"
      );
    }
  };

  return (
    <div className="max-w-md mx-auto mt-20 bg-white p-8 rounded-2xl shadow">
      <h1 className="text-3xl font-bold mb-6">
        Register
      </h1>

      <form
        onSubmit={handleSubmit}
        className="space-y-4"
      >
        <input
          placeholder="Name"
          className="w-full border p-3 rounded-lg"
          value={name}
          onChange={(e)=>
            setName(e.target.value)
          }
        />

        <input
          type="email"
          placeholder="Email"
          className="w-full border p-3 rounded-lg"
          value={email}
          onChange={(e)=>
            setEmail(e.target.value)
          }
        />

        <input
          type="password"
          placeholder="Password"
          className="w-full border p-3 rounded-lg"
          value={password}
          onChange={(e)=>
            setPassword(e.target.value)
          }
        />

        <button
          className="w-full bg-green-600 text-white py-3 rounded-lg"
        >
          Register
        </button>
      </form>
    </div>
  );
}

export default Register;