import React from 'react';
import { Swiper, SwiperSlide } from 'swiper/react';
import { Navigation, Autoplay } from 'swiper/modules';
import { motion } from 'framer-motion';
import ProductCard from '../common/ProductCard';
import { products } from '../../services/dummyData';

import 'swiper/css';
import 'swiper/css/navigation';

const TrendingProducts = () => {
  const trendingProducts = products.filter(p => p.trendBadge === 'viral' || p.trendBadge === 'trending').slice(0, 8);

  return (
    <section className="py-16 bg-gradient-to-b from-white to-gray-50">
      <div className="container-custom">
        {/* Section Header */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          className="text-center mb-12"
        >
          <div className="inline-flex items-center gap-2 bg-primary/10 px-4 py-2 rounded-full mb-4">
            <span className="text-primary text-lg">📈</span>
            <span className="text-primary font-semibold">Real-time Trend</span>
          </div>
          <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
            Sedang Viral Hari Ini 🔥
          </h2>
          <p className="text-gray-600 max-w-2xl mx-auto">
            Update produk yang sedang trending di media sosial dan marketplace. 
            Data diperbarui setiap 6 jam berdasarkan Google Trends & TikTok.
          </p>
        </motion.div>

        {/* Products Slider */}
        <Swiper
          modules={[Navigation, Autoplay]}
          spaceBetween={24}
          slidesPerView={1}
          navigation
          autoplay={{ delay: 5000, disableOnInteraction: false }}
          breakpoints={{
            640: { slidesPerView: 2 },
            768: { slidesPerView: 3 },
            1024: { slidesPerView: 4 },
          }}
          className="trending-slider"
        >
          {trendingProducts.map((product) => (
            <SwiperSlide key={product.id}>
              <ProductCard product={product} />
            </SwiperSlide>
          ))}
        </Swiper>
      </div>
    </section>
  );
};

export default TrendingProducts;