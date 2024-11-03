import React, { useState } from "react";
import { createItem } from "../api/items"; // Use createItem instead of createTask

const ItemForm = () => {
    const [title, setTitle] = useState("");
    const [description, setDescription] = useState("");

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            await createItem({ title, description }); // Create an item with title and description
            setTitle("");
            setDescription("");
        } catch (error) {
            console.error("Ошибка при создании элемента:", error);
        }
    };

    return (
        <form onSubmit={handleSubmit}>
            <h2>Create Item</h2>
            <div>
                <input
                    type="text"
                    placeholder="Name"
                    value={title}
                    onChange={(e) => setTitle(e.target.value)}
                    required
                />
            </div>
            <div>
                <input
                    type="text"
                    placeholder="Description"
                    value={description}
                    onChange={(e) => setDescription(e.target.value)}
                    required
                />
            </div>
            <button type="submit">Create Item</button>
        </form>
    );
};

export default ItemForm;
