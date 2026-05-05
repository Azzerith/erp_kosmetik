import React, { useState, useEffect, useCallback } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import { motion, AnimatePresence } from 'framer-motion';
import { 
  FaStar, FaStarHalfAlt, FaRegStar, FaShoppingCart, FaHeart, FaRegHeart,
  FaShare, FaTiktok, FaInstagram, FaWhatsapp, FaChevronLeft, FaChevronRight,
  FaCheck, FaTruck, FaShieldAlt, FaUndo, FaQuestionCircle, FaMinus, FaPlus
} from 'react-icons/fa';
import { products } from '../../services/dummyData';
import { formatPrice, formatDiscount } from '../../utils/formatPrice';
import { useCart } from '../../contexts/CartContext';
import Button, { IconButton } from '../common/Button';
import LoadingSpinner, { ProductDetailSkeleton } from '../common/LoadingSpinner';
import ProductCard from '../common/ProductCard';
import TrendBadge from '../common/TrendBadge';

// Tab Component
const TabButton = ({ active, onClick, children }) => (
  <button
    onClick={onClick}
    className={`py-3 px-1 font-semibold transition-all duration-300 border-b-2 ${
      active
        ? 'text-primary border-primary'
        : 'text-gray-500 border-transparent hover:text-gray-700 hover:border-gray-300'
    }`}
  >
    {children}
  </button>
);

// Rating Stars Component
const RatingStars = ({ rating, size = 'md', showNumber = true }) => {
  const fullStars = Math.floor(rating);
  const hasHalfStar = rating % 1 >= 0.5;
  const emptyStars = 5 - fullStars - (hasHalfStar ? 1 : 0);
  
  const sizeClasses = {
    sm: 'text-sm',
    md: 'text-base',
    lg: 'text-xl',
    xl: 'text-2xl',
  };
  
  return (
    <div className={`flex items-center gap-1 ${sizeClasses[size]}`}>
      {[...Array(fullStars)].map((_, i) => (
        <FaStar key={`full-${i}`} className="text-yellow-400" />
      ))}
      {hasHalfStar && <FaStarHalfAlt className="text-yellow-400" />}
      {[...Array(emptyStars)].map((_, i) => (
        <FaRegStar key={`empty-${i}`} className="text-yellow-400" />
      ))}
      {showNumber && <span className="ml-2 text-gray-600 text-sm">({rating})</span>}
    </div>
  );
};

// Quantity Selector Component
const QuantitySelector = ({ quantity, onIncrease, onDecrease, maxStock, disabled }) => {
  return (
    <div className="flex items-center gap-3">
      <div className="flex items-center border rounded-lg overflow-hidden">
        <button
          onClick={onDecrease}
          disabled={quantity <= 1 || disabled}
          className="p-2 px-4 hover:bg-gray-100 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <FaMinus size={12} />
        </button>
        <span className="w-12 text-center font-medium">{quantity}</span>
        <button
          onClick={onIncrease}
          disabled={quantity >= maxStock || disabled}
          className="p-2 px-4 hover:bg-gray-100 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <FaPlus size={12} />
        </button>
      </div>
      <span className="text-sm text-gray-500">
        Stok: {maxStock}
      </span>
    </div>
  );
};

// Variant Selector Component
const VariantSelector = ({ variants, selectedVariant, onSelect }) => {
  if (!variants || variants.length === 0) return null;
  
  return (
    <div className="space-y-2">
      <label className="font-medium text-gray-700">Varian:</label>
      <div className="flex flex-wrap gap-2">
        {variants.map((variant) => (
          <button
            key={variant.id}
            onClick={() => onSelect(variant)}
            className={`px-4 py-2 rounded-lg border transition-all ${
              selectedVariant?.id === variant.id
                ? 'border-primary bg-primary/5 text-primary'
                : 'border-gray-300 hover:border-primary hover:bg-primary/5'
            }`}
          >
            {variant.name}
            {variant.priceModifier > 0 && (
              <span className="ml-1 text-xs text-primary">+{formatPrice(variant.priceModifier)}</span>
            )}
            {variant.priceModifier < 0 && (
              <span className="ml-1 text-xs text-green-600">{formatPrice(variant.priceModifier)}</span>
            )}
          </button>
        ))}
      </div>
    </div>
  );
};

// Image Gallery Component
const ImageGallery = ({ images, productName, trendBadge, trendScore }) => {
  const [selectedImage, setSelectedImage] = useState(0);
  const [isZoomed, setIsZoomed] = useState(false);
  
  if (!images || images.length === 0) {
    images = ['https://picsum.photos/600/600'];
  }
  
  return (
    <div className="space-y-4">
      {/* Main Image */}
      <div 
        className="relative bg-gray-100 rounded-2xl overflow-hidden aspect-square cursor-zoom-in"
        onMouseEnter={() => setIsZoomed(true)}
        onMouseLeave={() => setIsZoomed(false)}
      >
        <img
          src={images[selectedImage]}
          alt={productName}
          className={`w-full h-full object-cover transition-transform duration-300 ${
            isZoomed ? 'scale-150' : 'scale-100'
          }`}
        />
        <TrendBadge type={trendBadge} score={trendScore} className="top-4 left-4" />
      </div>
      
      {/* Thumbnails */}
      {images.length > 1 && (
        <div className="flex gap-3 overflow-x-auto pb-2">
          {images.map((img, idx) => (
            <button
              key={idx}
              onClick={() => setSelectedImage(idx)}
              className={`w-20 h-20 rounded-lg overflow-hidden border-2 flex-shrink-0 transition-all ${
                selectedImage === idx ? 'border-primary shadow-md' : 'border-transparent opacity-70 hover:opacity-100'
              }`}
            >
              <img src={img} alt={`${productName} - ${idx + 1}`} className="w-full h-full object-cover" />
            </button>
          ))}
        </div>
      )}
    </div>
  );
};

// Review Item Component
const ReviewItem = ({ review }) => {
  return (
    <div className="border-b pb-6 last:border-0">
      <div className="flex items-start gap-4">
        <img
          src={review.avatar || `https://randomuser.me/api/portraits/${review.gender || 'women'}/${review.id}.jpg`}
          alt={review.name}
          className="w-10 h-10 rounded-full object-cover"
        />
        <div className="flex-1">
          <div className="flex items-center justify-between flex-wrap gap-2 mb-1">
            <p className="font-semibold text-gray-800">{review.name}</p>
            <span className="text-sm text-gray-400">{review.date}</span>
          </div>
          <RatingStars rating={review.rating} size="sm" showNumber={false} />
          <p className="text-gray-700 mt-2 leading-relaxed">{review.comment}</p>
          {review.images && review.images.length > 0 && (
            <div className="flex gap-2 mt-3">
              {review.images.map((img, idx) => (
                <img key={idx} src={img} alt="Review" className="w-16 h-16 rounded-lg object-cover" />
              ))}
            </div>
          )}
          {review.isVerifiedPurchase && (
            <div className="flex items-center gap-1 mt-2 text-xs text-green-600">
              <FaCheck size={10} />
              <span>Pembelian Terverifikasi</span>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

// Main ProductDetail Component
const ProductDetail = () => {
  const { slug } = useParams();
  const navigate = useNavigate();
  const { addToCart } = useCart();
  
  const [product, setProduct] = useState(null);
  const [loading, setLoading] = useState(true);
  const [quantity, setQuantity] = useState(1);
  const [selectedVariant, setSelectedVariant] = useState(null);
  const [activeTab, setActiveTab] = useState('description');
  const [isWishlisted, setIsWishlisted] = useState(false);
  const [isAddingToCart, setIsAddingToCart] = useState(false);
  
  // Fetch product data
  useEffect(() => {
    const fetchProduct = async () => {
      setLoading(true);
      try {
        // Simulate API call
        await new Promise(resolve => setTimeout(resolve, 500));
        const found = products.find(p => p.slug === slug);
        
        if (!found) {
          navigate('/products');
          return;
        }
        
        setProduct(found);
        if (found.variants && found.variants.length > 0) {
          setSelectedVariant(found.variants[0]);
        }
      } catch (error) {
        console.error('Error fetching product:', error);
      } finally {
        setLoading(false);
      }
    };
    
    if (slug) {
      fetchProduct();
    }
  }, [slug, navigate]);
  
  // Get related products
  const getRelatedProducts = useCallback(() => {
    if (!product) return [];
    return products
      .filter(p => p.categoryId === product.categoryId && p.id !== product.id)
      .slice(0, 4);
  }, [product]);
  
  // Handlers
  const handleQuantityIncrease = () => {
    const maxStock = selectedVariant?.stock || product?.stock || 0;
    if (quantity < maxStock) {
      setQuantity(prev => prev + 1);
    }
  };
  
  const handleQuantityDecrease = () => {
    if (quantity > 1) {
      setQuantity(prev => prev - 1);
    }
  };
  
  const handleAddToCart = async () => {
    setIsAddingToCart(true);
    try {
      await new Promise(resolve => setTimeout(resolve, 500));
      addToCart(product, selectedVariant, quantity);
    } catch (error) {
      console.error('Error adding to cart:', error);
    } finally {
      setIsAddingToCart(false);
    }
  };
  
  const handleBuyNow = () => {
    addToCart(product, selectedVariant, quantity);
    navigate('/checkout');
  };
  
  const handleShare = () => {
    if (navigator.share) {
      navigator.share({
        title: product?.name,
        text: product?.shortDescription,
        url: window.location.href,
      });
    } else {
      navigator.clipboard.writeText(window.location.href);
      alert('Link produk disalin!');
    }
  };
  
  if (loading) {
    return <ProductDetailSkeleton />;
  }
  
  if (!product) {
    return (
      <div className="container-custom py-20 text-center">
        <div className="text-6xl mb-4">🔍</div>
        <h2 className="text-2xl font-semibold mb-2">Produk Tidak Ditemukan</h2>
        <p className="text-gray-500 mb-6">Produk yang Anda cari mungkin telah dihapus atau tidak tersedia.</p>
        <Link to="/products" className="btn-primary inline-block">
          Lihat Semua Produk
        </Link>
      </div>
    );
  }
  
  const displayPrice = selectedVariant 
    ? (product.salePrice || product.basePrice) + (selectedVariant.priceModifier || 0)
    : (product.salePrice || product.basePrice);
  const originalPrice = product.basePrice;
  const discount = formatDiscount(originalPrice, displayPrice);
  const maxStock = selectedVariant?.stock || product.stock;
  const relatedProducts = getRelatedProducts();
  
  return (
    <div className="bg-white">
      <div className="container-custom py-6 md:py-8">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-gray-500 mb-6 overflow-x-auto">
          <Link to="/" className="hover:text-primary">Beranda</Link>
          <span>/</span>
          <Link to="/products" className="hover:text-primary">Produk</Link>
          <span>/</span>
          <Link to={`/products?category=${product.category}`} className="hover:text-primary">
            {product.category}
          </Link>
          <span>/</span>
          <span className="text-gray-700 truncate">{product.name}</span>
        </div>
        
        {/* Product Main Section */}
        <div className="grid lg:grid-cols-2 gap-8 lg:gap-12 mb-12">
          {/* Left - Image Gallery */}
          <ImageGallery 
            images={product.images} 
            productName={product.name}
            trendBadge={product.trendBadge}
            trendScore={product.trendScore}
          />
          
          {/* Right - Product Info */}
          <div>
            {/* Brand & Category */}
            <div className="flex items-center gap-2 text-sm text-gray-500 mb-2">
              <span className="font-medium text-gray-700">{product.brand}</span>
              <span>•</span>
              <Link to={`/products?category=${product.category}`} className="hover:text-primary">
                {product.category}
              </Link>
            </div>
            
            {/* Title */}
            <h1 className="text-2xl md:text-3xl font-bold text-gray-900 mb-3">
              {product.name}
            </h1>
            
            {/* Rating */}
            <div className="flex items-center gap-4 mb-4">
              <RatingStars rating={product.rating} />
              <Link to="#reviews" className="text-sm text-gray-500 hover:text-primary">
                ({product.totalReviews} ulasan)
              </Link>
              <div className="text-sm text-green-600">
                Terjual {product.totalSold.toLocaleString()}+
              </div>
            </div>
            
            {/* Price */}
            <div className="mb-6">
              <div className="flex items-center gap-3 flex-wrap">
                <span className="text-3xl md:text-4xl font-bold text-primary">
                  {formatPrice(displayPrice)}
                </span>
                {discount > 0 && (
                  <>
                    <span className="text-lg text-gray-400 line-through">
                      {formatPrice(originalPrice)}
                    </span>
                    <span className="bg-red-100 text-red-600 px-2 py-1 rounded-lg text-sm font-semibold">
                      Hemat {discount}%
                    </span>
                  </>
                )}
              </div>
              {product.salePrice && (
                <p className="text-sm text-green-600 mt-1">Harga spesial! Hemat hingga {discount}%</p>
              )}
            </div>
            
            {/* Trend Info */}
            {product.trendScore > 70 && (
              <div className="bg-gradient-to-r from-primary/10 to-secondary/10 rounded-xl p-4 mb-6">
                <div className="flex items-start gap-3">
                  <div className="text-2xl">📈</div>
                  <div>
                    <p className="font-semibold text-gray-800">Produk Trending!</p>
                    <p className="text-sm text-gray-600">
                      Trend Score {product.trendScore} - Pencarian naik 40% dalam 24 jam terakhir
                    </p>
                    <div className="mt-2 h-2 bg-gray-200 rounded-full overflow-hidden">
                      <div 
                        className="h-full bg-gradient-to-r from-primary to-secondary rounded-full"
                        style={{ width: `${product.trendScore}%` }}
                      />
                    </div>
                  </div>
                </div>
              </div>
            )}
            
            {/* Short Description */}
            <p className="text-gray-600 mb-6 leading-relaxed">
              {product.shortDescription || product.description?.substring(0, 200)}
            </p>
            
            {/* Variants */}
            {product.variants && product.variants.length > 0 && (
              <div className="mb-6">
                <VariantSelector 
                  variants={product.variants}
                  selectedVariant={selectedVariant}
                  onSelect={setSelectedVariant}
                />
              </div>
            )}
            
            {/* Quantity */}
            <div className="mb-6">
              <label className="font-medium text-gray-700 block mb-2">Jumlah:</label>
              <QuantitySelector 
                quantity={quantity}
                onIncrease={handleQuantityIncrease}
                onDecrease={handleQuantityDecrease}
                maxStock={maxStock}
                disabled={isAddingToCart}
              />
            </div>
            
            {/* Action Buttons */}
            <div className="flex flex-col sm:flex-row gap-4 mb-8">
              <Button
                variant="primary"
                size="lg"
                fullWidth
                isLoading={isAddingToCart}
                onClick={handleAddToCart}
                icon={<FaShoppingCart />}
              >
                Tambah ke Keranjang
              </Button>
              <Button
                variant="outline"
                size="lg"
                fullWidth
                onClick={handleBuyNow}
              >
                Beli Sekarang
              </Button>
              <IconButton
                icon={isWishlisted ? <FaHeart className="text-red-500" /> : <FaRegHeart />}
                variant={isWishlisted ? 'danger' : 'ghost'}
                size="lg"
                onClick={() => setIsWishlisted(!isWishlisted)}
                label="Wishlist"
                className="flex-shrink-0"
              />
              <IconButton
                icon={<FaShare />}
                variant="ghost"
                size="lg"
                onClick={handleShare}
                label="Share"
                className="flex-shrink-0"
              />
            </div>
            
            {/* Shipping Info */}
            <div className="grid grid-cols-2 gap-3 pt-4 border-t">
              <div className="flex items-center gap-2 text-sm text-gray-600">
                <FaTruck className="text-primary" />
                <span>Gratis Ongkir Min. Rp150K</span>
              </div>
              <div className="flex items-center gap-2 text-sm text-gray-600">
                <FaShieldAlt className="text-primary" />
                <span>Garansi 100% Original</span>
              </div>
              <div className="flex items-center gap-2 text-sm text-gray-600">
                <FaUndo className="text-primary" />
                <span>Pengembalian 14 Hari</span>
              </div>
              <div className="flex items-center gap-2 text-sm text-gray-600">
                <FaQuestionCircle className="text-primary" />
                <span>Customer Service 24/7</span>
              </div>
            </div>
          </div>
        </div>
        
        {/* Tabs Section */}
        <div className="border-t border-b mb-12">
          <div className="flex gap-6 overflow-x-auto">
            <TabButton active={activeTab === 'description'} onClick={() => setActiveTab('description')}>
              Deskripsi
            </TabButton>
            <TabButton active={activeTab === 'ingredients'} onClick={() => setActiveTab('ingredients')}>
              Komposisi & Cara Pakai
            </TabButton>
            <TabButton active={activeTab === 'reviews'} onClick={() => setActiveTab('reviews')}>
              Ulasan ({product.totalReviews})
            </TabButton>
          </div>
        </div>
        
        {/* Tab Content */}
        <div className="mb-12">
          <AnimatePresence mode="wait">
            {activeTab === 'description' && (
              <motion.div
                key="description"
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -20 }}
                className="prose max-w-none"
              >
                <p className="text-gray-700 leading-relaxed whitespace-pre-line">
                  {product.description}
                </p>
                
                {/* Certifications */}
                {(product.isBpomCertified || product.isHalalCertified) && (
                  <div className="mt-6 grid sm:grid-cols-2 gap-4 bg-gray-50 rounded-xl p-4">
                    {product.isBpomCertified && (
                      <div className="flex items-center gap-3">
                        <div className="w-10 h-10 bg-green-100 rounded-full flex items-center justify-center text-xl">
                          ✅
                        </div>
                        <div>
                          <p className="font-semibold">BPOM Certified</p>
                          <p className="text-sm text-gray-500">Telah terdaftar di BPOM RI</p>
                        </div>
                      </div>
                    )}
                    {product.isHalalCertified && (
                      <div className="flex items-center gap-3">
                        <div className="w-10 h-10 bg-green-100 rounded-full flex items-center justify-center text-xl">
                          🕌
                        </div>
                        <div>
                          <p className="font-semibold">Halal Certified</p>
                          <p className="text-sm text-gray-500">Bersertifikat Halal MUI</p>
                        </div>
                      </div>
                    )}
                    {product.isHerbal && (
                      <div className="flex items-center gap-3">
                        <div className="w-10 h-10 bg-green-100 rounded-full flex items-center justify-center text-xl">
                          🌿
                        </div>
                        <div>
                          <p className="font-semibold">Herbal</p>
                          <p className="text-sm text-gray-500">Terbuat dari bahan alami</p>
                        </div>
                      </div>
                    )}
                  </div>
                )}
              </motion.div>
            )}
            
            {activeTab === 'ingredients' && (
              <motion.div
                key="ingredients"
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -20 }}
                className="space-y-6"
              >
                <div>
                  <h3 className="font-semibold text-lg mb-3">Komposisi</h3>
                  <ul className="list-disc list-inside space-y-2 text-gray-700">
                    <li>Vitamin C (10%) - Mencerahkan kulit</li>
                    <li>Hyaluronic Acid - Melembabkan kulit</li>
                    <li>Niacinamide - Mengecilkan pori-pori</li>
                    <li>Licorice Extract - Menyamarkan noda hitam</li>
                    <li>Aloe Vera - Menenangkan kulit</li>
                  </ul>
                </div>
                
                <div>
                  <h3 className="font-semibold text-lg mb-3">Cara Penggunaan</h3>
                  <ol className="list-decimal list-inside space-y-2 text-gray-700">
                    <li>Bersihkan wajah terlebih dahulu</li>
                    <li>Aplikasikan toner secukupnya</li>
                    <li>Ambil 2-3 tetes serum, tepuk-tepuk lembut ke wajah</li>
                    <li>Lanjutkan dengan pelembab</li>
                    <li>Gunakan sunscreen di pagi hari</li>
                  </ol>
                </div>
              </motion.div>
            )}
            
            {activeTab === 'reviews' && (
              <motion.div
                key="reviews"
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -20 }}
              >
                {/* Rating Summary */}
                <div className="bg-gray-50 rounded-xl p-6 mb-6">
                  <div className="flex items-center gap-8 flex-wrap">
                    <div className="text-center">
                      <div className="text-4xl font-bold text-primary">{product.rating}</div>
                      <RatingStars rating={product.rating} size="sm" />
                      <p className="text-sm text-gray-500 mt-1">{product.totalReviews} ulasan</p>
                    </div>
                    <div className="flex-1 space-y-2">
                      {[5, 4, 3, 2, 1].map(star => {
                        const percentage = product.totalReviews > 0 
                          ? Math.round((product.totalReviews / 5) * 100) 
                          : 0;
                        return (
                          <div key={star} className="flex items-center gap-3">
                            <span className="text-sm w-8">{star} ★</span>
                            <div className="flex-1 h-2 bg-gray-200 rounded-full overflow-hidden">
                              <div 
                                className="h-full bg-yellow-400 rounded-full"
                                style={{ width: `${percentage}%` }}
                              />
                            </div>
                            <span className="text-sm text-gray-500 w-12">{percentage}%</span>
                          </div>
                        );
                      })}
                    </div>
                  </div>
                </div>
                
                {/* Review List */}
                <div className="space-y-6">
                  {[1, 2, 3].map((review) => (
                    <ReviewItem
                      key={review}
                      review={{
                        id: review,
                        name: ['Sarah Wijaya', 'Budi Santoso', 'Dewi Putri'][review - 1],
                        rating: 5,
                        comment: [
                          'Produknya bagus banget! Wajah jadi glowing setelah pemakaian rutin. Pengiriman cepat dan packing aman. Rekomendasi!',
                          'Kualitas produk sangat baik, sesuai dengan deskripsi. Harganya worth it untuk kualitas yang didapat.',
                          'Langsung jatuh cinta dengan produk ini! Efeknya terasa sejak pemakaian pertama. Will repurchase!'
                        ][review - 1],
                        date: '2 hari yang lalu',
                        isVerifiedPurchase: true,
                        gender: review === 2 ? 'men' : 'women',
                      }}
                    />
                  ))}
                </div>
                
                {/* Load More Reviews */}
                {product.totalReviews > 3 && (
                  <div className="text-center mt-6">
                    <Button variant="outline" onClick={() => {}}>
                      Lihat Semua Ulasan ({product.totalReviews})
                    </Button>
                  </div>
                )}
              </motion.div>
            )}
          </AnimatePresence>
        </div>
        
        {/* Related Products */}
        {relatedProducts.length > 0 && (
          <div>
            <h2 className="text-2xl font-bold mb-6">Produk Terkait</h2>
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
              {relatedProducts.map((related) => (
                <ProductCard key={related.id} product={related} />
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default ProductDetail;