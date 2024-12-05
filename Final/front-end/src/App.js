import React, { useState } from "react";
import { BrowserRouter as Router, Route, Routes, Link } from "react-router-dom";
import { AuthProvider, useAuth } from "./components/AuthContext";
import Home from "./components/Home";
import ItemList from "./components/ItemList";
import Register from "./components/Register";
import Login from "./components/Login";
import './App.css';

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  const handleLogin = () => {
    setIsAuthenticated(true);
  };

  return (
      <AuthProvider>
          <Router>
              <div className="container">
                  <nav>
                      <Link to="/">Home</Link>
                      <Link to="/register">Register</Link>
                      <Link to="/login">Login</Link>
                      <Link to="/items">Items</Link>
                  </nav>
                  <Routes>
                      <Route exact path="/" element={<Home/>}/>
                      <Route path="/register" element={<Register/>}/>
                      <Route path="/login" element={<Login onLogin={handleLogin}/>}/> {/* Передаем onLogin */}
                      <Route path="/items" element={<ItemList/>}/>
                  </Routes>
              </div>
          </Router>
      </AuthProvider>
  );
}

export default App;
