import React from 'react';
import HeroSection from '../components/home/HeroSection';
import TrendingProducts from '../components/home/TrendingProducts';
import CategoryShowcase from '../components/home/CategoryShowcase';
import ViralTikTokSection from '../components/home/ViralTikTokSection';
import FlashSaleSection from '../components/home/FlashSaleSection';
import BestSellerSection from '../components/home/BestSellerSection';
import TestimonialSection from '../components/home/TestimonialSection';
import TrustBadges from '../components/home/TrustBadges';

const HomePage = () => {
  return (
    <div>
      <HeroSection />
      <TrustBadges />
      <TrendingProducts />
      <CategoryShowcase />
      <FlashSaleSection />
      <ViralTikTokSection />
      <BestSellerSection />
      <TestimonialSection />
    </div>
  );
};

export default HomePage;