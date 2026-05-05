import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { motion, AnimatePresence } from 'framer-motion';
import { FaSearch, FaShoppingCart, FaUser, FaBars, FaTimes } from 'react-icons/fa';
import { useCart } from '../../contexts/CartContext';
import CartDrawer from '../cart/CartDrawer';

const Header = () => {
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const { getCartCount, setIsCartOpen } = useCart();
  const navigate = useNavigate();

  const handleSearch = (e) => {
    e.preventDefault();
    if (searchQuery.trim()) {
      navigate(`/products?search=${searchQuery}`);
      setSearchQuery('');
    }
  };

  return (
    <>
      <header className="sticky top-0 z-50 bg-white shadow-sm">
        <div className="container-custom">
          <div className="flex items-center justify-between py-4">
            {/* Logo */}
            <Link to="/" className="flex items-center gap-2">
              <div className="w-10 h-10 bg-gradient-to-r from-primary to-secondary rounded-full flex items-center justify-center">
                <span className="text-white font-bold text-xl">E</span>
              </div>
              <span className="font-heading font-bold text-xl text-gray-800">
                Erp<span className="text-primary">Cosmetics</span>
              </span>
            </Link>

            {/* Desktop Navigation */}
            <nav className="hidden md:flex items-center gap-8">
              <Link to="/" className="text-gray-600 hover:text-primary transition-colors">
                Beranda
              </Link>
              <Link to="/products" className="text-gray-600 hover:text-primary transition-colors">
                Produk
              </Link>
              <Link to="/trending" className="text-gray-600 hover:text-primary transition-colors flex items-center gap-1">
                <span className="text-accent">🔥</span>
                Trending
              </Link>
              <Link to="/flash-sale" className="text-gray-600 hover:text-primary transition-colors">
                Flash Sale
              </Link>
            </nav>

            {/* Search Bar */}
            <form onSubmit={handleSearch} className="hidden md:flex items-center bg-gray-100 rounded-full px-4 py-2 w-96">
              <input
                type="text"
                placeholder="Cari produk..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="bg-transparent flex-1 outline-none text-sm"
              />
              <button type="submit" className="text-gray-500 hover:text-primary">
                <FaSearch />
              </button>
            </form>

            {/* Actions */}
            <div className="flex items-center gap-4">
              <button
                onClick={() => setIsCartOpen(true)}
                className="relative text-gray-600 hover:text-primary transition-colors"
              >
                <FaShoppingCart className="text-xl" />
                {getCartCount() > 0 && (
                  <span className="absolute -top-2 -right-2 bg-primary text-white text-xs rounded-full w-5 h-5 flex items-center justify-center">
                    {getCartCount()}
                  </span>
                )}
              </button>
              
              <Link to="/login" className="hidden md:flex items-center gap-2 text-gray-600 hover:text-primary transition-colors">
                <FaUser />
                <span>Masuk</span>
              </Link>

              {/* Mobile Menu Button */}
              <button
                onClick={() => setIsMenuOpen(!isMenuOpen)}
                className="md:hidden text-gray-600"
              >
                {isMenuOpen ? <FaTimes className="text-xl" /> : <FaBars className="text-xl" />}
              </button>
            </div>
          </div>
        </div>

        {/* Mobile Menu */}
        <AnimatePresence>
          {isMenuOpen && (
            <motion.div
              initial={{ opacity: 0, height: 0 }}
              animate={{ opacity: 1, height: 'auto' }}
              exit={{ opacity: 0, height: 0 }}
              className="md:hidden bg-white border-t"
            >
              <div className="container-custom py-4">
                <div className="flex flex-col gap-4">
                  <Link to="/" className="text-gray-600 py-2" onClick={() => setIsMenuOpen(false)}>
                    Beranda
                  </Link>
                  <Link to="/products" className="text-gray-600 py-2" onClick={() => setIsMenuOpen(false)}>
                    Produk
                  </Link>
                  <Link to="/trending" className="text-gray-600 py-2" onClick={() => setIsMenuOpen(false)}>
                    Trending
                  </Link>
                  <Link to="/flash-sale" className="text-gray-600 py-2" onClick={() => setIsMenuOpen(false)}>
                    Flash Sale
                  </Link>
                  <Link to="/login" className="text-gray-600 py-2" onClick={() => setIsMenuOpen(false)}>
                    Masuk
                  </Link>
                  
                  {/* Mobile Search */}
                  <form onSubmit={handleSearch} className="flex items-center bg-gray-100 rounded-full px-4 py-2">
                    <input
                      type="text"
                      placeholder="Cari produk..."
                      value={searchQuery}
                      onChange={(e) => setSearchQuery(e.target.value)}
                      className="bg-transparent flex-1 outline-none text-sm"
                    />
                    <button type="submit" className="text-gray-500">
                      <FaSearch />
                    </button>
                  </form>
                </div>
              </div>
            </motion.div>
          )}
        </AnimatePresence>
      </header>

      <CartDrawer />
    </>
  );
};

export default Header;