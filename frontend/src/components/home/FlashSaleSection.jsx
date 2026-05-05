import React, { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import ProductCard from '../common/ProductCard';
import { flashSaleItems } from '../../services/dummyData';

const FlashSaleSection = () => {
  const [timeLeft, setTimeLeft] = useState({
    hours: 23,
    minutes: 59,
    seconds: 59,
  });

  useEffect(() => {
    const timer = setInterval(() => {
      setTimeLeft((prev) => {
        if (prev.seconds > 0) {
          return { ...prev, seconds: prev.seconds - 1 };
        } else if (prev.minutes > 0) {
          return { ...prev, minutes: prev.minutes - 1, seconds: 59 };
        } else if (prev.hours > 0) {
          return { hours: prev.hours - 1, minutes: 59, seconds: 59 };
        }
        return prev;
      });
    }, 1000);

    return () => clearInterval(timer);
  }, []);

  return (
    <section className="py-16 bg-gradient-to-r from-red-50 to-orange-50">
      <div className="container-custom">
        <div className="bg-white rounded-2xl shadow-xl overflow-hidden">
          {/* Header */}
          <div className="bg-gradient-to-r from-red-500 to-orange-500 p-8 text-white">
            <div className="flex flex-col md:flex-row justify-between items-center gap-4">
              <div>
                <div className="flex items-center gap-2 mb-2">
                  <span className="text-3xl">⚡</span>
                  <h2 className="text-2xl md:text-3xl font-bold">Flash Sale</h2>
                </div>
                <p className="opacity-90">Diskon terbatas! Jangan sampai kehabisan</p>
              </div>
              
              <div className="text-center">
                <p className="text-sm mb-2">Berakhir dalam:</p>
                <div className="flex gap-3">
                  <div className="bg-white/20 backdrop-blur rounded-lg px-4 py-2">
                    <span className="text-2xl font-bold">{String(timeLeft.hours).padStart(2, '0')}</span>
                    <span className="text-xs ml-1">Jam</span>
                  </div>
                  <div className="bg-white/20 backdrop-blur rounded-lg px-4 py-2">
                    <span className="text-2xl font-bold">{String(timeLeft.minutes).padStart(2, '0')}</span>
                    <span className="text-xs ml-1">Menit</span>
                  </div>
                  <div className="bg-white/20 backdrop-blur rounded-lg px-4 py-2">
                    <span className="text-2xl font-bold">{String(timeLeft.seconds).padStart(2, '0')}</span>
                    <span className="text-xs ml-1">Detik</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* Products Grid */}
          <div className="p-6">
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
              {flashSaleItems.map((product, index) => (
                <motion.div
                  key={product.id}
                  initial={{ opacity: 0, y: 20 }}
                  whileInView={{ opacity: 1, y: 0 }}
                  transition={{ delay: index * 0.1 }}
                >
                  <div className="relative">
                    <div className="absolute -top-2 -left-2 bg-red-500 text-white px-3 py-1 rounded-full text-sm font-bold z-10">
                      FLASH SALE
                    </div>
                    <ProductCard product={{ ...product, salePrice: product.flashPrice, basePrice: product.originalPrice }} />
                  </div>
                </motion.div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </section>
  );
};

export default FlashSaleSection;