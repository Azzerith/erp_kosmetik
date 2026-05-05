import React, { useEffect } from 'react';
import { Link } from 'react-router-dom';
import { motion } from 'framer-motion';
import { FaCheckCircle, FaCopy, FaWhatsapp, FaEnvelope } from 'react-icons/fa';
import toast from 'react-hot-toast';

const OrderSuccessPage = () => {
  const orderNumber = `ORD-${Date.now()}`;

  useEffect(() => {
    // In real app, this would track order completion
    window.scrollTo(0, 0);
  }, []);

  const copyOrderNumber = () => {
    navigator.clipboard.writeText(orderNumber);
    toast.success('Nomor pesanan disalin');
  };

  return (
    <div className="min-h-screen bg-gradient-to-b from-green-50 to-white py-12">
      <div className="container-custom max-w-2xl mx-auto">
        <motion.div
          initial={{ opacity: 0, scale: 0.9 }}
          animate={{ opacity: 1, scale: 1 }}
          className="bg-white rounded-2xl shadow-xl p-8 text-center"
        >
          <motion.div
            initial={{ scale: 0 }}
            animate={{ scale: 1 }}
            transition={{ type: 'spring', delay: 0.2 }}
            className="w-20 h-20 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-6"
          >
            <FaCheckCircle className="text-green-500 text-4xl" />
          </motion.div>

          <h1 className="text-2xl font-bold text-gray-900 mb-2">
            Pesanan Berhasil!
          </h1>
          <p className="text-gray-500 mb-6">
            Terima kasih telah berbelanja di ErpCosmetics
          </p>

          <div className="bg-gray-50 rounded-xl p-4 mb-6">
            <p className="text-sm text-gray-500 mb-1">Nomor Pesanan</p>
            <div className="flex items-center justify-center gap-2">
              <p className="text-lg font-semibold text-gray-900">{orderNumber}</p>
              <button onClick={copyOrderNumber} className="text-primary hover:text-primary/80">
                <FaCopy />
              </button>
            </div>
          </div>

          <div className="space-y-3 mb-8">
            <p className="text-gray-600">
              Email konfirmasi telah dikirim ke <strong>customer@example.com</strong>
            </p>
            <p className="text-gray-600">
              Kami akan segera memproses pesanan Anda. Status pesanan dapat dilihat di halaman pesanan.
            </p>
          </div>

          <div className="flex flex-col sm:flex-row gap-4 mb-8">
            <Link to="/products" className="btn-primary">
              Lanjutkan Belanja
            </Link>
            <Link to="/orders" className="btn-outline">
              Lihat Pesanan Saya
            </Link>
          </div>

          <div className="border-t pt-6">
            <p className="text-sm text-gray-500 mb-3">Butuh bantuan? Hubungi kami</p>
            <div className="flex justify-center gap-4">
              <a href="#" className="flex items-center gap-2 text-green-600 hover:text-green-700">
                <FaWhatsapp />
                <span>WhatsApp</span>
              </a>
              <a href="#" className="flex items-center gap-2 text-primary hover:text-primary/80">
                <FaEnvelope />
                <span>Email</span>
              </a>
            </div>
          </div>
        </motion.div>
      </div>
    </div>
  );
};

export default OrderSuccessPage;