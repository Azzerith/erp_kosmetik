import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { motion } from 'framer-motion';
import { FaCreditCard, FaUniversity, FaQrcode, FaWallet, FaTruck } from 'react-icons/fa';
import { useCart } from '../contexts/CartContext';
import { formatPrice } from '../utils/formatPrice';
import toast from 'react-hot-toast';

const CheckoutPage = () => {
  const navigate = useNavigate();
  const { cartItems, getCartTotal, clearCart } = useCart();
  const [step, setStep] = useState(1);
  const [shippingAddress, setShippingAddress] = useState({
    name: '',
    phone: '',
    address: '',
    city: '',
    postalCode: '',
  });
  const [paymentMethod, setPaymentMethod] = useState('');

  const subtotal = getCartTotal();
  const shippingCost = 20000;
  const total = subtotal + shippingCost;

  const handleAddressSubmit = (e) => {
    e.preventDefault();
    setStep(2);
  };

  const handlePayment = () => {
    if (!paymentMethod) {
      toast.error('Pilih metode pembayaran terlebih dahulu');
      return;
    }
    
    toast.success('Pesanan berhasil dibuat! (Demo)');
    clearCart();
    navigate('/order-success');
  };

  if (cartItems.length === 0) {
    navigate('/cart');
    return null;
  }

  return (
    <div className="container-custom py-8">
      <h1 className="text-2xl font-bold mb-6">Checkout</h1>

      {/* Steps */}
      <div className="flex items-center justify-center mb-8">
        {[1, 2].map((s) => (
          <div key={s} className="flex items-center">
            <div
              className={`w-10 h-10 rounded-full flex items-center justify-center font-semibold ${
                step >= s ? 'bg-primary text-white' : 'bg-gray-200 text-gray-500'
              }`}
            >
              {s}
            </div>
            {s < 2 && (
              <div className={`w-16 h-0.5 ${step > s ? 'bg-primary' : 'bg-gray-200'}`} />
            )}
          </div>
        ))}
      </div>

      <div className="grid lg:grid-cols-3 gap-8">
        {/* Form Section */}
        <div className="lg:col-span-2">
          {step === 1 && (
            <motion.div
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              className="bg-white rounded-xl shadow-sm p-6"
            >
              <h2 className="text-lg font-semibold mb-4 flex items-center gap-2">
                <FaTruck className="text-primary" />
                Alamat Pengiriman
              </h2>
              
              <form onSubmit={handleAddressSubmit} className="space-y-4">
                <div className="grid md:grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium mb-1">Nama Penerima</label>
                    <input
                      type="text"
                      required
                      value={shippingAddress.name}
                      onChange={(e) => setShippingAddress({ ...shippingAddress, name: e.target.value })}
                      className="input-field"
                      placeholder="Nama lengkap"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium mb-1">No. Telepon</label>
                    <input
                      type="tel"
                      required
                      value={shippingAddress.phone}
                      onChange={(e) => setShippingAddress({ ...shippingAddress, phone: e.target.value })}
                      className="input-field"
                      placeholder="08123456789"
                    />
                  </div>
                </div>
                
                <div>
                  <label className="block text-sm font-medium mb-1">Alamat Lengkap</label>
                  <textarea
                    required
                    value={shippingAddress.address}
                    onChange={(e) => setShippingAddress({ ...shippingAddress, address: e.target.value })}
                    className="input-field"
                    rows="3"
                    placeholder="Jalan, Gedung, No. Rumah, RT/RW"
                  />
                </div>
                
                <div className="grid md:grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium mb-1">Kota</label>
                    <input
                      type="text"
                      required
                      value={shippingAddress.city}
                      onChange={(e) => setShippingAddress({ ...shippingAddress, city: e.target.value })}
                      className="input-field"
                      placeholder="Kota/Kabupaten"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium mb-1">Kode Pos</label>
                    <input
                      type="text"
                      required
                      value={shippingAddress.postalCode}
                      onChange={(e) => setShippingAddress({ ...shippingAddress, postalCode: e.target.value })}
                      className="input-field"
                      placeholder="Kode pos"
                    />
                  </div>
                </div>
                
                <button type="submit" className="btn-primary w-full">
                  Lanjut ke Pembayaran
                </button>
              </form>
            </motion.div>
          )}

          {step === 2 && (
            <motion.div
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              className="bg-white rounded-xl shadow-sm p-6"
            >
              <h2 className="text-lg font-semibold mb-4">Metode Pembayaran</h2>
              
              <div className="space-y-3">
                {[
                  { id: 'qris', name: 'QRIS', icon: <FaQrcode />, desc: 'Scan QRIS via GoPay, OVO, Dana, dll' },
                  { id: 'bank', name: 'Transfer Bank', icon: <FaUniversity />, desc: 'BCA, Mandiri, BNI, BRI' },
                  { id: 'credit', name: 'Kartu Kredit', icon: <FaCreditCard />, desc: 'Visa, Mastercard' },
                  { id: 'ewallet', name: 'E-Wallet', icon: <FaWallet />, desc: 'GoPay, OVO, Dana, ShopeePay' },
                ].map((method) => (
                  <label
                    key={method.id}
                    className={`flex items-start gap-4 p-4 border rounded-xl cursor-pointer transition-all ${
                      paymentMethod === method.id
                        ? 'border-primary bg-primary/5'
                        : 'border-gray-200 hover:border-primary'
                    }`}
                  >
                    <input
                      type="radio"
                      name="payment"
                      value={method.id}
                      checked={paymentMethod === method.id}
                      onChange={(e) => setPaymentMethod(e.target.value)}
                      className="mt-1 text-primary focus:ring-primary"
                    />
                    <div className="flex-1">
                      <div className="flex items-center gap-2">
                        <span className="text-xl">{method.icon}</span>
                        <span className="font-semibold">{method.name}</span>
                      </div>
                      <p className="text-sm text-gray-500 mt-1">{method.desc}</p>
                    </div>
                  </label>
                ))}
              </div>
              
              <div className="flex gap-4 mt-6">
                <button
                  onClick={() => setStep(1)}
                  className="btn-outline flex-1"
                >
                  Kembali
                </button>
                <button onClick={handlePayment} className="btn-primary flex-1">
                  Bayar Sekarang
                </button>
              </div>
            </motion.div>
          )}
        </div>

        {/* Order Summary */}
        <div className="lg:col-span-1">
          <div className="bg-white rounded-xl shadow-sm p-6 sticky top-24">
            <h2 className="text-lg font-semibold mb-4">Pesanan Anda</h2>
            
            <div className="space-y-3 mb-4 max-h-80 overflow-y-auto">
              {cartItems.map((item) => (
                <div key={`${item.id}-${item.variantId}`} className="flex gap-3">
                  <img
                    src={item.image || 'https://picsum.photos/50/50'}
                    alt={item.name}
                    className="w-12 h-12 object-cover rounded"
                  />
                  <div className="flex-1">
                    <p className="text-sm font-medium">{item.name}</p>
                    <p className="text-xs text-gray-500">{item.quantity} x {formatPrice(item.price)}</p>
                  </div>
                  <span className="text-sm font-semibold">{formatPrice(item.price * item.quantity)}</span>
                </div>
              ))}
            </div>
            
            <div className="space-y-2 pt-4 border-t">
              <div className="flex justify-between text-gray-600">
                <span>Subtotal</span>
                <span>{formatPrice(subtotal)}</span>
              </div>
              <div className="flex justify-between text-gray-600">
                <span>Ongkos Kirim</span>
                <span>{formatPrice(shippingCost)}</span>
              </div>
              <div className="flex justify-between font-bold text-lg pt-2 border-t">
                <span>Total</span>
                <span className="text-primary">{formatPrice(total)}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default CheckoutPage;