import React, { useState, useContext } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import API from '../api/axios';
import { AuthContext } from '../context/AuthContext';
import toast from 'react-hot-toast';

const Login = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const { login } = useContext(AuthContext);
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    try {
      const res = await API.post('/auth/login', { email, password });
      login(res.data.user, res.data.token);
      toast.success(`Welcome back, ${res.data.user.name}!`);
      if (res.data.user.role === 'admin') {
        navigate('/admin');
      } else {
        navigate('/');
      }
    } catch (err) {
      setError(err.response?.data?.error || 'Login failed. Please check credentials.');
    }
  };

  const fillQuickLogin = (quickEmail, quickPass) => {
    setEmail(quickEmail);
    setPassword(quickPass);
  };

  return (
    <div className="auth-wrapper">
      <div className="auth-card">
        <h2 className="auth-title">Welcome Back 👋</h2>
        {error && <p style={{ color: '#ff4757', marginBottom: '1rem', textAlign: 'center' }}>{error}</p>}
        
        {/* Quick Demo Login Badges */}
        <div style={{ background: 'var(--bg-dark)', padding: '1rem', borderRadius: '12px', marginBottom: '1.5rem', border: '1px dashed var(--border-color)' }}>
          <div style={{ fontSize: '0.82rem', color: 'var(--text-muted)', fontWeight: '700', marginBottom: '0.6rem', textTransform: 'uppercase' }}>
            ⚡ Demo Accounts (Click to Autofill)
          </div>
          <div style={{ display: 'flex', gap: '0.5rem', flexWrap: 'wrap' }}>
            <button 
              type="button" 
              onClick={() => fillQuickLogin('admin@food.com', 'admin123')}
              style={{ background: 'rgba(255, 71, 87, 0.15)', color: '#ff4757', border: '1px solid #ff4757', padding: '0.3rem 0.6rem', borderRadius: '6px', cursor: 'pointer', fontSize: '0.82rem', fontWeight: '600' }}
            >
              👑 Admin Portal
            </button>
            <button 
              type="button" 
              onClick={() => fillQuickLogin('customer@food.com', 'customer123')}
              style={{ background: 'rgba(46, 213, 115, 0.15)', color: '#2ed573', border: '1px solid #2ed573', padding: '0.3rem 0.6rem', borderRadius: '6px', cursor: 'pointer', fontSize: '0.82rem', fontWeight: '600' }}
            >
              🍔 Customer
            </button>
          </div>
        </div>

        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label>Email Address</label>
            <input
              type="email"
              className="form-input"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="admin@food.com"
              required
            />
          </div>
          <div className="form-group">
            <label>Password</label>
            <input
              type="password"
              className="form-input"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="••••••••"
              required
            />
          </div>
          <button type="submit" className="btn-primary" style={{ width: '100%', marginTop: '1rem' }}>
            Login
          </button>
        </form>
        <p style={{ marginTop: '1.5rem', textAlign: 'center', color: 'var(--text-muted)' }}>
          Don't have an account? <Link to="/register" style={{ color: 'var(--primary)' }}>Register</Link>
        </p>
      </div>
    </div>
  );
};

export default Login;
