import React, { useEffect, useState } from 'react';
import API from '../api/axios';
import FoodCard from '../components/FoodCard';

const Home = () => {
  const [foods, setFoods] = useState([]);
  const [loading, setLoading] = useState(true);

  // Demo fallback food items if backend has no items yet
  const fallbackFoods = [
    { id: 1, name: 'Margherita Pizza', description: 'Classic cheese & tomato pizza', price: 299, image_url: 'https://images.unsplash.com/photo-1513104890138-7c749659a591?w=500' },
    { id: 2, name: 'Juicy Cheese Burger', description: 'Double beef patty with cheddar cheese', price: 199, image_url: 'https://images.unsplash.com/photo-1568901346375-23c9450c58cd?w=500' },
    { id: 3, name: 'Creamy Pasta Alfredo', description: 'Fettuccine in rich garlic cream sauce', price: 249, image_url: 'https://images.unsplash.com/photo-1621996346565-e3d5d6281288?w=500' },
  ];

  useEffect(() => {
    API.get('/food')
      .then((res) => {
        if (res.data.data && res.data.data.length > 0) {
          setFoods(res.data.data);
        } else {
          setFoods(fallbackFoods);
        }
      })
      .catch(() => {
        setFoods(fallbackFoods);
      })
      .finally(() => setLoading(false));
  }, []);

  return (
    <div className="container" style={{ padding: '2rem 1.5rem' }}>
      <h1 style={{ fontSize: '2rem', marginBottom: '0.5rem' }}>Delicious Food Delivered Fast 🚀</h1>
      <p style={{ color: 'var(--text-muted)', marginBottom: '1.5rem' }}>Choose your favorite meal from our top-rated menu.</p>

      {loading ? (
        <p>Loading tasty food...</p>
      ) : (
        <div className="food-grid">
          {foods.map((food) => (
            <FoodCard key={food.id} food={food} />
          ))}
        </div>
      )}
    </div>
  );
};

export default Home;
