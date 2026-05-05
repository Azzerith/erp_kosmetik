import React from 'react';
import { motion } from 'framer-motion';
import { Link } from 'react-router-dom';
import { FaArrowRight, FaPlay } from 'react-icons/fa';
import { trendingKeywords } from '../../services/dummyData';

const HeroSection = () => {
  return (
    <section className="relative bg-gradient-to-r from-primary/10 via-secondary/5 to-primary/10 overflow-hidden">
      {/* Background Pattern */}
      <div className="absolute inset-0 opacity-10">
        <div className="absolute top-20 left-10 w-72 h-72 bg-primary rounded-full blur-3xl"></div>
        <div className="absolute bottom-20 right-10 w-96 h-96 bg-secondary rounded-full blur-3xl"></div>
      </div>

      <div className="container-custom py-16 md:py-24 relative">
        <div className="grid md:grid-cols-2 gap-12 items-center">
          {/* Left Content */}
          <motion.div
            initial={{ opacity: 0, x: -50 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.6 }}
          >
            <div className="inline-flex items-center gap-2 bg-white rounded-full px-4 py-2 shadow-sm mb-6">
              <span className="relative flex h-3 w-3">
                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
                <span className="relative inline-flex rounded-full h-3 w-3 bg-green-500"></span>
              </span>
              <span className="text-sm font-medium text-gray-700">Live Trends Update</span>
            </div>

            <h1 className="text-4xl md:text-6xl font-bold text-gray-900 leading-tight mb-6">
              Belanja Produk{' '}
              <span className="bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">
                Viral & Trending
              </span>
            </h1>

            <p className="text-lg text-gray-600 mb-8">
              Temukan produk kecantikan dan herbal terbaik yang sedang viral di TikTok, Instagram, 
              dan marketplace. Update tren real-time langsung dari sumber terpercaya.
            </p>

            <div className="flex flex-col sm:flex-row gap-4 mb-8">
              <Link to="/products" className="btn-primary inline-flex items-center justify-center gap-2">
                Belanja Sekarang
                <FaArrowRight className="group-hover:translate-x-1 transition-transform" />
              </Link>
              <button className="btn-outline inline-flex items-center justify-center gap-2">
                <FaPlay />
                Lihat Video Viral
              </button>
            </div>

            {/* Trending Keywords Marquee */}
            <div className="bg-white/80 backdrop-blur-sm rounded-xl p-4">
              <p className="text-sm font-medium text-gray-500 mb-2">🔥 Sedang Trending di TikTok:</p>
              <div className="overflow-hidden">
                <div className="flex gap-4 animate-marquee whitespace-nowrap">
                  {trendingKeywords.map((item, idx) => (
                    <span
                      key={idx}
                      className="inline-flex items-center gap-1 px-3 py-1 bg-gray-100 rounded-full text-sm"
                    >
                      <span className="text-accent">#{item.keyword}</span>
                      <span className="text-green-600 text-xs">↑{item.score}%</span>
                    </span>
                  ))}
                </div>
              </div>
            </div>
          </motion.div>

          {/* Right Image */}
          <motion.div
            initial={{ opacity: 0, x: 50 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.6, delay: 0.2 }}
            className="relative"
          >
            <div className="relative rounded-2xl overflow-hidden shadow-2xl">
              <img
                src="https://picsum.photos/id/26/600/600"
                alt="Beauty Products"
                className="w-full h-auto"
              />
              <div className="absolute inset-0 bg-gradient-to-t from-black/20 to-transparent"></div>
            </div>
            
            {/* Floating Badge */}
            <motion.div
              animate={{ y: [0, -10, 0] }}
              transition={{ duration: 2, repeat: Infinity }}
              className="absolute -top-4 -right-4 bg-white rounded-xl shadow-lg p-3"
            >
              <div className="flex items-center gap-2">
                <div className="w-10 h-10 bg-primary/10 rounded-full flex items-center justify-center">
                  <span className="text-2xl">📈</span>
                </div>
                <div>
                  <p className="text-xs text-gray-500">Trend Score</p>
                  <p className="text-xl font-bold text-primary">95.5</p>
                </div>
              </div>
            </motion.div>

            {/* Floating Badge 2 */}
            <motion.div
              animate={{ y: [0, 10, 0] }}
              transition={{ duration: 2.5, repeat: Infinity, delay: 0.5 }}
              className="absolute -bottom-4 -left-4 bg-white rounded-xl shadow-lg p-3"
            >
              <div className="flex items-center gap-2">
                <div className="w-10 h-10 bg-secondary/10 rounded-full flex items-center justify-center">
                  <span className="text-2xl">⭐</span>
                </div>
                <div>
                  <p className="text-xs text-gray-500">Rated 4.8/5</p>
                  <p className="text-sm font-semibold">dari 10k+ pembeli</p>
                </div>
              </div>
            </motion.div>
          </motion.div>
        </div>
      </div>
    </section>
  );
};

export default HeroSection;