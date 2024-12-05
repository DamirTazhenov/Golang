import React, { useState } from "react";
import { login as login2 } from "../api/auth";
import { useNavigate } from "react-router-dom";
import { ToastContainer, toast } from "react-toastify";
import { useAuth } from "./AuthContext";
import 'react-toastify/dist/ReactToastify.css';
import './Login.css'; // Импортируем стили

const Login = ({ onLogin }) => {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const navigate = useNavigate();
    const { login } = useAuth();

    const handleLogin = async (e) => {
        e.preventDefault();
        try {
            const { token, role, user_id, name} = await login2(email, password);

            if (onLogin) {
                onLogin();
            }

            // Показать уведомление об успешной авторизации
            toast.success("Успешная авторизация!");
            login(token, role, user_id, name);
            setTimeout(() => {
                navigate("/");
            }, 2000); // 2 секунды задержки для показа уведомления
        } catch (error) {
            console.error("Ошибка при входе:", error);
            toast.error("Ошибка при входе");
        }
    };

    return (
        <div>
            <form onSubmit={handleLogin}>
                <h2>Login</h2>
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
                <button type="submit">Login</button>
            </form>

            <ToastContainer />
        </div>
    );
};

export default Login;
