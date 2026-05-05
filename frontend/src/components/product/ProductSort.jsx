import React from 'react';

const ProductSort = ({ sortBy, onSortChange }) => {
  const sortOptions = [
    { value: 'trending', label: '🔥 Trending' },
    { value: 'best_seller', label: '⭐ Best Seller' },
    { value: 'newest', label: '✨ Terbaru' },
    { value: 'price_asc', label: 'Harga: Rendah ke Tinggi' },
    { value: 'price_desc', label: 'Harga: Tinggi ke Rendah' },
    { value: 'rating', label: 'Rating Tertinggi' },
  ];

  return (
    <div className="flex items-center gap-2">
      <span className="text-sm text-gray-500 hidden sm:inline">Urutkan:</span>
      <select
        value={sortBy}
        onChange={(e) => onSortChange(e.target.value)}
        className="px-4 py-2 border rounded-lg text-sm focus:ring-primary focus:border-primary outline-none"
      >
        {sortOptions.map(option => (
          <option key={option.value} value={option.value}>
            {option.label}
          </option>
        ))}
      </select>
    </div>
  );
};

export default ProductSort;