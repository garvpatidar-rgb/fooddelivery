import React, { useContext } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { AuthContext } from '../context/AuthContext';
import { CartContext } from '../context/CartContext';
import { Utensils, ShoppingBag, LogOut, User, ClipboardList, Shield } from 'lucide-react';

const Navbar = () => {
  const { user, logout } = useContext(AuthContext);
  const { cartItems } = useContext(CartContext);
  const navigate = useNavigate();

  const totalQuantity = cartItems.reduce((sum, item) => sum + item.quantity, 0);

  return (
    <nav className="navbar">
      <div className="container nav-container">
        <Link to="/" className="logo">
          <Utensils size={28} />
          <span>TastyBites</span>
        </Link>

        <div className="nav-links">
          <Link to="/" className="nav-link">Menu</Link>

          {user ? (
            <>
              {user.role === 'admin' && (
                <Link 
                  to="/admin" 
                  className="nav-link" 
                  style={{ 
                    color: '#ff4757', 
                    fontWeight: 'bold', 
                    background: 'rgba(255, 71, 87, 0.15)', 
                    padding: '0.4rem 0.8rem', 
                    borderRadius: '8px', 
                    border: '1px solid rgba(255, 71, 87, 0.3)' 
                  }}
                >
                  <Shield size={18} />
                  Admin Portal
                </Link>
              )}
              <Link to="/orders" className="nav-link">
                <ClipboardList size={18} />
                Orders
              </Link>
              <div className="nav-link">
                <User size={18} />
                {user.name}
              </div>
              <button onClick={() => { logout(); navigate('/login'); }} className="nav-link" style={{ background: 'none', border: 'none', cursor: 'pointer' }}>
                <LogOut size={18} />
                Logout
              </button>
            </>
          ) : (
            <>
              <Link to="/login" className="nav-link">Login</Link>
              <Link to="/register" className="btn-primary">Register</Link>
            </>
          )}

          <Link to="/cart" className="cart-btn">
            <ShoppingBag size={20} />
            {totalQuantity > 0 && <span className="cart-badge">{totalQuantity}</span>}
          </Link>
        </div>
      </div>
    </nav>
  );
};

export default Navbar;
