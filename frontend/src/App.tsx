import { useMutation } from "@connectrpc/connect-query"
import React, { useState } from "react"
import { register } from "./gen/v1/user-UserService_connectquery"

export function App() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [name, setName] = useState("");
  const [response, setResponse] = useState<string>("");
  const useRegister = useMutation(register, {
    onSuccess: (data) => {
      setResponse(`User created with ID: ${data.id}`);
    },
    onError: (error) => {
      setResponse(`Error: ${error.message}`);
    },
  })

  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault();
      await useRegister.mutateAsync({
        email,
        password,
        name,
      });
  };

  return (
    <div style={{ padding: "2rem" }}>
      <h2>E-Commerce User Registration</h2>
      <form onSubmit={handleRegister}>
        <div>
          <input
            type="text"
            placeholder="Name"
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
        </div>
        <div>
          <input
            type="email"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />
        </div>
        <div>
          <input
            type="password"
            placeholder="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
        </div>
        <button type="submit">Register</button>
      </form>

      {response && <pre>{response}</pre>}
    </div>
  );
}

export default App;