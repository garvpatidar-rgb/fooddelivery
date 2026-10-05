import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import API from '../api/axios';
import toast from 'react-hot-toast';

const Register = () => {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [address, setAddress] = useState('');
  const [role, setRole] = useState('customer');
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    try {
      await API.post('/auth/register', { name, email, password, address, role });
      toast.success('Registration successful! Please login.');
      navigate('/login');
    } catch (err) {
      setError(err.response?.data?.error || 'Registration failed.');
    }
  };

  return (
    <div className="auth-wrapper">
      <div className="auth-card">
        <h2 className="auth-title">Create an Account 🎉</h2>
        {error && <p style={{ color: '#ff4757', marginBottom: '1rem', textAlign: 'center' }}>{error}</p>}
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label>Full Name</label>
            <input
              type="text"
              className="form-input"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
            />
          </div>
          <div className="form-group">
            <label>Email Address</label>
            <input
              type="email"
              className="form-input"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
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
              required
            />
          </div>
          <div className="form-group">
            <label>Account Type</label>
            <select
              className="form-input"
              value={role}
              onChange={(e) => setRole(e.target.value)}
            >
              <option value="customer">Customer (Order Food)</option>
              <option value="admin">Restaurant Admin (Manage Menu & Orders)</option>
            </select>
          </div>
          <div className="form-group">
            <label>Delivery Address</label>
            <input
              type="text"
              className="form-input"
              value={address}
              onChange={(e) => setAddress(e.target.value)}
              placeholder="Full home or office address..."
            />
          </div>
          <button type="submit" className="btn-primary" style={{ width: '100%', marginTop: '1rem' }}>
            Register
          </button>
        </form>
        <p style={{ marginTop: '1.5rem', textAlign: 'center', color: 'var(--text-muted)' }}>
          Already have an account? <Link to="/login" style={{ color: 'var(--primary)' }}>Login</Link>
        </p>
      </div>
    </div>
  );
};

export default Register;
