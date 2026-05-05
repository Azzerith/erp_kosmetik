import React, { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { FaTruck, FaCreditCard, FaUniversity, FaQrcode, FaWallet, FaArrowLeft, FaCheck } from 'react-icons/fa';

const CheckoutForm = ({ onSubmit, onBack, isLoading }) => {
  const [step, setStep] = useState(1);
  const [shippingAddress, setShippingAddress] = useState({
    name: '',
    phone: '',
    address: '',
    city: '',
    province: '',
    postalCode: '',
    note: '',
  });
  const [paymentMethod, setPaymentMethod] = useState('');

  const handleAddressSubmit = (e) => {
    e.preventDefault();
    if (validateAddress()) {
      setStep(2);
      window.scrollTo({ top: 0, behavior: 'smooth' });
    }
  };

  const validateAddress = () => {
    const required = ['name', 'phone', 'address', 'city', 'province', 'postalCode'];
    const missing = required.filter(field => !shippingAddress[field]);
    
    if (missing.length > 0) {
      alert('Harap lengkapi data alamat pengiriman');
      return false;
    }
    return true;
  };

  const handlePayment = () => {
    if (!paymentMethod) {
      alert('Pilih metode pembayaran terlebih dahulu');
      return;
    }
    onSubmit({ shippingAddress, paymentMethod });
  };

  const paymentMethods = [
    { id: 'qris', name: 'QRIS', icon: <FaQrcode className="text-2xl" />, description: 'Scan QRIS via GoPay, OVO, Dana, LinkAja' },
    { id: 'bank_transfer', name: 'Transfer Bank', icon: <FaUniversity className="text-2xl" />, description: 'BCA, Mandiri, BNI, BRI, CIMB Niaga' },
    { id: 'credit_card', name: 'Kartu Kredit', icon: <FaCreditCard className="text-2xl" />, description: 'Visa, Mastercard, JCB' },
    { id: 'ewallet', name: 'E-Wallet', icon: <FaWallet className="text-2xl" />, description: 'GoPay, OVO, Dana, ShopeePay' },
  ];

  return (
    <div className="bg-white rounded-xl shadow-sm overflow-hidden">
      {/* Steps Indicator */}
      <div className="bg-gray-50 px-6 py-4 border-b">
        <div className="flex items-center justify-center">
          {[1, 2].map((s) => (
            <React.Fragment key={s}>
              <div className="flex items-center">
                <div
                  className={`w-10 h-10 rounded-full flex items-center justify-center font-semibold transition-all ${
                    step >= s
                      ? 'bg-primary text-white shadow-lg'
                      : 'bg-gray-200 text-gray-500'
                  }`}
                >
                  {step > s ? <FaCheck /> : s}
                </div>
                <span
                  className={`ml-2 text-sm font-medium ${
                    step >= s ? 'text-primary' : 'text-gray-500'
                  }`}
                >
                  {s === 1 ? 'Alamat Pengiriman' : 'Pembayaran'}
                </span>
              </div>
              {s < 2 && (
                <div
                  className={`w-16 h-0.5 mx-4 ${
                    step > s ? 'bg-primary' : 'bg-gray-200'
                  }`}
                />
              )}
            </React.Fragment>
          ))}
        </div>
      </div>

      <div className="p-6">
        <AnimatePresence mode="wait">
          {/* Step 1: Shipping Address */}
          {step === 1 && (
            <motion.div
              key="shipping"
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              exit={{ opacity: 0, x: 20 }}
            >
              <div className="flex items-center gap-2 mb-4">
                <FaTruck className="text-primary" />
                <h2 className="text-lg font-semibold">Informasi Pengiriman</h2>
              </div>

              <form onSubmit={handleAddressSubmit} className="space-y-4">
                <div className="grid md:grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Nama Penerima <span className="text-red-500">*</span>
                    </label>
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
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      No. Telepon <span className="text-red-500">*</span>
                    </label>
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
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Alamat Lengkap <span className="text-red-500">*</span>
                  </label>
                  <textarea
                    required
                    value={shippingAddress.address}
                    onChange={(e) => setShippingAddress({ ...shippingAddress, address: e.target.value })}
                    className="input-field"
                    rows="3"
                    placeholder="Nama jalan, nomor rumah, RT/RW, nama gedung, dll."
                  />
                </div>

                <div className="grid md:grid-cols-3 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Provinsi <span className="text-red-500">*</span>
                    </label>
                    <input
                      type="text"
                      required
                      value={shippingAddress.province}
                      onChange={(e) => setShippingAddress({ ...shippingAddress, province: e.target.value })}
                      className="input-field"
                      placeholder="Provinsi"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Kota/Kabupaten <span className="text-red-500">*</span>
                    </label>
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
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Kode Pos <span className="text-red-500">*</span>
                    </label>
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

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Catatan untuk Kurir (Opsional)
                  </label>
                  <input
                    type="text"
                    value={shippingAddress.note}
                    onChange={(e) => setShippingAddress({ ...shippingAddress, note: e.target.value })}
                    className="input-field"
                    placeholder="Contoh: Tolong diletakkan di depan pintu"
                  />
                </div>

                <button type="submit" className="btn-primary w-full">
                  Lanjut ke Pembayaran
                </button>
              </form>
            </motion.div>
          )}

          {/* Step 2: Payment Method */}
          {step === 2 && (
            <motion.div
              key="payment"
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              exit={{ opacity: 0, x: 20 }}
            >
              <div className="flex items-center gap-2 mb-4">
                <FaCreditCard className="text-primary" />
                <h2 className="text-lg font-semibold">Metode Pembayaran</h2>
              </div>

              <div className="space-y-3 mb-6">
                {paymentMethods.map((method) => (
                  <label
                    key={method.id}
                    className={`flex items-start gap-4 p-4 border rounded-xl cursor-pointer transition-all ${
                      paymentMethod === method.id
                        ? 'border-primary bg-primary/5 shadow-sm'
                        : 'border-gray-200 hover:border-primary hover:bg-gray-50'
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
                      <div className="flex items-center gap-3">
                        <span className="text-gray-600">{method.icon}</span>
                        <div>
                          <span className="font-semibold">{method.name}</span>
                          <p className="text-sm text-gray-500 mt-0.5">{method.description}</p>
                        </div>
                      </div>
                    </div>
                    {paymentMethod === method.id && (
                      <div className="text-primary">
                        <FaCheck />
                      </div>
                    )}
                  </label>
                ))}
              </div>

              <div className="flex gap-4">
                <button onClick={onBack || (() => setStep(1))} className="btn-outline flex-1">
                  <FaArrowLeft className="inline mr-2" />
                  Kembali
                </button>
                <button onClick={handlePayment} className="btn-primary flex-1" disabled={isLoading}>
                  {isLoading ? 'Memproses...' : 'Bayar Sekarang'}
                </button>
              </div>
            </motion.div>
          )}
        </AnimatePresence>
      </div>
    </div>
  );
};

export default CheckoutForm;