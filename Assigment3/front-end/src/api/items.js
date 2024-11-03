import axios from "axios";
import { getToken } from "./auth";

const API_URL = "http://localhost:8080/api";

// Получение всех элементов
export const fetchAllItems = async () => {
    const response = await axios.get(`${API_URL}/items`, {
        headers: {
            Authorization: `Bearer ${getToken()}`,
        },
    });
    return response.data;
};

export const fetchItemById = async (itemId) => {
    try {
        const response = await axios.get(`${API_URL}/items/${itemId}`, {
            headers: {
                Authorization: `Bearer ${getToken()}`,
            },
        });
        return response.data;
    } catch (error) {
        console.error("Ошибка при получении элемента:", error);
        throw error;
    }
};

// Функция для создания нового элемента
export const createItem = async (item) => {
    try {
        const response = await axios.post(`${API_URL}/items`, item, {
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('accessToken')}`
            }
        });
        return response.data;
    } catch (error) {
        console.error("Ошибка при создании элемента:", error);
        throw error;
    }
};

// Функция для обновления элемента
export const updateItem = async (itemId, updatedItem) => {
    try {
        const response = await axios.put(`${API_URL}/items/${itemId}`, updatedItem, {
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('accessToken')}`
            }
        });
        return response.data;
    } catch (error) {
        console.error("Ошибка при обновлении элемента:", error);
        throw error;
    }
};

// Функция для удаления элемента
export const deleteItem = async (itemId) => {
    try {
        const response = await axios.delete(`${API_URL}/items/${itemId}`, {
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('accessToken')}`
            }
        });
        return response.data;
    } catch (error) {
        console.error("Ошибка при удалении элемента:", error);
        throw error;
    }
};
