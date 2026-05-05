import { useState, useEffect, useCallback } from 'react';
import { products, trendingKeywords } from '../services/dummyData';

// Hook untuk trending products
export const useTrendingProducts = (limit = 8) => {
  const [trendingProducts, setTrendingProducts] = useState([]);
  const [loading, setLoading] = useState(true);
  
  useEffect(() => {
    // Simulasi API call
    const fetchTrending = async () => {
      setLoading(true);
      try {
        await new Promise(resolve => setTimeout(resolve, 500));
        const trending = products
          .filter(p => p.trendScore > 70)
          .sort((a, b) => b.trendScore - a.trendScore)
          .slice(0, limit);
        setTrendingProducts(trending);
      } catch (error) {
        console.error('Error fetching trending products:', error);
      } finally {
        setLoading(false);
      }
    };
    
    fetchTrending();
  }, [limit]);
  
  return { trendingProducts, loading };
};

// Hook untuk trending keywords
export const useTrendingKeywords = () => {
  const [keywords, setKeywords] = useState([]);
  const [loading, setLoading] = useState(true);
  
  useEffect(() => {
    const fetchKeywords = async () => {
      setLoading(true);
      try {
        await new Promise(resolve => setTimeout(resolve, 300));
        setKeywords(trendingKeywords);
      } catch (error) {
        console.error('Error fetching trending keywords:', error);
      } finally {
        setLoading(false);
      }
    };
    
    fetchKeywords();
  }, []);
  
  return { keywords, loading };
};

// Hook untuk trend score by product
export const useTrendScore = (productId) => {
  const [trendScore, setTrendScore] = useState(null);
  const [trendHistory, setTrendHistory] = useState([]);
  const [loading, setLoading] = useState(true);
  
  useEffect(() => {
    const fetchTrendScore = async () => {
      setLoading(true);
      try {
        await new Promise(resolve => setTimeout(resolve, 500));
        const product = products.find(p => p.id === productId);
        setTrendScore(product?.trendScore || 0);
        
        // Generate dummy history data
        const history = [];
        for (let i = 30; i >= 0; i--) {
          history.push({
            date: new Date(Date.now() - i * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
            score: Math.max(0, (product?.trendScore || 50) + (Math.random() * 20 - 10)),
          });
        }
        setTrendHistory(history);
      } catch (error) {
        console.error('Error fetching trend score:', error);
      } finally {
        setLoading(false);
      }
    };
    
    if (productId) {
      fetchTrendScore();
    }
  }, [productId]);
  
  return { trendScore, trendHistory, loading };
};

// Hook untuk trending by category
export const useTrendingByCategory = (categoryId, limit = 4) => {
  const [products, setProducts] = useState([]);
  const [loading, setLoading] = useState(true);
  
  useEffect(() => {
    const fetchCategoryTrending = async () => {
      setLoading(true);
      try {
        await new Promise(resolve => setTimeout(resolve, 500));
        // In real app, filter by category from API
        const filtered = products
          .filter(p => p.categoryId === categoryId)
          .sort((a, b) => b.trendScore - a.trendScore)
          .slice(0, limit);
        setProducts(filtered);
      } catch (error) {
        console.error('Error fetching category trending:', error);
      } finally {
        setLoading(false);
      }
    };
    
    if (categoryId) {
      fetchCategoryTrending();
    }
  }, [categoryId, limit]);
  
  return { products, loading };
};

// Hook untuk real-time trend update (simulasi)
export const useRealTimeTrend = (productId) => {
  const [currentScore, setCurrentScore] = useState(0);
  const [isIncreasing, setIsIncreasing] = useState(true);
  
  useEffect(() => {
    const product = products.find(p => p.id === productId);
    setCurrentScore(product?.trendScore || 0);
    
    // Simulasi real-time updates every 30 seconds
    const interval = setInterval(() => {
      setCurrentScore(prev => {
        const change = (Math.random() * 6) - 3; // -3 to +3
        const newScore = Math.min(100, Math.max(0, prev + change));
        setIsIncreasing(change > 0);
        return newScore;
      });
    }, 30000);
    
    return () => clearInterval(interval);
  }, [productId]);
  
  return { 
    currentScore, 
    isIncreasing,
    trendDirection: isIncreasing ? 'up' : 'down',
    trendPercentage: isIncreasing ? '+2.5%' : '-1.2%',
  };
};