import React, { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import ProductCard from '../common/ProductCard';
import LoadingSpinner, { ProductCardSkeleton } from '../common/LoadingSpinner';

const ProductList = ({ products, loading, columns = 4, showLoadMore = true }) => {
  const [visibleCount, setVisibleCount] = useState(8);
  const [filteredProducts, setFilteredProducts] = useState([]);

  useEffect(() => {
    setFilteredProducts(products);
    setVisibleCount(8);
  }, [products]);

  const columnClasses = {
    2: 'grid-cols-1 sm:grid-cols-2',
    3: 'grid-cols-1 sm:grid-cols-2 lg:grid-cols-3',
    4: 'grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4',
    5: 'grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-5',
  };

  const displayedProducts = filteredProducts.slice(0, visibleCount);
  const hasMore = visibleCount < filteredProducts.length;

  const loadMore = () => {
    setVisibleCount(prev => Math.min(prev + 8, filteredProducts.length));
  };

  if (loading) {
    return (
      <div className={columnClasses[columns]}>
        {[...Array(8)].map((_, i) => (
          <ProductCardSkeleton key={i} />
        ))}
      </div>
    );
  }

  if (filteredProducts.length === 0) {
    return (
      <div className="text-center py-12">
        <div className="text-6xl mb-4">🔍</div>
        <h3 className="text-lg font-semibold text-gray-800 mb-2">Tidak Ada Produk</h3>
        <p className="text-gray-500">Coba filter atau kata kunci yang berbeda</p>
      </div>
    );
  }

  return (
    <div>
      <div className={columnClasses[columns]}>
        <AnimatePresence mode="wait">
          {displayedProducts.map((product, index) => (
            <motion.div
              key={product.id}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -20 }}
              transition={{ delay: (index % 8) * 0.05 }}
            >
              <ProductCard product={product} />
            </motion.div>
          ))}
        </AnimatePresence>
      </div>

      {showLoadMore && hasMore && (
        <div className="text-center mt-8">
          <button
            onClick={loadMore}
            className="btn-outline px-8 py-3"
          >
            Lihat Lebih Banyak
            <span className="ml-2">↓</span>
          </button>
        </div>
      )}
    </div>
  );
};

export default ProductList;