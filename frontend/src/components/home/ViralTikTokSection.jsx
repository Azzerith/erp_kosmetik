import React from 'react';
import { motion } from 'framer-motion';
import { FaTiktok, FaHeart, FaEye } from 'react-icons/fa';
import { viralTikTok, products } from '../../services/dummyData';

const ViralTikTokSection = () => {
  const viralProducts = products.filter(p => p.trendBadge === 'viral');

  return (
    <section className="py-16 bg-gradient-to-b from-pink-50 to-white">
      <div className="container-custom">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          className="text-center mb-12"
        >
          <div className="inline-flex items-center gap-2 bg-black text-white px-4 py-2 rounded-full mb-4">
            <FaTiktok className="text-xl" />
            <span className="font-semibold">Viral di TikTok</span>
          </div>
          <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
            Produk yang Lagi Viral di TikTok
          </h2>
          <p className="text-gray-600 max-w-2xl mx-auto">
            Produk-produk yang sedang ramai dibicarakan dan menjadi tantangan di TikTok Indonesia
          </p>
        </motion.div>

        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
          {viralProducts.map((product, index) => (
            <motion.div
              key={product.id}
              initial={{ opacity: 0, scale: 0.9 }}
              whileInView={{ opacity: 1, scale: 1 }}
              viewport={{ once: true }}
              transition={{ delay: index * 0.1 }}
              whileHover={{ y: -8 }}
              className="bg-white rounded-2xl overflow-hidden shadow-lg"
            >
              <div className="relative aspect-video bg-gradient-to-br from-purple-500 to-pink-500 flex items-center justify-center">
                <div className="text-center text-white">
                  <FaTiktok className="text-6xl mb-2" />
                  <p className="text-sm">Video Viral</p>
                  <p className="text-xs opacity-75">#{product.name.replace(/\s/g, '')}</p>
                </div>
                <div className="absolute bottom-2 right-2 bg-black/50 backdrop-blur rounded-full px-3 py-1 text-white text-xs flex items-center gap-2">
                  <FaEye />
                  <span>2.5M views</span>
                </div>
              </div>
              <div className="p-4">
                <h3 className="font-semibold text-gray-800 mb-2">{product.name}</h3>
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-4 text-sm text-gray-600">
                    <span className="flex items-center gap-1">
                      <FaHeart className="text-red-500" />
                      {Math.floor(Math.random() * 200) + 100}K
                    </span>
                    <span className="text-green-600">
                      ↑ {Math.floor(Math.random() * 50) + 20}% growth
                    </span>
                  </div>
                  <a
                    href={`https://tiktok.com/search?q=${product.name}`}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-primary hover:underline text-sm font-semibold"
                  >
                    Lihat Video →
                  </a>
                </div>
              </div>
            </motion.div>
          ))}
        </div>
      </div>
    </section>
  );
};

export default ViralTikTokSection;