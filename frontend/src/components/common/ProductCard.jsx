import React from 'react';
import { Link } from 'react-router-dom';
import { motion } from 'framer-motion';
import { FaStar, FaShoppingCart } from 'react-icons/fa';
import TrendBadge from './TrendBadge';
import { formatPrice, formatDiscount } from '../../utils/formatPrice';
import { useCart } from '../../contexts/CartContext';

const ProductCard = ({ product, variant = 'default' }) => {
  const { addToCart } = useCart();
  const displayPrice = product.salePrice || product.basePrice;
  const discount = formatDiscount(product.basePrice, product.salePrice);

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      whileHover={{ y: -8 }}
      transition={{ duration: 0.3 }}
      className="group relative bg-white rounded-2xl overflow-hidden shadow-sm hover:shadow-xl transition-all duration-300"
    >
      {/* Image Container */}
      <Link to={`/product/${product.slug}`} className="block relative overflow-hidden bg-gray-100 aspect-[4/5]">
        <img
          src={product.images?.[0] || 'https://picsum.photos/400/500'}
          alt={product.name}
          className="w-full h-full object-cover group-hover:scale-110 transition-transform duration-500"
        />
        
        {/* Trend Badge */}
        <TrendBadge type={product.trendBadge} score={product.trendScore} />
        
        {/* Discount Badge */}
        {discount > 0 && (
          <div className="absolute top-2 right-2 bg-red-500 text-white px-2 py-1 rounded-lg text-xs font-bold">
            -{discount}%
          </div>
        )}
        
        {/* Quick Actions */}
        <div className="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/60 to-transparent p-4 translate-y-full group-hover:translate-y-0 transition-transform duration-300">
          <button
            onClick={() => addToCart(product)}
            className="w-full bg-white text-primary py-2 rounded-full font-semibold hover:bg-primary hover:text-white transition-colors flex items-center justify-center gap-2"
          >
            <FaShoppingCart />
            Tambah ke Keranjang
          </button>
        </div>
      </Link>

      {/* Product Info */}
      <div className="p-4">
        {/* Brand & Category */}
        <div className="text-xs text-gray-500 mb-1">
          {product.brand} • {product.category}
        </div>
        
        {/* Name */}
        <Link to={`/product/${product.slug}`}>
          <h3 className="font-semibold text-gray-800 mb-2 line-clamp-2 hover:text-primary transition-colors">
            {product.name}
          </h3>
        </Link>
        
        {/* Rating */}
        <div className="flex items-center gap-2 mb-2">
          <div className="flex items-center text-yellow-400">
            <FaStar className="fill-current" />
            <span className="text-sm font-semibold ml-1 text-gray-700">{product.rating}</span>
          </div>
          <span className="text-xs text-gray-400">({product.totalReviews} ulasan)</span>
        </div>
        
        {/* Price */}
        <div className="flex items-center gap-2">
          <span className="text-lg font-bold text-primary">{formatPrice(displayPrice)}</span>
          {product.salePrice && (
            <span className="text-sm text-gray-400 line-through">{formatPrice(product.basePrice)}</span>
          )}
        </div>
        
        {/* Sold Count */}
        <div className="text-xs text-gray-500 mt-2 flex items-center gap-1">
          <span>Terjual {product.totalSold}+</span>
        </div>
      </div>
    </motion.div>
  );
};

export default ProductCard;