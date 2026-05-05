import React from 'react';
import { FaFire, FaTiktok, FaChartLine } from 'react-icons/fa';

const TrendBadge = ({ type, score, className = '' }) => {
  const badges = {
    viral: {
      icon: <FaTiktok className="mr-1" />,
      text: '🔥 VIRAL TIKTOK',
      color: 'bg-gradient-to-r from-pink-500 to-purple-500',
    },
    trending: {
      icon: <FaChartLine className="mr-1" />,
      text: `📈 TRENDING ${score ? `+${score}%` : ''}`,
      color: 'bg-accent',
    },
    best_seller: {
      icon: <FaFire className="mr-1" />,
      text: '⭐ BEST SELLER',
      color: 'bg-primary',
    },
    hot: {
      icon: <FaFire className="mr-1" />,
      text: '🔥 HOT',
      color: 'bg-orange-500',
    },
    none: {
      icon: null,
      text: '',
      color: '',
    },
  };

  const badge = badges[type] || badges.none;
  
  if (type === 'none') return null;

  return (
    <div className={`trend-badge ${badge.color} text-white ${className}`}>
      <div className="flex items-center text-xs">
        {badge.icon}
        {badge.text}
      </div>
    </div>
  );
};

export default TrendBadge;