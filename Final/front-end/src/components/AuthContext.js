import React, { createContext, useContext, useEffect, useState } from "react";
import { getToken, getRole, getUserId, logout as performLogout } from "../api/auth";

const AuthContext = createContext();

export const useAuth = () => useContext(AuthContext);

export const AuthProvider = ({ children }) => {
    const [isAuthenticated, setIsAuthenticated] = useState(!!getToken());
    const [role, setRole] = useState(getRole());
    const [userId, setUserId] = useState(getUserId());
    const [userName, setName] = useState(localStorage.getItem("userName"));

    const login = (token, role, user_id) => {
        localStorage.setItem("accessToken", token);
        localStorage.setItem("userRole", role);
        localStorage.setItem("userId", user_id);
        localStorage.setItem("userName", userName);

        setIsAuthenticated(true);
        setRole(role);
        setUserId(user_id);
        setName(userName);
    };

    const logout = () => {
        performLogout();
        setIsAuthenticated(false);
        setRole(null);
        setUserId(null);
        setName(null);
    };

    useEffect(() => {
        if (getToken()) {
            setIsAuthenticated(true);
            setRole(getRole());
            setUserId(getUserId());
            setName(localStorage.getItem("userName"));
        }
    }, []);

    return (
        <AuthContext.Provider value={{ isAuthenticated, role, userId, userName, login, logout }}>
            {children}
        </AuthContext.Provider>
    );
};
