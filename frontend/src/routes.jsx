import React from 'react';
import { createBrowserRouter, Navigate } from 'react-router-dom';
import Layout from './components/layout/Layout';
import HomePage from './pages/HomePage';
import ProductsPage from './pages/ProductsPage';
import ProductDetailPage from './pages/ProductDetailPage';
import CartPage from './pages/CartPage';
import CheckoutPage from './pages/CheckoutPage';
import LoginPage from './pages/LoginPage';
import RegisterPage from './pages/RegisterPage';
import OrderSuccessPage from './pages/OrderSuccessPage';

// Import protected route component
const ProtectedRoute = ({ children, requireAuth = true, requireAdmin = false }) => {
  const token = localStorage.getItem('token');
  const user = JSON.parse(localStorage.getItem('user') || '{}');
  
  if (requireAuth && !token) {
    return <Navigate to="/login" replace />;
  }
  
  if (requireAdmin && user.role !== 'admin' && user.role !== 'super_admin') {
    return <Navigate to="/" replace />;
  }
  
  return children;
};

// Route configuration
export const router = createBrowserRouter([
  {
    path: '/',
    element: <Layout />,
    children: [
      // Public routes
      {
        index: true,
        element: <HomePage />,
      },
      {
        path: 'products',
        element: <ProductsPage />,
      },
      {
        path: 'product/:slug',
        element: <ProductDetailPage />,
      },
      {
        path: 'cart',
        element: <CartPage />,
      },
      {
        path: 'login',
        element: <LoginPage />,
      },
      {
        path: 'register',
        element: <RegisterPage />,
      },
      
      // Protected routes (require authentication)
      {
        path: 'checkout',
        element: (
          <ProtectedRoute requireAuth>
            <CheckoutPage />
          </ProtectedRoute>
        ),
      },
      {
        path: 'order-success',
        element: (
          <ProtectedRoute requireAuth>
            <OrderSuccessPage />
          </ProtectedRoute>
        ),
      },
      {
        path: 'profile',
        element: (
          <ProtectedRoute requireAuth>
            <div>Profile Page (Coming Soon)</div>
          </ProtectedRoute>
        ),
      },
      {
        path: 'orders',
        element: (
          <ProtectedRoute requireAuth>
            <div>Orders Page (Coming Soon)</div>
          </ProtectedRoute>
        ),
      },
      {
        path: 'wishlist',
        element: (
          <ProtectedRoute requireAuth>
            <div>Wishlist Page (Coming Soon)</div>
          </ProtectedRoute>
        ),
      },
      
      // Admin routes
      {
        path: 'admin',
        element: (
          <ProtectedRoute requireAuth requireAdmin>
            <div>Admin Dashboard (Coming Soon)</div>
          </ProtectedRoute>
        ),
      },
      {
        path: 'admin/products',
        element: (
          <ProtectedRoute requireAuth requireAdmin>
            <div>Admin Products (Coming Soon)</div>
          </ProtectedRoute>
        ),
      },
      {
        path: 'admin/orders',
        element: (
          <ProtectedRoute requireAuth requireAdmin>
            <div>Admin Orders (Coming Soon)</div>
          </ProtectedRoute>
        ),
      },
      {
        path: 'admin/inventory',
        element: (
          <ProtectedRoute requireAuth requireAdmin>
            <div>Admin Inventory (Coming Soon)</div>
          </ProtectedRoute>
        ),
      },
      {
        path: 'admin/trends',
        element: (
          <ProtectedRoute requireAuth requireAdmin>
            <div>Admin Trends (Coming Soon)</div>
          </ProtectedRoute>
        ),
      },
      {
        path: 'admin/reports',
        element: (
          <ProtectedRoute requireAuth requireAdmin>
            <div>Admin Reports (Coming Soon)</div>
          </ProtectedRoute>
        ),
      },
      
      // 404 Not Found
      {
        path: '*',
        element: (
          <div className="container-custom py-20 text-center">
            <div className="text-6xl mb-4">🔍</div>
            <h1 className="text-4xl font-bold text-gray-900 mb-2">404</h1>
            <p className="text-gray-500 mb-6">Halaman yang Anda cari tidak ditemukan</p>
            <a href="/" className="btn-primary inline-block">
              Kembali ke Beranda
            </a>
          </div>
        ),
      },
    ],
  },
]);

// Export route paths as constants for easy reference
export const ROUTES = {
  HOME: '/',
  PRODUCTS: '/products',
  PRODUCT_DETAIL: '/product/:slug',
  CART: '/cart',
  CHECKOUT: '/checkout',
  ORDER_SUCCESS: '/order-success',
  LOGIN: '/login',
  REGISTER: '/register',
  PROFILE: '/profile',
  ORDERS: '/orders',
  WISHLIST: '/wishlist',
  ADMIN: '/admin',
  ADMIN_PRODUCTS: '/admin/products',
  ADMIN_ORDERS: '/admin/orders',
  ADMIN_INVENTORY: '/admin/inventory',
  ADMIN_TRENDS: '/admin/trends',
  ADMIN_REPORTS: '/admin/reports',
};

// Helper function to generate dynamic routes
export const getProductDetailPath = (slug) => `/product/${slug}`;
export const getOrderDetailPath = (orderNumber) => `/orders/${orderNumber}`;

export default router;