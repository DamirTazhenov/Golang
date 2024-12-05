import React, { useState } from "react";
import { register } from "../api/auth";
import { useNavigate } from "react-router-dom"; // Используем useNavigate вместо useHistory

const Register = () => {
    const [name, setName] = useState("");
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [role, setRole] = useState("user"); // Роль по умолчанию
    const [error, setError] = useState(""); // Переменная для ошибки
    const navigate = useNavigate(); // Используем useNavigate для перенаправления

    const handleRegister = async (e) => {
        e.preventDefault();
        try {
            await register(name, email, password, role);

            navigate("/login");
        } catch (error) {
            console.error("Ошибка при регистрации:", error);
            setError("Ошибка при регистрации. Пожалуйста, попробуйте еще раз."); // Устанавливаем сообщение об ошибке
        }
    };

    return (
        <form onSubmit={handleRegister}>
            <h2>Регистрация</h2>
            {error && <p className="error">{error}</p>}
            <div>
                <input
                    type="text"
                    placeholder="Имя"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    required
                />
            </div>
            <div>
                <input
                    type="email"
                    placeholder="Электронная почта"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    required
                />
            </div>
            <div>
                <input
                    type="password"
                    placeholder="Пароль"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    required
                />
            </div>
            <div>
                <label>Роль</label>
                <select value={role} onChange={(e) => setRole(e.target.value)}>
                    <option value="user">User</option>
                    <option value="manager">Manager</option>
                    <option value="admin">Admin</option>
                </select>
            </div>
            <button type="submit">Зарегистрироваться</button>
        </form>
    );
};

export default Register;
