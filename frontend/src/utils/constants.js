// App Constants

export const APP_NAME = 'ErpCosmetics';
export const APP_DESCRIPTION = 'Platform belanja produk kecantikan dan herbal berbasis trend terkini';

// Routes
export const ROUTES = {
  HOME: '/',
  PRODUCTS: '/products',
  PRODUCT_DETAIL: '/product/:slug',
  CART: '/cart',
  CHECKOUT: '/checkout',
  ORDER_SUCCESS: '/order-success',
  LOGIN: '/login',
  REGISTER: '/register',
  FORGOT_PASSWORD: '/forgot-password',
  RESET_PASSWORD: '/reset-password',
  PROFILE: '/profile',
  ORDERS: '/orders',
  ORDER_DETAIL: '/orders/:orderNumber',
  WISHLIST: '/wishlist',
  TRENDING: '/trending',
  FLASH_SALE: '/flash-sale',
};

// Order Status
export const ORDER_STATUS = {
  PENDING_PAYMENT: 'pending_payment',
  PAID: 'paid',
  PROCESSING: 'processing',
  SHIPPED: 'shipped',
  DELIVERED: 'delivered',
  COMPLETED: 'completed',
  CANCELLED: 'cancelled',
  REFUNDED: 'refunded',
};

export const ORDER_STATUS_LABELS = {
  [ORDER_STATUS.PENDING_PAYMENT]: { label: 'Menunggu Pembayaran', color: 'warning' },
  [ORDER_STATUS.PAID]: { label: 'Sudah Dibayar', color: 'info' },
  [ORDER_STATUS.PROCESSING]: { label: 'Diproses', color: 'info' },
  [ORDER_STATUS.SHIPPED]: { label: 'Dikirim', color: 'primary' },
  [ORDER_STATUS.DELIVERED]: { label: 'Telah Sampai', color: 'success' },
  [ORDER_STATUS.COMPLETED]: { label: 'Selesai', color: 'success' },
  [ORDER_STATUS.CANCELLED]: { label: 'Dibatalkan', color: 'error' },
  [ORDER_STATUS.REFUNDED]: { label: 'Dikembalikan', color: 'error' },
};

// Payment Methods
export const PAYMENT_METHODS = [
  { id: 'qris', name: 'QRIS', icon: '📱' },
  { id: 'bank_transfer', name: 'Transfer Bank', icon: '🏦' },
  { id: 'credit_card', name: 'Kartu Kredit', icon: '💳' },
  { id: 'ewallet', name: 'E-Wallet', icon: '👛' },
  { id: 'cod', name: 'COD (Bayar di Tempat)', icon: '💰' },
];

// Shipping Couriers
export const COURIERS = [
  { id: 'jne', name: 'JNE', icon: '📦' },
  { id: 'jnt', name: 'J&T Express', icon: '📦' },
  { id: 'sicepat', name: 'SiCepat', icon: '📦' },
  { id: 'pos', name: 'Pos Indonesia', icon: '📮' },
];

// Trend Badge Types
export const TREND_BADGES = {
  viral: { label: 'VIRAL TikTok', icon: '🔥', color: 'gradient-pink' },
  trending: { label: 'TRENDING', icon: '📈', color: 'accent' },
  best_seller: { label: 'BEST SELLER', icon: '⭐', color: 'primary' },
  hot: { label: 'HOT', icon: '🔥', color: 'orange' },
};

// Product Certifications
export const CERTIFICATIONS = [
  { id: 'bpom', label: 'BPOM Certified', icon: '✅' },
  { id: 'halal', label: 'Halal Certified', icon: '🕌' },
  { id: 'vegan', label: 'Vegan', icon: '🌱' },
  { id: 'herbal', label: 'Herbal', icon: '🌿' },
];

// Local Storage Keys
export const STORAGE_KEYS = {
  TOKEN: 'token',
  REFRESH_TOKEN: 'refreshToken',
  USER: 'user',
  CART: 'cart',
  THEME: 'theme',
  LANGUAGE: 'language',
};

// Pagination
export const DEFAULT_PAGE_SIZE = 12;
export const PAGE_SIZE_OPTIONS = [12, 24, 48, 96];

// Price Format
export const CURRENCY = 'IDR';
export const CURRENCY_SYMBOL = 'Rp';
export const THOUSAND_SEPARATOR = '.';
export const DECIMAL_SEPARATOR = ',';

// Image Placeholders
export const PLACEHOLDER_IMAGE = 'https://picsum.photos/400/500';
export const PLACEHOLDER_AVATAR = 'https://picsum.photos/100/100';

// API Endpoints
export const API_ENDPOINTS = {
  AUTH: '/auth',
  PRODUCTS: '/products',
  ORDERS: '/orders',
  PAYMENTS: '/payments',
  CART: '/cart',
  TRENDS: '/trends',
  SHIPPING: '/shipping',
  VOUCHERS: '/vouchers',
  REVIEWS: '/reviews',
  USERS: '/users',
  REPORTS: '/reports',
};

// Regex Patterns
export const PATTERNS = {
  EMAIL: /^[^\s@]+@([^\s@.,]+\.)+[^\s@.,]{2,}$/,
  PHONE: /^(\+62|62|0)8[1-9][0-9]{6,10}$/,
  PASSWORD: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d]{8,}$/,
  POSTAL_CODE: /^[0-9]{5}$/,
};

// Error Messages
export const ERROR_MESSAGES = {
  REQUIRED: 'Field ini wajib diisi',
  EMAIL_INVALID: 'Format email tidak valid',
  PHONE_INVALID: 'Format nomor telepon tidak valid',
  PASSWORD_WEAK: 'Password minimal 8 karakter, mengandung huruf besar, huruf kecil, dan angka',
  PASSWORDS_NOT_MATCH: 'Password tidak cocok',
  NETWORK_ERROR: 'Terjadi kesalahan jaringan. Silakan coba lagi',
  SERVER_ERROR: 'Terjadi kesalahan pada server. Silakan coba lagi nanti',
  UNAUTHORIZED: 'Sesi Anda telah berakhir. Silakan login kembali',
  FORBIDDEN: 'Anda tidak memiliki akses ke halaman ini',
  NOT_FOUND: 'Data tidak ditemukan',
};

// Success Messages
export const SUCCESS_MESSAGES = {
  LOGIN: 'Login berhasil!',
  REGISTER: 'Pendaftaran berhasil! Silakan login',
  LOGOUT: 'Berhasil logout',
  ADD_TO_CART: 'Produk ditambahkan ke keranjang',
  REMOVE_FROM_CART: 'Produk dihapus dari keranjang',
  UPDATE_CART: 'Keranjang berhasil diperbarui',
  ORDER_CREATED: 'Pesanan berhasil dibuat',
  PAYMENT_SUCCESS: 'Pembayaran berhasil',
  REVIEW_SUBMITTED: 'Ulasan berhasil dikirim',
  PROFILE_UPDATED: 'Profil berhasil diperbarui',
  PASSWORD_CHANGED: 'Password berhasil diubah',
};

// Theme
export const THEMES = {
  LIGHT: 'light',
  DARK: 'dark',
};

// Breakpoints (Tailwind)
export const BREAKPOINTS = {
  sm: 640,
  md: 768,
  lg: 1024,
  xl: 1280,
  '2xl': 1536,
};

// Social Media Links
export const SOCIAL_LINKS = {
  instagram: 'https://instagram.com/erpcosmetics',
  tiktok: 'https://tiktok.com/@erpcosmetics',
  facebook: 'https://facebook.com/erpcosmetics',
  twitter: 'https://twitter.com/erpcosmetics',
  youtube: 'https://youtube.com/@erpcosmetics',
};

// Contact Info
export const CONTACT_INFO = {
  email: 'support@erpcosmetics.com',
  phone: '+62 812 3456 7890',
  whatsapp: '6281234567890',
  address: 'Jl. Contoh No. 123, Jakarta Selatan, Indonesia',
};

export default {
  APP_NAME,
  APP_DESCRIPTION,
  ROUTES,
  ORDER_STATUS,
  ORDER_STATUS_LABELS,
  PAYMENT_METHODS,
  COURIERS,
  TREND_BADGES,
  CERTIFICATIONS,
  STORAGE_KEYS,
  DEFAULT_PAGE_SIZE,
  PAGE_SIZE_OPTIONS,
  CURRENCY,
  CURRENCY_SYMBOL,
  PLACEHOLDER_IMAGE,
  API_ENDPOINTS,
  PATTERNS,
  ERROR_MESSAGES,
  SUCCESS_MESSAGES,
};