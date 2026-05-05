import React from 'react';
import { motion } from 'framer-motion';
import ProductCard from '../common/ProductCard';
import { products } from '../../services/dummyData';

const BestSellerSection = () => {
  const bestSellers = products.filter(p => p.trendBadge === 'best_seller').slice(0, 4);

  return (
    <section className="py-16 bg-white">
      <div className="container-custom">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          className="text-center mb-12"
        >
          <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
            Best Seller ⭐
          </h2>
          <p className="text-gray-600 max-w-2xl mx-auto">
            Produk favorit yang paling banyak dibeli oleh pelanggan setia kami
          </p>
        </motion.div>

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
          {bestSellers.map((product, index) => (
            <motion.div
              key={product.id}
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: index * 0.1 }}
            >
              <ProductCard product={product} />
            </motion.div>
          ))}
        </div>
      </div>
    </section>
  );
};

export default BestSellerSection;