import React from 'react';
import { Link } from 'react-router-dom';
import { motion } from 'framer-motion';
import { FaTrash, FaPlus, FaMinus, FaShoppingBag, FaArrowLeft } from 'react-icons/fa';
import { useCart } from '../contexts/CartContext';
import { formatPrice } from '../utils/formatPrice';

const CartPage = () => {
  const { cartItems, removeFromCart, updateQuantity, getCartTotal, clearCart } = useCart();
  
  const shippingCost = 20000;
  const discount = 0;
  const total = getCartTotal() + shippingCost - discount;

  if (cartItems.length === 0) {
    return (
      <div className="container-custom py-20 text-center">
        <div className="text-6xl mb-4">🛒</div>
        <h2 className="text-2xl font-semibold mb-2">Keranjang Belanja Kosong</h2>
        <p className="text-gray-500 mb-6">Yuk, mulai belanja produk favorit Anda!</p>
        <Link to="/products" className="btn-primary inline-block">
          Mulai Belanja
        </Link>
      </div>
    );
  }

  return (
    <div className="container-custom py-8">
      <div className="flex items-center gap-4 mb-6">
        <Link to="/products" className="text-gray-500 hover:text-primary">
          <FaArrowLeft />
        </Link>
        <h1 className="text-2xl font-bold">Keranjang Belanja</h1>
      </div>

      <div className="grid lg:grid-cols-3 gap-8">
        {/* Cart Items */}
        <div className="lg:col-span-2 space-y-4">
          {cartItems.map((item) => (
            <motion.div
              key={`${item.id}-${item.variantId}`}
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              className="bg-white rounded-xl shadow-sm p-4 flex gap-4"
            >
              <img
                src={item.image || 'https://picsum.photos/100/100'}
                alt={item.name}
                className="w-24 h-24 object-cover rounded-lg"
              />
              
              <div className="flex-1">
                <div className="flex justify-between mb-1">
                  <h3 className="font-semibold text-gray-800">{item.name}</h3>
                  <button
                    onClick={() => removeFromCart(item.id, item.variantId)}
                    className="text-red-500 hover:text-red-700"
                  >
                    <FaTrash />
                  </button>
                </div>
                
                {item.variantName && (
                  <p className="text-sm text-gray-500 mb-2">Varian: {item.variantName}</p>
                )}
                
                <div className="flex justify-between items-center mt-2">
                  <div className="flex items-center gap-3 border rounded-lg">
                    <button
                      onClick={() => updateQuantity(item.id, item.variantId, item.quantity - 1)}
                      className="p-2 px-3 hover:bg-gray-100"
                    >
                      <FaMinus className="text-xs" />
                    </button>
                    <span className="w-8 text-center">{item.quantity}</span>
                    <button
                      onClick={() => updateQuantity(item.id, item.variantId, item.quantity + 1)}
                      className="p-2 px-3 hover:bg-gray-100"
                    >
                      <FaPlus className="text-xs" />
                    </button>
                  </div>
                  <span className="font-bold text-primary">
                    {formatPrice(item.price * item.quantity)}
                  </span>
                </div>
              </div>
            </motion.div>
          ))}
          
          <div className="text-right">
            <button
              onClick={clearCart}
              className="text-red-500 hover:underline text-sm"
            >
              Hapus Semua
            </button>
          </div>
        </div>

        {/* Order Summary */}
        <div className="lg:col-span-1">
          <div className="bg-white rounded-xl shadow-sm p-6 sticky top-24">
            <h2 className="text-lg font-semibold mb-4">Ringkasan Belanja</h2>
            
            <div className="space-y-3 mb-4">
              <div className="flex justify-between text-gray-600">
                <span>Subtotal</span>
                <span>{formatPrice(getCartTotal())}</span>
              </div>
              <div className="flex justify-between text-gray-600">
                <span>Ongkos Kirim</span>
                <span>{formatPrice(shippingCost)}</span>
              </div>
              <div className="flex justify-between text-gray-600">
                <span>Diskon</span>
                <span>-{formatPrice(discount)}</span>
              </div>
              <div className="border-t pt-3">
                <div className="flex justify-between font-bold text-lg">
                  <span>Total</span>
                  <span className="text-primary">{formatPrice(total)}</span>
                </div>
              </div>
            </div>

            <Link to="/checkout" className="btn-primary w-full block text-center">
              Checkout
            </Link>
            
            <div className="mt-4 text-center text-sm text-gray-500">
              <FaShoppingBag className="inline mr-1" />
              Gratis ongkir minimal belanja Rp 150.000
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default CartPage;