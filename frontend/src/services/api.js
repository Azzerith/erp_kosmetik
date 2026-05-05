import axios from 'axios';
import toast from 'react-hot-toast';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

// Create axios instance
const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor untuk menambahkan token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor untuk handle error
api.interceptors.response.use(
  (response) => {
    return response.data;
  },
  async (error) => {
    const originalRequest = error.config;
    
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;
      
      try {
        const refreshToken = localStorage.getItem('refreshToken');
        const response = await axios.post(`${API_BASE_URL}/auth/refresh`, {
          refresh_token: refreshToken,
        });
        
        const { access_token, refresh_token } = response.data;
        localStorage.setItem('token', access_token);
        localStorage.setItem('refreshToken', refresh_token);
        
        originalRequest.headers.Authorization = `Bearer ${access_token}`;
        return api(originalRequest);
      } catch (refreshError) {
        localStorage.removeItem('token');
        localStorage.removeItem('refreshToken');
        localStorage.removeItem('user');
        window.location.href = '/login';
        return Promise.reject(refreshError);
      }
    }
    
    // Show error toast untuk non-auth errors
    if (error.response?.data?.message) {
      toast.error(error.response.data.message);
    } else if (error.message) {
      toast.error(error.message);
    }
    
    return Promise.reject(error);
  }
);

// ======================================================
// AUTH API
// ======================================================

export const authAPI = {
  register: (data) => api.post('/auth/register', data),
  login: (data) => api.post('/auth/login', data),
  logout: () => api.post('/auth/logout'),
  refreshToken: (refreshToken) => api.post('/auth/refresh', { refresh_token: refreshToken }),
  forgotPassword: (email) => api.post('/auth/forgot-password', { email }),
  resetPassword: (token, password) => api.post('/auth/reset-password', { token, password }),
  getMe: () => api.get('/auth/me'),
  updateProfile: (data) => api.put('/auth/profile', data),
  changePassword: (data) => api.post('/auth/change-password', data),
};

// ======================================================
// PRODUCT API
// ======================================================

export const productAPI = {
  getProducts: (params) => api.get('/products', { params }),
  getProductBySlug: (slug) => api.get(`/products/${slug}`),
  getProductById: (id) => api.get(`/products/id/${id}`),
  getTrendingProducts: (limit = 8) => api.get('/products/trending', { params: { limit } }),
  getFlashSale: () => api.get('/products/flash-sale'),
  getBestSellers: (limit = 8) => api.get('/products/best-sellers', { params: { limit } }),
  getCategories: () => api.get('/categories'),
  getBrands: () => api.get('/brands'),
  searchProducts: (query) => api.get('/products/search', { params: { q: query } }),
  
  // Admin only
  createProduct: (data) => api.post('/admin/products', data),
  updateProduct: (id, data) => api.put(`/admin/products/${id}`, data),
  deleteProduct: (id) => api.delete(`/admin/products/${id}`),
};

// ======================================================
// ORDER API
// ======================================================

export const orderAPI = {
  createOrder: (data) => api.post('/orders', data),
  getOrders: (params) => api.get('/orders', { params }),
  getOrderById: (orderNumber) => api.get(`/orders/${orderNumber}`),
  cancelOrder: (id) => api.post(`/orders/${id}/cancel`),
  getOrderStatus: (orderNumber) => api.get(`/orders/${orderNumber}/status`),
  
  // Admin only
  getAllOrders: (params) => api.get('/admin/orders', { params }),
  updateOrderStatus: (id, status) => api.put(`/admin/orders/${id}/status`, { status }),
  updateTracking: (id, trackingNumber, courier) => api.put(`/admin/orders/${id}/tracking`, { tracking_number: trackingNumber, courier }),
};

// ======================================================
// PAYMENT API
// ======================================================

export const paymentAPI = {
  initiatePayment: (orderId) => api.post('/payments/initiate', { order_id: orderId }),
  getPaymentStatus: (orderId) => api.get(`/payments/${orderId}/status`),
  getPaymentHistory: (orderId) => api.get(`/payments/${orderId}/history`),
};

// ======================================================
// CART API
// ======================================================

export const cartAPI = {
  getCart: () => api.get('/cart'),
  addToCart: (productId, variantId, quantity) => api.post('/cart/items', { product_id: productId, variant_id: variantId, quantity }),
  updateCartItem: (itemId, quantity) => api.put(`/cart/items/${itemId}`, { quantity }),
  removeCartItem: (itemId) => api.delete(`/cart/items/${itemId}`),
  clearCart: () => api.delete('/cart/clear'),
};

// ======================================================
// TREND API
// ======================================================

export const trendAPI = {
  getTrendingKeywords: () => api.get('/trends/keywords'),
  getTrendingProducts: () => api.get('/trends/products'),
  getTrendScore: (productId) => api.get(`/trends/products/${productId}/score`),
  getTrendHistory: (productId, days = 30) => api.get(`/trends/products/${productId}/history`, { params: { days } }),
};

// ======================================================
// SHIPPING API
// ======================================================

export const shippingAPI = {
  getProvinces: () => api.get('/shipping/provinces'),
  getCities: (provinceId) => api.get(`/shipping/cities/${provinceId}`),
  calculateCost: (origin, destination, weight, courier) => api.post('/shipping/calculate', {
    origin, destination, weight, courier,
  }),
};

// ======================================================
// VOUCHER API
// ======================================================

export const voucherAPI = {
  getVouchers: () => api.get('/vouchers'),
  validateVoucher: (code, orderAmount) => api.post('/vouchers/validate', { code, order_amount: orderAmount }),
  applyVoucher: (code, orderId) => api.post('/vouchers/apply', { code, order_id: orderId }),
};

// ======================================================
// REVIEW API
// ======================================================

export const reviewAPI = {
  getReviews: (productId, params) => api.get(`/reviews/products/${productId}`, { params }),
  createReview: (data) => api.post('/reviews', data),
  updateReview: (id, data) => api.put(`/reviews/${id}`, data),
  deleteReview: (id) => api.delete(`/reviews/${id}`),
  markHelpful: (id) => api.post(`/reviews/${id}/helpful`),
};

export default api;