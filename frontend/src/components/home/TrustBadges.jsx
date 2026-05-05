import React from 'react';
import { motion } from 'framer-motion';
import { FaShieldAlt, FaTruck, FaUndo, FaHeadset } from 'react-icons/fa';

const TrustBadges = () => {
  const badges = [
    {
      icon: <FaShieldAlt className="text-3xl" />,
      title: '100% BPOM & Halal',
      description: 'Produk bersertifikasi resmi',
    },
    {
      icon: <FaTruck className="text-3xl" />,
      title: 'Gratis Ongkir',
      description: 'Minimal belanja Rp 150K',
    },
    {
      icon: <FaUndo className="text-3xl" />,
      title: 'Garansi 14 Hari',
      description: 'Pengembalian dana',
    },
    {
      icon: <FaHeadset className="text-3xl" />,
      title: 'Support 24/7',
      description: 'Customer service siap membantu',
    },
  ];

  return (
    <section className="py-12 bg-white border-y border-gray-100">
      <div className="container-custom">
        <div className="grid grid-cols-2 md:grid-cols-4 gap-6">
          {badges.map((badge, index) => (
            <motion.div
              key={index}
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: index * 0.1 }}
              className="text-center"
            >
              <div className="text-primary mb-3 flex justify-center">{badge.icon}</div>
              <h3 className="font-semibold text-gray-800 mb-1">{badge.title}</h3>
              <p className="text-sm text-gray-500">{badge.description}</p>
            </motion.div>
          ))}
        </div>
      </div>
    </section>
  );
};

export default TrustBadges;