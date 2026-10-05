import React, { useEffect, useState, useContext } from 'react';
import API from '../api/axios';
import { AuthContext } from '../context/AuthContext';
import { Package } from 'lucide-react';

const Orders = () => {
  const [orders, setOrders] = useState([]);
  const [loading, setLoading] = useState(true);
  const { user } = useContext(AuthContext);

  useEffect(() => {
    if (user) {
      API.get('/orders')
        .then((res) => {
          setOrders(res.data.data);
        })
        .catch((err) => console.error(err))
        .finally(() => setLoading(false));
    }
  }, [user]);

  if (loading) return <div className="container" style={{ padding: '2rem' }}>Loading orders...</div>;

  return (
    <div className="container" style={{ padding: '2rem 1.5rem' }}>
      <h1 style={{ marginBottom: '2rem' }}>My Orders 📦</h1>
      
      {orders.length === 0 ? (
        <div style={{ textAlign: 'center', padding: '3rem', background: 'var(--bg-card)', borderRadius: '16px' }}>
          <Package size={48} style={{ color: 'var(--text-muted)', marginBottom: '1rem' }} />
          <h3 style={{ color: 'var(--text-muted)' }}>You haven't placed any orders yet.</h3>
        </div>
      ) : (
        <div style={{ display: 'grid', gap: '1.5rem' }}>
          {orders.map((order) => (
            <div key={order.id} style={{ background: 'var(--bg-card)', padding: '1.5rem', borderRadius: '12px', border: '1px solid var(--border-color)' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: '1rem', borderBottom: '1px solid var(--border-color)', paddingBottom: '1rem' }}>
                <div>
                  <h3 style={{ marginBottom: '0.25rem' }}>Order #{order.id}</h3>
                  <p style={{ color: 'var(--text-muted)', fontSize: '0.9rem' }}>{new Date(order.created_at).toLocaleString()}</p>
                </div>
                <div style={{ 
                  background: order.status === 'delivered' ? 'rgba(46, 213, 115, 0.2)' : 'rgba(255, 165, 2, 0.2)',
                  color: order.status === 'delivered' ? '#2ed573' : '#ffa502',
                  padding: '0.4rem 0.8rem',
                  borderRadius: '999px',
                  fontWeight: 'bold',
                  textTransform: 'capitalize',
                  fontSize: '0.9rem'
                }}>
                  {order.status}
                </div>
              </div>

              <div style={{ marginBottom: '1rem' }}>
                {order.items?.map((item, idx) => (
                  <div key={idx} style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '0.5rem', color: 'var(--text-muted)' }}>
                    <span>{item.quantity}x {item.food_item?.name || 'Item'}</span>
                    <span>₹{item.price * item.quantity}</span>
                  </div>
                ))}
              </div>

              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', paddingTop: '1rem', borderTop: '1px dashed var(--border-color)' }}>
                <span style={{ color: 'var(--text-muted)' }}>Deliver to: {order.delivery_address}</span>
                <span style={{ fontSize: '1.25rem', fontWeight: 'bold', color: 'var(--accent)' }}>Total: ₹{order.total_amount}</span>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default Orders;
