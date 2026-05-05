import React, { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { motion } from 'framer-motion';
import { FaStar, FaShoppingCart, FaHeart, FaShare, FaTiktok, FaInstagram, FaWhatsapp } from 'react-icons/fa';
import { products } from '../services/dummyData';
import { formatPrice } from '../utils/formatPrice';
import { useCart } from '../contexts/CartContext';
import TrendBadge from '../components/common/TrendBadge';
import ProductCard from '../components/common/ProductCard';

const ProductDetailPage = () => {
  const { slug } = useParams();
  const [product, setProduct] = useState(null);
  const [selectedVariant, setSelectedVariant] = useState(null);
  const [quantity, setQuantity] = useState(1);
  const [activeTab, setActiveTab] = useState('description');
  const [selectedImage, setSelectedImage] = useState(0);
  const { addToCart } = useCart();

  useEffect(() => {
    const found = products.find(p => p.slug === slug);
    setProduct(found);
    if (found && found.variants && found.variants.length > 0) {
      setSelectedVariant(found.variants[0]);
    }
  }, [slug]);

  if (!product) {
    return (
      <div className="container-custom py-20 text-center">
        <div className="text-4xl mb-4">🔍</div>
        <p className="text-gray-500">Produk tidak ditemukan</p>
        <Link to="/products" className="btn-primary inline-block mt-4">
          Lihat Produk Lain
        </Link>
      </div>
    );
  }

  const displayPrice = selectedVariant 
    ? product.salePrice + (selectedVariant.priceModifier || 0)
    : (product.salePrice || product.basePrice);

  const originalPrice = product.basePrice;

  const handleAddToCart = () => {
    addToCart(product, selectedVariant, quantity);
  };

  const relatedProducts = products
    .filter(p => p.categoryId === product.categoryId && p.id !== product.id)
    .slice(0, 4);

  return (
    <div className="bg-white">
      <div className="container-custom py-8">
        {/* Breadcrumb */}
        <div className="text-sm text-gray-500 mb-6">
          <Link to="/" className="hover:text-primary">Beranda</Link>
          <span className="mx-2">/</span>
          <Link to="/products" className="hover:text-primary">Produk</Link>
          <span className="mx-2">/</span>
          <span className="text-gray-700">{product.name}</span>
        </div>

        {/* Product Main */}
        <div className="grid md:grid-cols-2 gap-8 mb-12">
          {/* Image Gallery */}
          <div className="space-y-4">
            <div className="relative bg-gray-100 rounded-2xl overflow-hidden aspect-square">
              <img
                src={product.images?.[selectedImage] || 'https://picsum.photos/400/500'}
                alt={product.name}
                className="w-full h-full object-cover"
              />
              <TrendBadge type={product.trendBadge} score={product.trendScore} className="top-4 left-4" />
            </div>
            {product.images && product.images.length > 1 && (
              <div className="flex gap-3">
                {product.images.map((img, idx) => (
                  <button
                    key={idx}
                    onClick={() => setSelectedImage(idx)}
                    className={`w-20 h-20 rounded-lg overflow-hidden border-2 ${
                      selectedImage === idx ? 'border-primary' : 'border-transparent'
                    }`}
                  >
                    <img src={img} alt="" className="w-full h-full object-cover" />
                  </button>
                ))}
              </div>
            )}
          </div>

          {/* Product Info */}
          <div>
            <div className="flex items-center gap-2 mb-2">
              <span className="text-sm text-gray-500">{product.brand}</span>
              <span className="text-gray-300">|</span>
              <span className="text-sm text-gray-500">{product.category}</span>
            </div>

            <h1 className="text-3xl font-bold text-gray-900 mb-4">{product.name}</h1>

            {/* Rating */}
            <div className="flex items-center gap-4 mb-4">
              <div className="flex items-center gap-1">
                <FaStar className="text-yellow-400" />
                <span className="font-semibold">{product.rating}</span>
                <span className="text-gray-400">({product.totalReviews} ulasan)</span>
              </div>
              <div className="text-sm text-green-600">
                Terjual {product.totalSold}+
              </div>
            </div>

            {/* Price */}
            <div className="mb-6">
              {product.salePrice ? (
                <div className="flex items-center gap-3">
                  <span className="text-3xl font-bold text-primary">{formatPrice(displayPrice)}</span>
                  <span className="text-lg text-gray-400 line-through">{formatPrice(originalPrice)}</span>
                  <span className="bg-red-100 text-red-600 px-2 py-1 rounded-lg text-sm font-semibold">
                    Hemat {formatPrice(originalPrice - displayPrice)}
                  </span>
                </div>
              ) : (
                <span className="text-3xl font-bold text-primary">{formatPrice(displayPrice)}</span>
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
                  </div>
                </div>
              </div>
            )}

            {/* Variants */}
            {product.variants && product.variants.length > 0 && (
              <div className="mb-6">
                <h3 className="font-semibold mb-2">Varian:</h3>
                <div className="flex gap-2">
                  {product.variants.map(variant => (
                    <button
                      key={variant.id}
                      onClick={() => setSelectedVariant(variant)}
                      className={`px-4 py-2 rounded-lg border transition-all ${
                        selectedVariant?.id === variant.id
                          ? 'border-primary bg-primary/5 text-primary'
                          : 'border-gray-300 hover:border-primary'
                      }`}
                    >
                      {variant.name}
                    </button>
                  ))}
                </div>
              </div>
            )}

            {/* Quantity */}
            <div className="mb-6">
              <h3 className="font-semibold mb-2">Jumlah:</h3>
              <div className="flex items-center gap-3">
                <button
                  onClick={() => setQuantity(Math.max(1, quantity - 1))}
                  className="w-10 h-10 border rounded-lg hover:bg-gray-100"
                >
                  -
                </button>
                <span className="w-12 text-center text-lg">{quantity}</span>
                <button
                  onClick={() => setQuantity(quantity + 1)}
                  className="w-10 h-10 border rounded-lg hover:bg-gray-100"
                >
                  +
                </button>
                <span className="text-sm text-gray-500">
                  Stok: {selectedVariant?.stock || product.stock}
                </span>
              </div>
            </div>

            {/* Actions */}
            <div className="flex gap-4 mb-6">
              <button onClick={handleAddToCart} className="btn-primary flex-1 flex items-center justify-center gap-2">
                <FaShoppingCart />
                Tambah ke Keranjang
              </button>
              <button className="btn-outline px-6">
                <FaHeart />
              </button>
              <button className="btn-outline px-6">
                <FaShare />
              </button>
            </div>

            {/* Social Proof */}
            <div className="flex gap-6 pt-4 border-t">
              <div className="flex items-center gap-2">
                <FaTiktok className="text-gray-600" />
                <span className="text-sm">2.5M views</span>
              </div>
              <div className="flex items-center gap-2">
                <FaInstagram className="text-gray-600" />
                <span className="text-sm">10K+ posts</span>
              </div>
              <div className="flex items-center gap-2">
                <FaWhatsapp className="text-gray-600" />
                <span className="text-sm">Customer Service</span>
              </div>
            </div>
          </div>
        </div>

        {/* Tabs - Description & Reviews */}
        <div className="border-t border-b mb-12">
          <div className="flex gap-8">
            <button
              onClick={() => setActiveTab('description')}
              className={`py-3 font-semibold transition-colors ${
                activeTab === 'description'
                  ? 'text-primary border-b-2 border-primary'
                  : 'text-gray-500 hover:text-gray-700'
              }`}
            >
              Deskripsi
            </button>
            <button
              onClick={() => setActiveTab('ingredients')}
              className={`py-3 font-semibold transition-colors ${
                activeTab === 'ingredients'
                  ? 'text-primary border-b-2 border-primary'
                  : 'text-gray-500 hover:text-gray-700'
              }`}
            >
              Komposisi
            </button>
            <button
              onClick={() => setActiveTab('reviews')}
              className={`py-3 font-semibold transition-colors ${
                activeTab === 'reviews'
                  ? 'text-primary border-b-2 border-primary'
                  : 'text-gray-500 hover:text-gray-700'
              }`}
            >
              Ulasan ({product.totalReviews})
            </button>
          </div>
        </div>

        <div className="mb-12">
          {activeTab === 'description' && (
            <div className="prose max-w-none">
              <p className="text-gray-700 leading-relaxed">{product.description}</p>
              <div className="mt-6 grid md:grid-cols-2 gap-4">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 bg-green-100 rounded-full flex items-center justify-center">✅</div>
                  <div>
                    <p className="font-semibold">BPOM Certified</p>
                    <p className="text-sm text-gray-500">Telah terdaftar di BPOM</p>
                  </div>
                </div>
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 bg-green-100 rounded-full flex items-center justify-center">✅</div>
                  <div>
                    <p className="font-semibold">Halal Certified</p>
                    <p className="text-sm text-gray-500">Bersertifikat Halal MUI</p>
                  </div>
                </div>
              </div>
            </div>
          )}

          {activeTab === 'ingredients' && (
            <div>
              <p className="text-gray-700 mb-4">
                {product.name} diformulasikan dengan bahan-bahan berkualitas:
              </p>
              <ul className="list-disc list-inside space-y-2 text-gray-700">
                <li>Vitamin C (10%) - Mencerahkan kulit</li>
                <li>Hyaluronic Acid - Melembabkan kulit</li>
                <li>Niacinamide - Mengecilkan pori-pori</li>
                <li>Licorice Extract - Menyamarkan noda hitam</li>
                <li>Aloe Vera - Menenangkan kulit</li>
              </ul>
            </div>
          )}

          {activeTab === 'reviews' && (
            <div className="space-y-6">
              {/* Sample reviews */}
              {[1, 2, 3].map((review) => (
                <div key={review} className="border-b pb-6">
                  <div className="flex items-start gap-4">
                    <img
                      src={`https://randomuser.me/api/portraits/women/${review}.jpg`}
                      alt="User"
                      className="w-10 h-10 rounded-full"
                    />
                    <div className="flex-1">
                      <div className="flex items-center justify-between mb-1">
                        <p className="font-semibold">Sarah Wijaya</p>
                        <span className="text-sm text-gray-400">2 hari yang lalu</span>
                      </div>
                      <div className="flex items-center gap-1 mb-2">
                        {[...Array(5)].map((_, i) => (
                          <FaStar key={i} className={i < 5 ? 'text-yellow-400' : 'text-gray-300'} size={14} />
                        ))}
                      </div>
                      <p className="text-gray-700">Produknya bagus banget! Wajah jadi glowing setelah pemakaian rutin. Pengiriman cepat dan packing aman. Rekomendasi!</p>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Related Products */}
        {relatedProducts.length > 0 && (
          <div>
            <h2 className="text-2xl font-bold mb-6">Produk Terkait</h2>
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
              {relatedProducts.map(related => (
                <ProductCard key={related.id} product={related} />
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default ProductDetailPage;