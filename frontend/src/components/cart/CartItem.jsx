import React from 'react';
import { motion } from 'framer-motion';
import { FaTrash, FaPlus, FaMinus } from 'react-icons/fa';
import { formatPrice } from '../../utils/formatPrice';

const CartItem = ({ item, onUpdateQuantity, onRemove }) => {
  return (
    <motion.div
      initial={{ opacity: 0, x: -20 }}
      animate={{ opacity: 1, x: 0 }}
      exit={{ opacity: 0, x: 20 }}
      className="flex gap-4 py-4 border-b last:border-0"
    >
      {/* Product Image */}
      <div className="w-24 h-24 bg-gray-100 rounded-lg overflow-hidden flex-shrink-0">
        <img
          src={item.image || 'https://picsum.photos/100/100'}
          alt={item.name}
          className="w-full h-full object-cover"
        />
      </div>

      {/* Product Info */}
      <div className="flex-1">
        <div className="flex justify-between items-start mb-1">
          <div>
            <h3 className="font-semibold text-gray-800">{item.name}</h3>
            {item.variantName && (
              <p className="text-sm text-gray-500">Varian: {item.variantName}</p>
            )}
          </div>
          <button
            onClick={() => onRemove(item.id, item.variantId)}
            className="text-red-500 hover:text-red-700 transition-colors"
            aria-label="Hapus item"
          >
            <FaTrash size={14} />
          </button>
        </div>

        {/* Price */}
        <p className="text-primary font-bold mb-2">{formatPrice(item.price)}</p>

        {/* Quantity Controls */}
        <div className="flex items-center gap-3">
          <div className="flex items-center border rounded-lg overflow-hidden">
            <button
              onClick={() => onUpdateQuantity(item.id, item.variantId, item.quantity - 1)}
              className="px-3 py-1 hover:bg-gray-100 transition-colors"
              disabled={item.quantity <= 1}
              aria-label="Kurangi jumlah"
            >
              <FaMinus size={10} />
            </button>
            <span className="w-10 text-center text-sm">{item.quantity}</span>
            <button
              onClick={() => onUpdateQuantity(item.id, item.variantId, item.quantity + 1)}
              className="px-3 py-1 hover:bg-gray-100 transition-colors"
              disabled={item.quantity >= item.maxStock}
              aria-label="Tambah jumlah"
            >
              <FaPlus size={10} />
            </button>
          </div>
          
          <span className="text-sm text-gray-500">
            Stok tersedia: {item.maxStock}
          </span>
        </div>
      </div>

      {/* Subtotal */}
      <div className="text-right">
        <p className="font-bold text-gray-800">
          {formatPrice(item.price * item.quantity)}
        </p>
        {item.quantity > 1 && (
          <p className="text-xs text-gray-400">
            @{formatPrice(item.price)}
          </p>
        )}
      </div>
    </motion.div>
  );
};

export default CartItem;