import React, { useEffect, useState } from "react";
import { fetchAllItems, createItem, updateItem, deleteItem, fetchItemById } from "../api/items";
import { useAuth } from "./AuthContext"; // Импортируем контекст
import './ItemsList.css';
import { ToastContainer, toast } from "react-toastify";
import 'react-toastify/dist/ReactToastify.css';

const ItemsList = () => {
    const { role } = useAuth(); // Получаем роль из контекста
    const [itemsByMonth, setItemsByMonth] = useState({});
    const [showModal, setShowModal] = useState(false);
    const [showEditModal, setShowEditModal] = useState(false);
    const [showDeleteModal, setShowDeleteModal] = useState(false);
    const [newItem, setNewItem] = useState({ Name: "", Price: "" });
    const [itemToEdit, setItemToEdit] = useState(null);
    const [originalItem, setOriginalItem] = useState(null);
    const [itemToDelete, setItemToDelete] = useState(null);
    const [searchId, setSearchId] = useState("");
    const [searchResult, setSearchResult] = useState(null);

    useEffect(() => {
        const loadItems = async () => {
            try {
                const data = await fetchAllItems();
                groupItemsByMonth(data);
            } catch (error) {
                console.error("Ошибка при загрузке элементов:", error);
            }
        };
        loadItems();
    }, []);

    const groupItemsByMonth = (items) => {
        const groupedItems = items.reduce((acc, item) => {
            if (item.CreatedAt) {
                const itemDate = new Date(item.CreatedAt);
                if (!isNaN(itemDate)) {
                    const monthYear = new Intl.DateTimeFormat('en-US', { month: 'long', year: 'numeric' }).format(itemDate);

                    if (!acc[monthYear]) {
                        acc[monthYear] = [];
                    }
                    acc[monthYear].push(item);
                } else {
                    console.error(`Некорректное значение даты: ${item.CreatedAt}`);
                }
            } else {
                console.error("Отсутствует поле CreatedAt у элемента:", item);
            }
            return acc;
        }, {});

        setItemsByMonth(groupedItems);
    };

    const handleCreateItem = async () => {
        try {
            const createdItem = await createItem({
                Name: newItem.Name,
                Price: parseFloat(newItem.Price) // Convert Price to float
            });
            const updatedItems = { ...itemsByMonth };
            const monthYear = new Intl.DateTimeFormat('en-US', { month: 'long', year: 'numeric' }).format(new Date(createdItem.CreatedAt));

            if (!updatedItems[monthYear]) {
                updatedItems[monthYear] = [];
            }
            updatedItems[monthYear].push(createdItem);

            setItemsByMonth(updatedItems);
            setNewItem({ Name: "", Price: "" });
            setShowModal(false);
        } catch (error) {
            toast.error("Ошибка при создании элемента:");
            console.error("Ошибка при создании элемента:", error);
        }
    };

    const handleEditItem = async () => {
        try {
            const updatedFields = {};
            if (itemToEdit.Name !== originalItem.Name) {
                updatedFields.name = itemToEdit.Name;
            }
            if (itemToEdit.Price !== originalItem.Price) {
                updatedFields.price = parseFloat(itemToEdit.Price); // Convert Price to float
            }

            if (Object.keys(updatedFields).length > 0) {
                const updatedItem = await updateItem(itemToEdit.ID, updatedFields);
                const updatedItems = { ...itemsByMonth };
                const monthYear = new Intl.DateTimeFormat('en-US', { month: 'long', year: 'numeric' }).format(new Date(updatedItem.CreatedAt));

                if (!updatedItems[monthYear]) {
                    updatedItems[monthYear] = [];
                }

                updatedItems[monthYear] = updatedItems[monthYear].map(item =>
                    item.ID === updatedItem.ID ? updatedItem : item
                );

                setItemsByMonth(updatedItems);
                setItemToEdit(null);
                setShowEditModal(false);
            }
        } catch (error) {
            toast.error("Ошибка при изменении элемента")
            console.error("Ошибка при изменении элемента:", error);
        }
    };

    const handleSearchById = async () => {
        try {
            const item = await fetchItemById(searchId);
            setSearchResult(item);
        } catch (error) {
            console.error("Ошибка при поиске элемента:", error);
            toast.error("Ошибка при поиске элемента")
            setSearchResult(null);
        }
    };

    const handleDeleteItem = async () => {
        try {
            await deleteItem(itemToDelete);
            const updatedItems = { ...itemsByMonth };

            Object.keys(updatedItems).forEach((month) => {
                updatedItems[month] = updatedItems[month].filter((item) => item.ID !== itemToDelete);
            });

            setItemsByMonth(updatedItems);
            setShowDeleteModal(false);
            setItemToDelete(null);
        } catch (error) {
            console.error("Ошибка при удалении элемента:", error);
        }
    };

    const handleOpenEditModal = (item) => {
        setOriginalItem({ ...item });
        setItemToEdit({ ...item });
        setShowEditModal(true);
    };

    return (
        <div className="items-list-container">
            <h2>All Items</h2>
            <ToastContainer />
            {/* Поиск по ID */}
            <div className="search-container">
                <input
                    type="text"
                    placeholder="Search item by ID"
                    value={searchId}
                    onChange={(e) => setSearchId(e.target.value)}
                />
                <button onClick={handleSearchById}>Search</button>
            </div>

            {searchResult && (
                <div className="search-result">
                    <h2>Search Result:</h2>
                    <div className="item-card search-card">
                        <h3>{searchResult.Name}</h3>
                        <p>Price: {searchResult.Price}</p>
                    </div>
                </div>
            )}

            <button className="create-item-btn" onClick={() => setShowModal(true)}>
                Create New Item
            </button>

            {Object.keys(itemsByMonth).map((month) => (
                <div key={month} className="item-month-group">
                    <h2>{month}:</h2>
                    <div className="item-list">
                        {itemsByMonth[month].map((item) => (
                            <div className="item-card-wrapper" key={item.ID}>
                                <div className="item-card">
                                    <h3>{item.Name}</h3>
                                    <p>Price: {item.Price || "No price available"}</p>
                                    {item.CreatedAt && (
                                        <p>{new Intl.DateTimeFormat('en-US', {
                                            year: 'numeric', month: 'long', day: 'numeric',
                                            hour: 'numeric', minute: 'numeric'
                                        }).format(new Date(item.CreatedAt))}</p>
                                    )}
                                    {/* Кнопка "Edit Item" только для ролей admin и manager */}
                                    {(role === "admin" || role === "manager") && (
                                        <button onClick={() => handleOpenEditModal(item)}>
                                            Edit Item
                                        </button>
                                    )}
                                    {/* Кнопка "Delete Item" только для роли admin */}
                                    {role === "admin" && (
                                        <button onClick={() => {
                                            setItemToDelete(item.ID);
                                            setShowDeleteModal(true);
                                        }}>
                                            Delete Item
                                        </button>
                                    )}
                                </div>
                            </div>
                        ))}
                    </div>
                </div>
            ))}

            {/* Модальные окна для создания и редактирования элементов */}
            {showModal && (
                <div className="modal-overlay">
                    <div className="modal">
                        <h2>Create New Item</h2>
                        <form onSubmit={(e) => { e.preventDefault(); handleCreateItem(); }}>
                            <div>
                                <label>Item Name</label>
                                <input
                                    type="text"
                                    value={newItem.Name}
                                    onChange={(e) => setNewItem({ ...newItem, Name: e.target.value })}
                                    required
                                />
                            </div>
                            <div>
                                <label>Item Price</label>
                                <input
                                    type="number"
                                    step="0.01"
                                    value={newItem.Price}
                                    onChange={(e) => setNewItem({ ...newItem, Price: e.target.value })}
                                    required
                                />
                            </div>
                            <div>
                                <button type="submit">Create Item</button>
                                <button type="button" onClick={() => setShowModal(false)}>Cancel</button>
                            </div>
                        </form>
                    </div>
                </div>
            )}

            {/* Модальные окна для редактирования и удаления */}
            {showEditModal && itemToEdit && (
                <div className="modal-overlay">
                    <div className="modal">
                        <h2>Edit Item</h2>
                        <form onSubmit={(e) => { e.preventDefault(); handleEditItem(); }}>
                            <div>
                                <label>Item Name</label>
                                <input
                                    type="text"
                                    value={itemToEdit.Name}
                                    onChange={(e) => setItemToEdit({ ...itemToEdit, Name: e.target.value })}
                                    required
                                />
                            </div>
                            <div>
                                <label>Item Price</label>
                                <input
                                    type="number"
                                    step="0.01"
                                    value={itemToEdit.Price}
                                    onChange={(e) => setItemToEdit({ ...itemToEdit, Price: e.target.value })}
                                    required
                                />
                            </div>
                            <div>
                                <button type="submit">Save Changes</button>
                                <button type="button" onClick={() => setShowEditModal(false)}>Cancel</button>
                            </div>
                        </form>
                    </div>
                </div>
            )}

            {showDeleteModal && (
                <div className="modal-overlay">
                    <div className="modal">
                        <h2>Are you sure you want to delete this item?</h2>
                        <div>
                            <button onClick={handleDeleteItem}>Yes, Delete</button>
                            <button onClick={() => setShowDeleteModal(false)}>Cancel</button>
                        </div>
                    </div>
                </div>
            )}
        </div>

    );
};

export default ItemsList;
