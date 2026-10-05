import React, { useState, useEffect, useContext } from 'react';
import API from '../api/axios';
import { AuthContext } from '../context/AuthContext';
import { useNavigate } from 'react-router-dom';
import { 
  ShoppingBag, 
  DollarSign, 
  Utensils, 
  Clock, 
  Plus, 
  Edit3, 
  Trash2, 
  Search, 
  RefreshCw,
  CheckCircle,
  XCircle,
  Truck,
  ChefHat
} from 'lucide-react';
import toast from 'react-hot-toast';

const AdminDashboard = () => {
  const { user } = useContext(AuthContext);
  const navigate = useNavigate();

  const [activeTab, setActiveTab] = useState('orders'); // 'orders' | 'menu'
  
  // Orders State
  const [orders, setOrders] = useState([]);
  const [ordersLoading, setOrdersLoading] = useState(true);
  const [statusFilter, setStatusFilter] = useState('all');

  // Food Menu State
  const [foods, setFoods] = useState([]);
  const [foodsLoading, setFoodsLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState('');
  const [categoryFilter, setCategoryFilter] = useState('all');

  // Modal State
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingItem, setEditingItem] = useState(null);
  const [formData, setFormData] = useState({
    name: '',
    category: 'Burger',
    price: '',
    description: '',
    image_url: '',
    is_available: true
  });
  const [formSubmitting, setFormSubmitting] = useState(false);

  // Redirect if not admin
  useEffect(() => {
    if (!user || user.role !== 'admin') {
      toast.error('Admin access required');
      navigate('/');
    }
  }, [user, navigate]);

  // Fetch Orders
  const fetchOrders = async () => {
    setOrdersLoading(true);
    try {
      const res = await API.get('/admin/orders');
      setOrders(res.data.data || []);
    } catch (err) {
      console.error('Failed to fetch admin orders:', err);
      toast.error(err.response?.data?.error || 'Failed to load orders');
    } finally {
      setOrdersLoading(false);
    }
  };

  // Fetch Foods
  const fetchFoods = async () => {
    setFoodsLoading(true);
    try {
      const res = await API.get('/food?all=true');
      setFoods(res.data.data || []);
    } catch (err) {
      console.error('Failed to fetch food items:', err);
      toast.error('Failed to load food menu');
    } finally {
      setFoodsLoading(false);
    }
  };

  useEffect(() => {
    fetchOrders();
    fetchFoods();
  }, []);

  // Update Order Status
  const handleUpdateOrderStatus = async (orderId, newStatus) => {
    try {
      await API.put(`/orders/${orderId}/status`, { status: newStatus });
      toast.success(`Order #${orderId} status set to ${newStatus}`);
      setOrders(prev => prev.map(o => o.id === orderId ? { ...o, status: newStatus } : o));
    } catch (err) {
      toast.error('Failed to update status');
    }
  };

  // Open Modal for Create or Edit
  const openModal = (item = null) => {
    if (item) {
      setEditingItem(item);
      setFormData({
        name: item.name,
        category: item.category,
        price: item.price,
        description: item.description || '',
        image_url: item.image_url || '',
        is_available: item.is_available ?? true
      });
    } else {
      setEditingItem(null);
      setFormData({
        name: '',
        category: 'Burger',
        price: '',
        description: '',
        image_url: '',
        is_available: true
      });
    }
    setIsModalOpen(true);
  };

  const closeModal = () => {
    setIsModalOpen(false);
    setEditingItem(null);
  };

  // Save Food Item (Create / Edit)
  const handleSaveFood = async (e) => {
    e.preventDefault();
    setFormSubmitting(true);
    try {
      const payload = {
        name: formData.name,
        category: formData.category,
        price: parseFloat(formData.price),
        description: formData.description,
        image_url: formData.image_url,
        is_available: formData.is_available
      };

      if (editingItem) {
        const res = await API.put(`/admin/food/${editingItem.id}`, payload);
        toast.success('Food item updated!');
        setFoods(prev => prev.map(f => f.id === editingItem.id ? (res.data.data || { ...f, ...payload }) : f));
      } else {
        const res = await API.post('/admin/food', payload);
        toast.success('Food item created!');
        if (res.data.data) {
          setFoods(prev => [res.data.data, ...prev]);
        } else {
          fetchFoods();
        }
      }
      closeModal();
    } catch (err) {
      toast.error(err.response?.data?.error || 'Failed to save food item');
    } finally {
      setFormSubmitting(false);
    }
  };

  // Delete Food Item
  const handleDeleteFood = async (id, name) => {
    if (!window.confirm(`Are you sure you want to delete "${name}"?`)) return;
    try {
      await API.delete(`/admin/food/${id}`);
      toast.success('Item deleted successfully');
      setFoods(prev => prev.filter(f => f.id !== id));
    } catch (err) {
      toast.error('Failed to delete item');
    }
  };

  // Toggle Stock Availability
  const handleToggleAvailability = async (food) => {
    const updatedStatus = !food.is_available;
    try {
      await API.put(`/admin/food/${food.id}`, {
        name: food.name,
        category: food.category,
        price: food.price,
        description: food.description,
        image_url: food.image_url,
        is_available: updatedStatus
      });
      toast.success(`${food.name} is now ${updatedStatus ? 'In Stock' : 'Out of Stock'}`);
      setFoods(prev => prev.map(f => f.id === food.id ? { ...f, is_available: updatedStatus } : f));
    } catch (err) {
      toast.error('Failed to update availability');
    }
  };

  // Calculated Stats
  const totalRevenue = orders.reduce((sum, o) => o.status !== 'cancelled' ? sum + (o.total_amount || 0) : sum, 0);
  const pendingOrdersCount = orders.filter(o => o.status === 'pending').length;

  // Filtered Orders
  const filteredOrders = orders.filter(o => {
    if (statusFilter === 'all') return true;
    return o.status === statusFilter;
  });

  // Filtered Foods
  const categories = ['all', ...new Set(foods.map(f => f.category))];
  const filteredFoods = foods.filter(f => {
    const matchesSearch = f.name.toLowerCase().includes(searchQuery.toLowerCase()) || 
                          (f.description && f.description.toLowerCase().includes(searchQuery.toLowerCase()));
    const matchesCategory = categoryFilter === 'all' || f.category === categoryFilter;
    return matchesSearch && matchesCategory;
  });

  const getStatusBadgeStyle = (status) => {
    switch (status) {
      case 'pending': return { bg: 'rgba(255, 159, 26, 0.15)', color: '#ff9f1a', border: '#ff9f1a' };
      case 'preparing': return { bg: 'rgba(24, 144, 255, 0.15)', color: '#1890ff', border: '#1890ff' };
      case 'out_for_delivery': return { bg: 'rgba(114, 46, 209, 0.15)', color: '#722ed1', border: '#722ed1' };
      case 'delivered': return { bg: 'rgba(46, 213, 115, 0.15)', color: '#2ed573', border: '#2ed573' };
      case 'cancelled': return { bg: 'rgba(255, 71, 87, 0.15)', color: '#ff4757', border: '#ff4757' };
      default: return { bg: 'rgba(255, 255, 255, 0.1)', color: '#fff', border: '#aaa' };
    }
  };

  return (
    <div className="container admin-container" style={{ padding: '2rem 1.5rem' }}>
      
      {/* Admin Header */}
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '2rem', flexWrap: 'wrap', gap: '1rem' }}>
        <div>
          <h1 style={{ fontSize: '2.2rem', fontWeight: '800', background: 'linear-gradient(135deg, #ff4757, #ff6b81)', WebkitBackgroundClip: 'text', WebkitTextFillColor: 'transparent' }}>
            Restaurant Admin Portal 👑
          </h1>
          <p style={{ color: 'var(--text-muted)', fontSize: '1rem' }}>Manage orders, update live statuses, and curate your food menu.</p>
        </div>
        <button 
          onClick={() => { fetchOrders(); fetchFoods(); toast.success('Refreshed data!'); }} 
          className="btn-secondary"
          style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', padding: '0.6rem 1.2rem', borderRadius: '10px' }}
        >
          <RefreshCw size={18} /> Refresh Data
        </button>
      </div>

      {/* Overview Stats Cards */}
      <div className="admin-stats-grid" style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(220px, 1fr))', gap: '1.25rem', marginBottom: '2.5rem' }}>
        <div className="stat-card" style={{ background: 'var(--bg-card)', padding: '1.5rem', borderRadius: '16px', border: '1px solid var(--border-color)', display: 'flex', alignItems: 'center', gap: '1.2rem' }}>
          <div style={{ width: '52px', height: '52px', borderRadius: '14px', background: 'rgba(255, 71, 87, 0.15)', color: '#ff4757', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
            <DollarSign size={28} />
          </div>
          <div>
            <div style={{ color: 'var(--text-muted)', fontSize: '0.88rem', fontWeight: '600' }}>Total Revenue</div>
            <div style={{ fontSize: '1.6rem', fontWeight: '800', color: 'var(--text-main)', marginTop: '0.2rem' }}>₹{totalRevenue.toLocaleString()}</div>
          </div>
        </div>

        <div className="stat-card" style={{ background: 'var(--bg-card)', padding: '1.5rem', borderRadius: '16px', border: '1px solid var(--border-color)', display: 'flex', alignItems: 'center', gap: '1.2rem' }}>
          <div style={{ width: '52px', height: '52px', borderRadius: '14px', background: 'rgba(24, 144, 255, 0.15)', color: '#1890ff', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
            <ShoppingBag size={28} />
          </div>
          <div>
            <div style={{ color: 'var(--text-muted)', fontSize: '0.88rem', fontWeight: '600' }}>Total Orders</div>
            <div style={{ fontSize: '1.6rem', fontWeight: '800', color: 'var(--text-main)', marginTop: '0.2rem' }}>{orders.length}</div>
          </div>
        </div>

        <div className="stat-card" style={{ background: 'var(--bg-card)', padding: '1.5rem', borderRadius: '16px', border: '1px solid var(--border-color)', display: 'flex', alignItems: 'center', gap: '1.2rem' }}>
          <div style={{ width: '52px', height: '52px', borderRadius: '14px', background: 'rgba(255, 159, 26, 0.15)', color: '#ff9f1a', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
            <Clock size={28} />
          </div>
          <div>
            <div style={{ color: 'var(--text-muted)', fontSize: '0.88rem', fontWeight: '600' }}>Pending Orders</div>
            <div style={{ fontSize: '1.6rem', fontWeight: '800', color: 'var(--text-main)', marginTop: '0.2rem' }}>{pendingOrdersCount}</div>
          </div>
        </div>

        <div className="stat-card" style={{ background: 'var(--bg-card)', padding: '1.5rem', borderRadius: '16px', border: '1px solid var(--border-color)', display: 'flex', alignItems: 'center', gap: '1.2rem' }}>
          <div style={{ width: '52px', height: '52px', borderRadius: '14px', background: 'rgba(46, 213, 115, 0.15)', color: '#2ed573', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
            <Utensils size={28} />
          </div>
          <div>
            <div style={{ color: 'var(--text-muted)', fontSize: '0.88rem', fontWeight: '600' }}>Menu Items</div>
            <div style={{ fontSize: '1.6rem', fontWeight: '800', color: 'var(--text-main)', marginTop: '0.2rem' }}>{foods.length}</div>
          </div>
        </div>
      </div>

      {/* Navigation Tabs */}
      <div style={{ display: 'flex', gap: '1rem', borderBottom: '2px solid var(--border-color)', marginBottom: '2rem' }}>
        <button
          onClick={() => setActiveTab('orders')}
          style={{
            padding: '0.8rem 1.6rem',
            fontWeight: '700',
            fontSize: '1.05rem',
            background: 'none',
            border: 'none',
            borderBottom: activeTab === 'orders' ? '3px solid #ff4757' : '3px solid transparent',
            color: activeTab === 'orders' ? '#ff4757' : 'var(--text-muted)',
            cursor: 'pointer',
            display: 'flex',
            alignItems: 'center',
            gap: '0.5rem',
            transition: 'all 0.2s ease'
          }}
        >
          <ShoppingBag size={20} /> Order Management ({orders.length})
        </button>

        <button
          onClick={() => setActiveTab('menu')}
          style={{
            padding: '0.8rem 1.6rem',
            fontWeight: '700',
            fontSize: '1.05rem',
            background: 'none',
            border: 'none',
            borderBottom: activeTab === 'menu' ? '3px solid #ff4757' : '3px solid transparent',
            color: activeTab === 'menu' ? '#ff4757' : 'var(--text-muted)',
            cursor: 'pointer',
            display: 'flex',
            alignItems: 'center',
            gap: '0.5rem',
            transition: 'all 0.2s ease'
          }}
        >
          <Utensils size={20} /> Food Items Menu ({foods.length})
        </button>
      </div>

      {/* TAB 1: ORDER MANAGEMENT */}
      {activeTab === 'orders' && (
        <div>
          {/* Order Status Filters */}
          <div style={{ display: 'flex', gap: '0.6rem', flexWrap: 'wrap', marginBottom: '1.8rem' }}>
            {['all', 'pending', 'preparing', 'out_for_delivery', 'delivered', 'cancelled'].map(st => (
              <button
                key={st}
                onClick={() => setStatusFilter(st)}
                style={{
                  padding: '0.5rem 1rem',
                  borderRadius: '20px',
                  fontWeight: '600',
                  fontSize: '0.88rem',
                  border: '1px solid var(--border-color)',
                  background: statusFilter === st ? '#ff4757' : 'var(--bg-card)',
                  color: statusFilter === st ? '#fff' : 'var(--text-muted)',
                  cursor: 'pointer',
                  textTransform: 'capitalize',
                  transition: 'all 0.2s ease'
                }}
              >
                {st.replace(/_/g, ' ')}
              </button>
            ))}
          </div>

          {ordersLoading ? (
            <div style={{ textAlign: 'center', padding: '3rem', color: 'var(--text-muted)' }}>Loading live orders...</div>
          ) : filteredOrders.length === 0 ? (
            <div style={{ textAlign: 'center', padding: '3rem', background: 'var(--bg-card)', borderRadius: '16px', border: '1px solid var(--border-color)' }}>
              <ShoppingBag size={48} style={{ color: 'var(--text-muted)', marginBottom: '1rem' }} />
              <h3>No orders found for this status.</h3>
            </div>
          ) : (
            <div style={{ display: 'grid', gap: '1.5rem' }}>
              {filteredOrders.map((order) => {
                const badge = getStatusBadgeStyle(order.status);
                return (
                  <div key={order.id} style={{ background: 'var(--bg-card)', padding: '1.5rem', borderRadius: '16px', border: '1px solid var(--border-color)', boxShadow: '0 4px 16px rgba(0,0,0,0.1)' }}>
                    
                    {/* Header */}
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', flexWrap: 'wrap', gap: '1rem', marginBottom: '1rem', borderBottom: '1px solid var(--border-color)', paddingBottom: '1rem' }}>
                      <div>
                        <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem' }}>
                          <h3 style={{ fontSize: '1.25rem', fontWeight: '700' }}>Order #{order.id}</h3>
                          <span style={{ background: badge.bg, color: badge.color, border: `1px solid ${badge.border}`, padding: '0.25rem 0.75rem', borderRadius: '20px', fontSize: '0.82rem', fontWeight: '700', textTransform: 'uppercase' }}>
                            {order.status.replace(/_/g, ' ')}
                          </span>
                        </div>
                        <p style={{ color: 'var(--text-muted)', fontSize: '0.88rem', marginTop: '0.3rem' }}>
                          Placed on: {new Date(order.created_at).toLocaleString()}
                        </p>
                      </div>

                      {/* Customer Info */}
                      <div style={{ textAlign: 'right' }}>
                        <div style={{ fontWeight: '700', fontSize: '1.1rem' }}>{order.user?.name || 'Customer'}</div>
                        <div style={{ color: 'var(--text-muted)', fontSize: '0.88rem' }}>{order.user?.email}</div>
                      </div>
                    </div>

                    {/* Order Items */}
                    <div style={{ marginBottom: '1.2rem', background: 'var(--bg-dark)', padding: '1rem', borderRadius: '12px' }}>
                      <div style={{ fontSize: '0.88rem', fontWeight: '700', color: 'var(--text-muted)', marginBottom: '0.5rem', textTransform: 'uppercase' }}>Items Ordered</div>
                      {order.items?.map((item, idx) => (
                        <div key={idx} style={{ display: 'flex', justifyContent: 'space-between', fontSize: '0.95rem', marginBottom: '0.4rem' }}>
                          <span>{item.quantity}x {item.food_item?.name || `Item #${item.food_item_id}`}</span>
                          <span style={{ fontWeight: '600' }}>₹{(item.price * item.quantity).toFixed(2)}</span>
                        </div>
                      ))}
                    </div>

                    {/* Delivery & Total */}
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: '1rem', marginBottom: '1.2rem' }}>
                      <div style={{ fontSize: '0.9rem', color: 'var(--text-muted)' }}>
                        📍 <strong>Delivery Address:</strong> {order.delivery_address || 'Not specified'}
                      </div>
                      <div style={{ fontSize: '1.3rem', fontWeight: '800', color: '#ff4757' }}>
                        Total Amount: ₹{order.total_amount?.toFixed(2)}
                      </div>
                    </div>

                    {/* Action Buttons for Live Status Updates */}
                    <div style={{ display: 'flex', gap: '0.6rem', flexWrap: 'wrap', paddingTop: '1rem', borderTop: '1px dashed var(--border-color)', alignItems: 'center' }}>
                      <span style={{ fontSize: '0.88rem', fontWeight: '700', color: 'var(--text-muted)', marginRight: '0.5rem' }}>Update Status:</span>

                      {order.status === 'pending' && (
                        <button
                          onClick={() => handleUpdateOrderStatus(order.id, 'preparing')}
                          className="btn-primary"
                          style={{ padding: '0.4rem 0.9rem', fontSize: '0.85rem', display: 'flex', alignItems: 'center', gap: '0.4rem', background: '#1890ff' }}
                        >
                          <ChefHat size={16} /> Accept & Prepare
                        </button>
                      )}

                      {order.status === 'preparing' && (
                        <button
                          onClick={() => handleUpdateOrderStatus(order.id, 'out_for_delivery')}
                          className="btn-primary"
                          style={{ padding: '0.4rem 0.9rem', fontSize: '0.85rem', display: 'flex', alignItems: 'center', gap: '0.4rem', background: '#722ed1' }}
                        >
                          <Truck size={16} /> Dispatch for Delivery
                        </button>
                      )}

                      {order.status === 'out_for_delivery' && (
                        <button
                          onClick={() => handleUpdateOrderStatus(order.id, 'delivered')}
                          className="btn-primary"
                          style={{ padding: '0.4rem 0.9rem', fontSize: '0.85rem', display: 'flex', alignItems: 'center', gap: '0.4rem', background: '#2ed573' }}
                        >
                          <CheckCircle size={16} /> Mark as Delivered
                        </button>
                      )}

                      {order.status !== 'cancelled' && order.status !== 'delivered' && (
                        <button
                          onClick={() => handleUpdateOrderStatus(order.id, 'cancelled')}
                          style={{ padding: '0.4rem 0.9rem', fontSize: '0.85rem', display: 'flex', alignItems: 'center', gap: '0.4rem', background: 'rgba(255, 71, 87, 0.15)', color: '#ff4757', border: '1px solid #ff4757', borderRadius: '8px', cursor: 'pointer', fontWeight: '600' }}
                        >
                          <XCircle size={16} /> Cancel Order
                        </button>
                      )}

                      {/* Dropdown status selector fallback */}
                      <select
                        value={order.status}
                        onChange={(e) => handleUpdateOrderStatus(order.id, e.target.value)}
                        style={{ marginLeft: 'auto', padding: '0.4rem 0.8rem', borderRadius: '8px', background: 'var(--bg-dark)', color: 'var(--text-main)', border: '1px solid var(--border-color)', fontWeight: '600', fontSize: '0.85rem' }}
                      >
                        <option value="pending">Pending</option>
                        <option value="preparing">Preparing</option>
                        <option value="out_for_delivery">Out for Delivery</option>
                        <option value="delivered">Delivered</option>
                        <option value="cancelled">Cancelled</option>
                      </select>
                    </div>

                  </div>
                );
              })}
            </div>
          )}
        </div>
      )}

      {/* TAB 2: FOOD ITEMS MENU MANAGEMENT */}
      {activeTab === 'menu' && (
        <div>
          {/* Controls Bar: Search, Category Filter, and Add Button */}
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: '1rem', marginBottom: '2rem' }}>
            <div style={{ display: 'flex', gap: '1rem', flex: '1', minWidth: '280px' }}>
              <div style={{ position: 'relative', flex: '1' }}>
                <Search size={18} style={{ position: 'absolute', left: '12px', top: '50%', transform: 'translateY(-50%)', color: 'var(--text-muted)' }} />
                <input
                  type="text"
                  placeholder="Search food items..."
                  className="form-input"
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  style={{ paddingLeft: '2.4rem' }}
                />
              </div>

              <select
                className="form-input"
                value={categoryFilter}
                onChange={(e) => setCategoryFilter(e.target.value)}
                style={{ width: '160px' }}
              >
                {categories.map(cat => (
                  <option key={cat} value={cat}>{cat === 'all' ? 'All Categories' : cat}</option>
                ))}
              </select>
            </div>

            <button onClick={() => openModal()} className="btn-primary" style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', padding: '0.8rem 1.4rem' }}>
              <Plus size={20} /> Add New Food Item
            </button>
          </div>

          {/* Foods Table / Grid */}
          {foodsLoading ? (
            <div style={{ textAlign: 'center', padding: '3rem', color: 'var(--text-muted)' }}>Loading menu...</div>
          ) : filteredFoods.length === 0 ? (
            <div style={{ textAlign: 'center', padding: '3rem', background: 'var(--bg-card)', borderRadius: '16px', border: '1px solid var(--border-color)' }}>
              <Utensils size={48} style={{ color: 'var(--text-muted)', marginBottom: '1rem' }} />
              <h3>No food items match your search.</h3>
            </div>
          ) : (
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(280px, 1fr))', gap: '1.5rem' }}>
              {filteredFoods.map((food) => (
                <div key={food.id} style={{ background: 'var(--bg-card)', borderRadius: '16px', overflow: 'hidden', border: '1px solid var(--border-color)', display: 'flex', flexDirection: 'column' }}>
                  
                  {/* Image */}
                  <div style={{ position: 'relative', height: '160px', overflow: 'hidden', background: '#222' }}>
                    <img 
                      src={food.image_url || 'https://images.unsplash.com/photo-1546069901-ba9599a7e63c?w=500'} 
                      alt={food.name}
                      style={{ width: '100%', height: '100%', objectFit: 'cover' }}
                      onError={(e) => { e.target.src = 'https://images.unsplash.com/photo-1546069901-ba9599a7e63c?w=500'; }}
                    />
                    <span style={{ position: 'absolute', top: '12px', left: '12px', background: 'rgba(0,0,0,0.7)', color: '#fff', padding: '0.2rem 0.6rem', borderRadius: '6px', fontSize: '0.78rem', fontWeight: '700' }}>
                      {food.category}
                    </span>
                    <span style={{ position: 'absolute', top: '12px', right: '12px', background: food.is_available ? '#2ed573' : '#ff4757', color: '#fff', padding: '0.2rem 0.6rem', borderRadius: '6px', fontSize: '0.78rem', fontWeight: '700' }}>
                      {food.is_available ? 'In Stock' : 'Out of Stock'}
                    </span>
                  </div>

                  {/* Body */}
                  <div style={{ padding: '1.25rem', flex: '1', display: 'flex', flexDirection: 'column' }}>
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: '0.5rem' }}>
                      <h3 style={{ fontSize: '1.15rem', fontWeight: '700' }}>{food.name}</h3>
                      <span style={{ fontSize: '1.2rem', fontWeight: '800', color: '#ff4757' }}>₹{food.price}</span>
                    </div>

                    <p style={{ color: 'var(--text-muted)', fontSize: '0.88rem', marginBottom: '1.2rem', flex: '1', lineHeight: '1.4' }}>
                      {food.description || 'No description provided.'}
                    </p>

                    {/* Stock & Actions */}
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', borderTop: '1px solid var(--border-color)', paddingTop: '0.8rem' }}>
                      <button
                        onClick={() => handleToggleAvailability(food)}
                        style={{
                          background: 'none',
                          border: 'none',
                          color: food.is_available ? '#2ed573' : '#ff4757',
                          cursor: 'pointer',
                          fontSize: '0.85rem',
                          fontWeight: '700'
                        }}
                      >
                        {food.is_available ? '● Available' : '○ Unavailable'}
                      </button>

                      <div style={{ display: 'flex', gap: '0.5rem' }}>
                        <button
                          onClick={() => openModal(food)}
                          style={{ background: 'rgba(255, 255, 255, 0.1)', border: 'none', color: '#fff', padding: '0.4rem 0.7rem', borderRadius: '8px', cursor: 'pointer' }}
                          title="Edit Item"
                        >
                          <Edit3 size={16} />
                        </button>
                        <button
                          onClick={() => handleDeleteFood(food.id, food.name)}
                          style={{ background: 'rgba(255, 71, 87, 0.15)', border: 'none', color: '#ff4757', padding: '0.4rem 0.7rem', borderRadius: '8px', cursor: 'pointer' }}
                          title="Delete Item"
                        >
                          <Trash2 size={16} />
                        </button>
                      </div>
                    </div>
                  </div>

                </div>
              ))}
            </div>
          )}
        </div>
      )}

      {/* CREATE / EDIT FOOD ITEM MODAL */}
      {isModalOpen && (
        <div style={{ position: 'fixed', inset: 0, background: 'rgba(0,0,0,0.75)', backdropFilter: 'blur(4px)', display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 1000, padding: '1rem' }}>
          <div style={{ background: 'var(--bg-card)', border: '1px solid var(--border-color)', borderRadius: '20px', width: '100%', maxWidth: '520px', padding: '2rem', boxShadow: '0 20px 40px rgba(0,0,0,0.5)' }}>
            <h2 style={{ marginBottom: '1.5rem', fontSize: '1.5rem', fontWeight: '800' }}>
              {editingItem ? '✏️ Edit Food Item' : '➕ Add New Food Item'}
            </h2>

            <form onSubmit={handleSaveFood}>
              <div className="form-group" style={{ marginBottom: '1rem' }}>
                <label style={{ display: 'block', marginBottom: '0.4rem', fontWeight: '600' }}>Item Name *</label>
                <input
                  type="text"
                  className="form-input"
                  required
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  placeholder="e.g. Cheesy Paneer Burger"
                />
              </div>

              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem', marginBottom: '1rem' }}>
                <div className="form-group">
                  <label style={{ display: 'block', marginBottom: '0.4rem', fontWeight: '600' }}>Category *</label>
                  <input
                    type="text"
                    className="form-input"
                    required
                    value={formData.category}
                    onChange={(e) => setFormData({ ...formData, category: e.target.value })}
                    placeholder="e.g. Burger, Pizza"
                  />
                </div>

                <div className="form-group">
                  <label style={{ display: 'block', marginBottom: '0.4rem', fontWeight: '600' }}>Price (₹) *</label>
                  <input
                    type="number"
                    step="0.01"
                    className="form-input"
                    required
                    value={formData.price}
                    onChange={(e) => setFormData({ ...formData, price: e.target.value })}
                    placeholder="299.00"
                  />
                </div>
              </div>

              <div className="form-group" style={{ marginBottom: '1rem' }}>
                <label style={{ display: 'block', marginBottom: '0.4rem', fontWeight: '600' }}>Image URL</label>
                <input
                  type="url"
                  className="form-input"
                  value={formData.image_url}
                  onChange={(e) => setFormData({ ...formData, image_url: e.target.value })}
                  placeholder="https://images.unsplash.com/..."
                />
              </div>

              <div className="form-group" style={{ marginBottom: '1.2rem' }}>
                <label style={{ display: 'block', marginBottom: '0.4rem', fontWeight: '600' }}>Description</label>
                <textarea
                  className="form-input"
                  rows="3"
                  value={formData.description}
                  onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                  placeholder="Brief description of taste, ingredients..."
                ></textarea>
              </div>

              <div style={{ display: 'flex', alignItems: 'center', gap: '0.6rem', marginBottom: '1.5rem' }}>
                <input
                  type="checkbox"
                  id="is_available_check"
                  checked={formData.is_available}
                  onChange={(e) => setFormData({ ...formData, is_available: e.target.checked })}
                  style={{ width: '18px', height: '18px', cursor: 'pointer' }}
                />
                <label htmlFor="is_available_check" style={{ fontWeight: '600', cursor: 'pointer' }}>Item is currently Available in Stock</label>
              </div>

              <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '1rem' }}>
                <button type="button" onClick={closeModal} className="btn-secondary" style={{ padding: '0.7rem 1.2rem', borderRadius: '10px' }}>
                  Cancel
                </button>
                <button type="submit" disabled={formSubmitting} className="btn-primary" style={{ padding: '0.7rem 1.4rem', borderRadius: '10px' }}>
                  {formSubmitting ? 'Saving...' : (editingItem ? 'Update Item' : 'Create Item')}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

    </div>
  );
};

export default AdminDashboard;
