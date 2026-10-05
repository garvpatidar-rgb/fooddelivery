import React, { useContext } from 'react';
import { CartContext } from '../context/CartContext';
import { Plus } from 'lucide-react';

const FoodCard = ({ food }) => {
  const { addToCart } = useContext(CartContext);

  const defaultImage = "https://images.unsplash.com/photo-1546069901-ba9599a7e63c?auto=format&fit=crop&w=600&q=80";

  return (
    <div className="food-card">
      <img
        src={food.image_url || defaultImage}
        alt={food.name}
        className="food-img"
      />
      <div className="food-body">
        <h3 className="food-title">{food.name}</h3>
        <p className="food-desc">{food.description || 'Delicious meal prepared with fresh ingredients.'}</p>
        <div className="food-footer">
          <span className="food-price">₹{food.price}</span>
          <button onClick={() => addToCart(food)} className="btn-primary" style={{ display: 'flex', alignItems: 'center', gap: '0.25rem' }}>
            <Plus size={18} /> Add
          </button>
        </div>
      </div>
    </div>
  );
};

export default FoodCard;
