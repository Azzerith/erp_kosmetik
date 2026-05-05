import React, { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { useSearchParams } from 'react-router-dom';
import ProductCard from '../components/common/ProductCard';
import ProductFilter from '../components/product/ProductFilter';
import ProductSort from '../components/product/ProductSort';
import { products } from '../services/dummyData';

const ProductsPage = () => {
  const [searchParams] = useSearchParams();
  const [filteredProducts, setFilteredProducts] = useState(products);
  const [sortBy, setSortBy] = useState('trending');

  // Filter by search query from URL
  useEffect(() => {
    const searchQuery = searchParams.get('search');
    const category = searchParams.get('category');
    
    let filtered = [...products];
    
    if (searchQuery) {
      filtered = filtered.filter(p => 
        p.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        p.brand.toLowerCase().includes(searchQuery.toLowerCase())
      );
    }
    
    if (category) {
      filtered = filtered.filter(p => p.slug === category);
    }
    
    setFilteredProducts(filtered);
  }, [searchParams]);

  // Sort products
  const getSortedProducts = () => {
    const sorted = [...filteredProducts];
    switch (sortBy) {
      case 'trending':
        return sorted.sort((a, b) => b.trendScore - a.trendScore);
      case 'best_seller':
        return sorted.sort((a, b) => b.totalSold - a.totalSold);
      case 'newest':
        return sorted.sort((a, b) => b.id - a.id);
      case 'price_asc':
        return sorted.sort((a, b) => (a.salePrice || a.basePrice) - (b.salePrice || b.basePrice));
      case 'price_desc':
        return sorted.sort((a, b) => (b.salePrice || b.basePrice) - (a.salePrice || a.basePrice));
      case 'rating':
        return sorted.sort((a, b) => b.rating - a.rating);
      default:
        return sorted;
    }
  };

  const displayProducts = getSortedProducts();

  return (
    <div className="container-custom py-8">
      <div className="flex flex-col lg:flex-row gap-8">
        {/* Sidebar Filter */}
        <div className="lg:w-1/4">
          <ProductFilter onFilter={setFilteredProducts} />
        </div>

        {/* Product Grid */}
        <div className="lg:w-3/4">
          <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
            <h1 className="text-2xl font-bold text-gray-900">
              Semua Produk
              <span className="text-sm font-normal text-gray-500 ml-2">
                ({displayProducts.length} produk)
              </span>
            </h1>
            <ProductSort sortBy={sortBy} onSortChange={setSortBy} />
          </div>

          {displayProducts.length === 0 ? (
            <div className="text-center py-12">
              <div className="text-6xl mb-4">🔍</div>
              <p className="text-gray-500">Tidak ada produk yang ditemukan</p>
            </div>
          ) : (
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
              {displayProducts.map((product, index) => (
                <motion.div
                  key={product.id}
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: index * 0.05 }}
                >
                  <ProductCard product={product} />
                </motion.div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default ProductsPage;