import React, { useContext, useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { CartContext } from '../context/CartContext';
import { AuthContext } from '../context/AuthContext';
import API from '../api/axios';
import { Minus, Plus, Trash2 } from 'lucide-react';

const Cart = () => {
  const { cartItems, updateQuantity, removeFromCart, totalAmount, clearCart } = useContext(CartContext);
  const { user } = useContext(AuthContext);
  const [address, setAddress] = useState(user?.address || '');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const handleCheckout = async () => {
    if (!user) {
      navigate('/login');
      return;
    }
    if (!address) {
      setError('Please provide a delivery address');
      return;
    }
    
    setLoading(true);
    setError('');
    try {
      const items = cartItems.map(item => ({
        food_item_id: item.id,
        quantity: item.quantity
      }));
      
      await API.post('/orders', { items, delivery_address: address });
      clearCart();
      navigate('/orders');
    } catch (err) {
      setError(err.response?.data?.error || 'Checkout failed');
    } finally {
      setLoading(false);
    }
  };

  if (cartItems.length === 0) {
    return (
      <div className="container" style={{ padding: '4rem 1.5rem', textAlign: 'center' }}>
        <h2 style={{ marginBottom: '1rem' }}>Your Cart is Empty 🛒</h2>
        <p style={{ color: 'var(--text-muted)', marginBottom: '2rem' }}>Looks like you haven't added any food yet.</p>
        <Link to="/" className="btn-primary">Browse Menu</Link>
      </div>
    );
  }

  return (
    <div className="container" style={{ padding: '2rem 1.5rem' }}>
      <h1 style={{ marginBottom: '2rem' }}>Checkout 🛒</h1>
      
      <div style={{ display: 'grid', gap: '2rem', gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))' }}>
        <div>
          {cartItems.map((item) => (
            <div key={item.id} style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', background: 'var(--bg-card)', padding: '1rem', borderRadius: '12px', marginBottom: '1rem', border: '1px solid var(--border-color)' }}>
              <div>
                <h4 style={{ fontSize: '1.1rem', marginBottom: '0.25rem' }}>{item.name}</h4>
                <div style={{ color: 'var(--accent)', fontWeight: 'bold' }}>₹{item.price * item.quantity}</div>
              </div>
              
              <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
                <div style={{ display: 'flex', alignItems: 'center', background: 'var(--bg-dark)', borderRadius: '8px', padding: '0.25rem' }}>
                  <button onClick={() => updateQuantity(item.id, item.quantity - 1)} style={{ background: 'transparent', color: 'white', padding: '0.25rem' }}>
                    <Minus size={16} />
                  </button>
                  <span style={{ margin: '0 0.75rem', fontWeight: 'bold' }}>{item.quantity}</span>
                  <button onClick={() => updateQuantity(item.id, item.quantity + 1)} style={{ background: 'transparent', color: 'white', padding: '0.25rem' }}>
                    <Plus size={16} />
                  </button>
                </div>
                <button onClick={() => removeFromCart(item.id)} style={{ background: 'transparent', color: '#ff4757' }}>
                  <Trash2 size={20} />
                </button>
              </div>
            </div>
          ))}
        </div>

        <div style={{ background: 'var(--bg-card)', padding: '2rem', borderRadius: '16px', border: '1px solid var(--border-color)', height: 'fit-content' }}>
          <h3 style={{ marginBottom: '1.5rem', fontSize: '1.5rem', borderBottom: '1px solid var(--border-color)', paddingBottom: '1rem' }}>Order Summary</h3>
          
          <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '1rem', fontSize: '1.1rem' }}>
            <span>Subtotal</span>
            <span>₹{totalAmount}</span>
          </div>
          <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '1.5rem', fontSize: '1.1rem', color: 'var(--accent)', fontWeight: 'bold' }}>
            <span>Total</span>
            <span>₹{totalAmount}</span>
          </div>

          {error && <p style={{ color: '#ff4757', marginBottom: '1rem' }}>{error}</p>}

          <div className="form-group">
            <label>Delivery Address</label>
            <textarea
              className="form-input"
              rows="3"
              value={address}
              onChange={(e) => setAddress(e.target.value)}
              placeholder="Enter your full address..."
              required
            ></textarea>
          </div>

          <button onClick={handleCheckout} disabled={loading} className="btn-primary" style={{ width: '100%', padding: '1rem', fontSize: '1.1rem', marginTop: '1rem' }}>
            {loading ? 'Processing...' : 'Place Order'}
          </button>
        </div>
      </div>
    </div>
  );
};

export default Cart;
