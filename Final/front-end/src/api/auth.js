import axios from "axios";

const API_URL = "http://localhost:8080";

// Авторизация (получение JWT токена)
export const login = async (email, password) => {
    try {
        const response = await axios.post(`${API_URL}/login`, { email, password });
        const { token, role, user_id } = response.data;
        console.log(response.data)
        if (token) {
            localStorage.setItem("accessToken", token);
            localStorage.setItem("userRole", role);
            localStorage.setItem("userId", user_id);

            return { token, role, user_id };
        }
    } catch (error) {
        console.error("Ошибка при логине:", error);
        throw new Error("Неверные учетные данные");
    }
};

export const getToken = () => localStorage.getItem("accessToken");

export const getRole = () => localStorage.getItem("userRole");

export const getUserId = () => localStorage.getItem("userId");

// Регистрация с указанием роли
export const register = async (name, email, password, role = "user") => {
    try {
        const response = await axios.post(`${API_URL}/register`, {
            name,
            email,
            password,
            role,  // Передача роли в запросе
        });

        return response.data;
    } catch (error) {
        console.error("Ошибка при регистрации:", error.response?.data || error.message);
        throw new Error("Ошибка при регистрации");
    }
};

// Выход (удаление токена)
export const logout = () => {
    localStorage.removeItem("accessToken");
    localStorage.removeItem("userRole");
};
