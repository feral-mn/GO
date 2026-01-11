import { useState } from "react";
import axios from "axios";

function App() {
  const [name, setName] = useState("");
  const [age, setAge] = useState("");
  const [users, setUsers] = useState([]);

  const saveUser = async () => {
    await axios.post("http://localhost:8080/users", {
      name: name,
      age: Number(age),
    });

    setName("");
    setAge("");
  };

  const fetchUsers = async () => {
    const res = await axios.get("http://localhost:8080/users");
    setUsers(res.data);
  };

  return (
    <div style={{ padding: 20 }}>
      <h2>User Form</h2>

      <input
        placeholder="Name"
        value={name}
        onChange={(e) => setName(e.target.value)}
      />

      <input
        type="number"
        placeholder="Age"
        value={age}
        onChange={(e) => setAge(e.target.value)}
      />

      <br /><br />

      <button onClick={saveUser}>Save</button>
      <button onClick={fetchUsers}>Fetch</button>

      <ul>
        {users.map((u) => (
          <li key={u.id}>
            {u.name} - {u.age}
          </li>
        ))}
      </ul>
    </div>
  );
}

export default App;
