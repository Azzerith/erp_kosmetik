# Product Requirements Document (PRD)
## ERP Web — Penjualan Kosmetik & Obat Tradisional Berbasis Trend

**Versi:** 1.0.0  
**Tanggal:** Mei 2026  
**Status:** Draft  
**Penulis:** Product Team  

---

## Daftar Isi

1. [Executive Summary](#1-executive-summary)
2. [Latar Belakang & Tujuan](#2-latar-belakang--tujuan)
3. [Target Pengguna](#3-target-pengguna)
4. [Fitur Utama (Feature Overview)](#4-fitur-utama-feature-overview)
5. [Arsitektur Sistem](#5-arsitektur-sistem)
6. [Tech Stack](#6-tech-stack)
7. [Integrasi API Eksternal](#7-integrasi-api-eksternal)
8. [Modul & Spesifikasi Fungsional](#8-modul--spesifikasi-fungsional)
9. [Desain UI/UX & First Impression](#9-desain-uiux--first-impression)
10. [Sistem Autentikasi & OAuth](#10-sistem-autentikasi--oauth)
11. [Sistem Transaksi & Pembayaran](#11-sistem-transaksi--pembayaran)
12. [Database Schema (ERD Overview)](#12-database-schema-erd-overview)
13. [API Endpoints](#13-api-endpoints)
14. [Non-Functional Requirements](#14-non-functional-requirements)
15. [Keamanan](#15-keamanan)
16. [Roadmap & Milestone](#16-roadmap--milestone)
17. [Risiko & Mitigasi](#17-risiko--mitigasi)
18. [Glosarium](#18-glosarium)

---

## 1. Executive Summary

Platform ERP ini dirancang khusus untuk bisnis kosmetik dan obat tradisional (UMKM hingga skala menengah) yang ingin meningkatkan penjualan melalui pendekatan **data-driven trend**. Sistem ini mengintegrasikan data tren dari platform sosial media dan marketplace, sehingga produk yang ditampilkan kepada pembeli selalu relevan dengan apa yang sedang viral dan diminati pasar.

Platform ini menggabungkan:
- **Storefront** dengan first impression visual yang kuat
- **ERP internal** (inventori, produk, laporan penjualan)
- **Engine rekomendasi berbasis trend** (integrasi API eksternal)
- **Transaksi & pembayaran online** yang mulus
- **Autentikasi OAuth** untuk kemudahan login

---

## 2. Latar Belakang & Tujuan

### 2.1 Latar Belakang

Industri kosmetik dan obat tradisional di Indonesia mengalami pertumbuhan yang signifikan, terutama karena tingginya pengaruh media sosial (TikTok, Instagram, YouTube) dalam membentuk keputusan pembelian konsumen. Produk yang viral di media sosial dapat terjual habis dalam hitungan jam. Bisnis yang tidak dapat membaca tren ini akan tertinggal.

Namun, sebagian besar UMKM di segmen ini masih mengelola operasional secara manual tanpa sistem yang terintegrasi, sehingga:
- Tidak bisa memprediksi produk yang akan laris
- Stok sering tidak tersinkronisasi dengan permintaan
- Tidak memiliki data penjualan yang terstruktur
- Kehilangan potensi pembeli karena tampilan toko online yang tidak menarik

### 2.2 Tujuan Produk

| # | Tujuan | Metrik Sukses |
|---|--------|---------------|
| 1 | Meningkatkan konversi penjualan melalui rekomendasi produk berbasis trend | Conversion rate ≥ 3.5% |
| 2 | Memberikan first impression visual yang profesional kepada pembeli | Bounce rate < 40% |
| 3 | Mempermudah manajemen operasional penjual | Waktu pengerjaan order berkurang 50% |
| 4 | Menyediakan insight tren produk secara real-time | Dashboard trend diakses ≥ 5x/hari |
| 5 | Memudahkan transaksi dan pembayaran online | Payment success rate ≥ 95% |

---

## 3. Target Pengguna

### 3.1 User Personas

**Persona A — Pemilik Toko (Admin/Seller)**
- Usia: 25–45 tahun
- Profil: Pemilik bisnis kosmetik UMKM atau distributor obat tradisional
- Pain point: Tidak tahu produk apa yang sedang trending, stok sering salah perkiraan
- Kebutuhan: Dashboard sederhana, laporan penjualan, manajemen produk & stok

**Persona B — Pembeli (Customer)**
- Usia: 18–40 tahun
- Profil: Konsumen yang aktif di media sosial, mencari produk kecantikan/kesehatan
- Pain point: Sulit menemukan produk yang sedang viral di satu tempat
- Kebutuhan: Pengalaman belanja yang cepat, visual produk menarik, checkout mudah

**Persona C — Staff/Kasir**
- Profil: Karyawan toko yang mengelola order dan pengiriman
- Kebutuhan: Akses terbatas, manajemen order, konfirmasi pembayaran

### 3.2 Role & Permission

| Role | Akses |
|------|-------|
| Super Admin | Full access semua modul |
| Admin/Owner | Semua modul kecuali pengaturan sistem |
| Staff | Order management, konfirmasi pembayaran |
| Customer | Storefront, cart, checkout, order history |
| Guest | Storefront read-only |

---

## 4. Fitur Utama (Feature Overview)

### 4.1 Untuk Pembeli (Storefront)
- **Hero Section dengan Produk Trending** — menampilkan produk viral secara real-time
- **Trend Badge** — label "🔥 Trending", "✨ Viral TikTok", "⭐ Best Seller" pada produk
- **Pencarian & Filter Cerdas** — berdasarkan kategori, harga, bahan, manfaat
- **Detail Produk** — galeri foto, deskripsi, komposisi, review pembeli
- **Cart & Wishlist**
- **Checkout dengan Multiple Payment Method**
- **Order Tracking**
- **Login via OAuth (Google, Facebook)**

### 4.2 Untuk Seller/Admin (Dashboard ERP)
- **Dashboard Analitik** — penjualan harian/mingguan/bulanan, produk terlaris
- **Trend Monitor** — grafik tren produk dari API eksternal
- **Manajemen Produk** — CRUD produk, kategori, varian, harga
- **Manajemen Stok/Inventori** — stok masuk, keluar, low stock alert
- **Manajemen Order** — status order, konfirmasi, pengiriman
- **Laporan Keuangan** — revenue, profit margin, grafik tren penjualan
- **Manajemen Promo** — diskon, voucher, flash sale
- **Manajemen Pelanggan** — database pembeli, history transaksi

---

## 5. Arsitektur Sistem

```
┌─────────────────────────────────────────────────────────┐
│                    CLIENT LAYER                          │
│  ┌──────────────────┐    ┌──────────────────────────┐   │
│  │  Storefront      │    │  Admin Dashboard ERP     │   │
│  │  (React + Vite)  │    │  (React + Vite)          │   │
│  └────────┬─────────┘    └────────────┬─────────────┘   │
└───────────┼─────────────────────────── ┼────────────────┘
            │  REST API (HTTPS/JSON)      │
┌───────────▼─────────────────────────── ▼────────────────┐
│                  BACKEND LAYER                           │
│  ┌──────────────────────────────────────────────────┐   │
│  │           Golang — Gin Framework                 │   │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────────────┐  │   │
│  │  │  Auth    │ │ Product  │ │  Order/Payment   │  │   │
│  │  │  Module  │ │ Module   │ │  Module          │  │   │
│  │  └──────────┘ └──────────┘ └──────────────────┘  │   │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────────────┐  │   │
│  │  │ Inventory│ │  Trend   │ │  Report/         │  │   │
│  │  │ Module   │ │  Engine  │ │  Analytics       │  │   │
│  │  └──────────┘ └──────────┘ └──────────────────┘  │   │
│  └──────────────────────────────────────────────────┘   │
└───────────┬─────────────────────────── ┬────────────────┘
            │                             │
┌───────────▼──────────┐   ┌─────────────▼───────────────┐
│    DATABASE LAYER    │   │     EXTERNAL SERVICES        │
│  ┌────────────────┐  │   │  ┌────────────────────────┐  │
│  │  MySQL (GORM)  │  │   │  │  Google Trends API     │  │
│  └────────────────┘  │   │  │  TikTok Research API   │  │
│  ┌────────────────┐  │   │  │  Midtrans / Xendit     │  │
│  │  Redis Cache   │  │   │  │  OAuth (Google/FB)     │  │
│  └────────────────┘  │   │  │  RajaOngkir (ongkos    │  │
└──────────────────────┘   │  │  kirim)                │  │
                           │  └────────────────────────┘  │
                           └─────────────────────────────┘
```

### 5.1 Prinsip Arsitektur
- **RESTful API** — komunikasi stateless antara frontend dan backend
- **Separation of Concern** — frontend dan backend terpisah sepenuhnya
- **JWT Authentication** — token-based auth untuk setiap request
- **Redis Caching** — cache data tren dan produk populer untuk performa optimal
- **Layered Architecture di Backend** — Handler → Service → Repository → Database

---

## 6. Tech Stack

### 6.1 Frontend
| Komponen | Teknologi | Versi |
|----------|-----------|-------|
| Framework | React + Vite | React 18, Vite 5 |
| Styling | Tailwind CSS | v3 |
| State Management | Zustand / React Query | Latest |
| HTTP Client | Axios | Latest |
| UI Component | Shadcn/UI + Custom | - |
| Routing | React Router v6 | - |
| Form Handling | React Hook Form + Zod | - |
| Chart | Recharts | - |
| Animation | Framer Motion | - |
| Image Optimization | Vite Image Tools | - |

### 6.2 Backend
| Komponen | Teknologi | Versi |
|----------|-----------|-------|
| Language | Golang | 1.22+ |
| Framework | Gin | v1.9+ |
| ORM | GORM | v2 |
| Database | MySQL | 8.0+ |
| Cache | Redis | 7+ |
| Auth | JWT (golang-jwt) | - |
| OAuth | golang.org/x/oauth2 | - |
| Config | Viper | - |
| Logging | Zap | - |
| Testing | Testify | - |

### 6.3 Infrastructure
| Komponen | Pilihan |
|----------|---------|
| Containerization | Docker + Docker Compose |
| Web Server | Nginx (reverse proxy) |
| File Storage | Local / AWS S3 / Cloudinary |
| CI/CD | GitHub Actions |
| Environment | .env + Viper |

---

## 7. Integrasi API Eksternal

### 7.1 Trend Data APIs

#### A. Google Trends API (via SerpAPI / PyTrends-equivalent)
- **Tujuan:** Mendapatkan data tren pencarian produk kecantikan & kesehatan di Indonesia
- **Endpoint yang digunakan:** `/search?q={keyword}&geo=ID&hl=id`
- **Data yang diambil:** Interest over time, related queries, trending searches
- **Frekuensi polling:** Setiap 6 jam (di-cache di Redis)
- **Mapping ke produk:** Keyword dari Google Trends dicocokkan dengan tag/kategori produk

#### B. TikTok Research API (Official)
- **Tujuan:** Mendapatkan data hashtag & video trending terkait kosmetik/herbal
- **Endpoint:** `/research/video/query/`, `/research/hashtag/query/`
- **Data yang diambil:** Trending hashtags (#skincare, #jamu, #herbal), video count, like count
- **Autentikasi:** OAuth 2.0 client credentials
- **Fallback:** Jika API tidak tersedia, gunakan data statis atau scraping publik

#### C. Marketplace Trend (opsional — Tokopedia/Shopee Affiliate API)
- **Tujuan:** Data best-seller dari marketplace lokal
- **Data yang diambil:** Top selling products, kategori terlaris

### 7.2 Trend Scoring Engine

Setiap produk akan mendapatkan **Trend Score** yang dihitung dari:

```
TrendScore = (GoogleTrendIndex × 0.4) 
           + (TikTokEngagementScore × 0.4) 
           + (InternalSalesVelocity × 0.2)
```

Score ini diperbarui setiap 6 jam dan digunakan untuk:
- Menentukan urutan produk di halaman utama
- Menentukan badge yang ditampilkan (Trending, Viral, Hot)
- Memberi sinyal ke pemilik toko untuk manajemen stok

### 7.3 Payment Gateway

**Pilihan utama: Midtrans**
- Mendukung: Transfer Bank, QRIS, GoPay, OVO, Dana, Kartu Kredit, Indomaret/Alfamart
- Integrasi: Snap API (popup checkout)
- Webhook: Notifikasi status pembayaran real-time
- Sandbox testing tersedia

**Alternatif: Xendit**
- Mendukung: Virtual Account, E-wallet, QRIS, Kartu Kredit
- REST API yang clean dan dokumentasi lengkap

### 7.4 Ongkos Kirim — RajaOngkir API
- Mendukung JNE, J&T, SiCepat, Pos Indonesia, dll.
- Kalkulasi ongkos kirim berdasarkan berat dan lokasi tujuan

### 7.5 OAuth Providers
- **Google OAuth 2.0** — via `accounts.google.com`
- **Facebook Login** — via `graph.facebook.com`
- **Fallback:** Email + Password (manual registration)

---

## 8. Modul & Spesifikasi Fungsional

### 8.1 Modul Autentikasi

**FR-AUTH-01:** User dapat mendaftar dengan email & password  
**FR-AUTH-02:** User dapat login dengan akun Google (OAuth 2.0)  
**FR-AUTH-03:** User dapat login dengan akun Facebook (OAuth 2.0)  
**FR-AUTH-04:** Sistem menerbitkan JWT access token (expire: 1 jam) dan refresh token (expire: 7 hari)  
**FR-AUTH-05:** Admin dapat mengelola role dan permission user  
**FR-AUTH-06:** Password direset melalui email OTP  
**FR-AUTH-07:** Session logout dari semua device  

### 8.2 Modul Produk

**FR-PROD-01:** Admin dapat menambah, mengedit, menghapus produk  
**FR-PROD-02:** Setiap produk memiliki: nama, deskripsi, kategori, brand, komposisi/bahan, harga, stok, gambar (multi-gambar), berat, SKU  
**FR-PROD-03:** Produk mendukung varian (contoh: ukuran 50ml/100ml, warna)  
**FR-PROD-04:** Produk dapat diaktifkan/nonaktifkan  
**FR-PROD-05:** Produk memiliki Trend Score yang diperbarui otomatis  
**FR-PROD-06:** Admin dapat men-tag produk dengan keyword trend secara manual  
**FR-PROD-07:** Produk dapat diberi label khusus: "BPOM Certified", "Halal", "Herbal", "VEGAN"  
**FR-PROD-08:** Pembeli dapat memberikan rating (1–5) dan ulasan teks pada produk  

### 8.3 Modul Inventori

**FR-INV-01:** Sistem mencatat setiap perubahan stok (masuk/keluar/adjustment)  
**FR-INV-02:** Low stock alert otomatis ketika stok di bawah threshold  
**FR-INV-03:** Stok terpotong otomatis saat order dikonfirmasi  
**FR-INV-04:** Stok dikembalikan otomatis jika order dibatalkan  
**FR-INV-05:** Admin dapat melakukan stok opname (manual adjustment)  
**FR-INV-06:** Riwayat mutasi stok tersimpan lengkap  

### 8.4 Modul Trend Engine

**FR-TREND-01:** Sistem melakukan polling API Google Trends setiap 6 jam  
**FR-TREND-02:** Sistem melakukan polling TikTok Research API setiap 6 jam  
**FR-TREND-03:** Trend Score dihitung dan disimpan per produk  
**FR-TREND-04:** Dashboard menampilkan grafik tren keyword per kategori  
**FR-TREND-05:** Admin mendapat notifikasi produk yang tren-nya naik signifikan (>20% dalam 24 jam)  
**FR-TREND-06:** Produk dengan TrendScore tertinggi tampil di bagian "Trending Now"  

### 8.5 Modul Order & Transaksi

**FR-ORD-01:** Pembeli dapat menambah produk ke keranjang (cart)  
**FR-ORD-02:** Pembeli dapat mengatur jumlah produk di cart  
**FR-ORD-03:** Pembeli mengisi alamat pengiriman dan memilih ekspedisi  
**FR-ORD-04:** Sistem menghitung total harga + ongkir secara real-time  
**FR-ORD-05:** Pembeli dapat menggunakan voucher/kode promo  
**FR-ORD-06:** Sistem membuat order dengan status awal: `PENDING_PAYMENT`  
**FR-ORD-07:** Setelah pembayaran dikonfirmasi, status berubah ke `PAID`  
**FR-ORD-08:** Admin mengkonfirmasi dan memproses pengiriman, status: `PROCESSING` → `SHIPPED`  
**FR-ORD-09:** Order dapat dibatalkan oleh pembeli sebelum diproses  
**FR-ORD-10:** Pembeli dapat melacak status order secara real-time  

**Status Order Flow:**
```
PENDING_PAYMENT → PAID → PROCESSING → SHIPPED → DELIVERED → COMPLETED
                       ↘ CANCELLED (sebelum PROCESSING)
                                    ↘ RETURN_REQUESTED → REFUNDED
```

### 8.6 Modul Pembayaran

**FR-PAY-01:** Integrasi Midtrans Snap untuk checkout popup  
**FR-PAY-02:** Mendukung metode: QRIS, GoPay, OVO, Dana, Transfer Bank, Kartu Kredit, Indomaret  
**FR-PAY-03:** Webhook Midtrans untuk update status pembayaran otomatis  
**FR-PAY-04:** Sistem menyimpan log setiap event pembayaran  
**FR-PAY-05:** Expired time pembayaran: 1 jam (dapat dikonfigurasi)  
**FR-PAY-06:** Invoice PDF digenerate otomatis setelah pembayaran sukses  

### 8.7 Modul Promo & Voucher

**FR-PROMO-01:** Admin dapat membuat diskon persentase atau nominal per produk  
**FR-PROMO-02:** Admin dapat membuat voucher kode dengan batas penggunaan  
**FR-PROMO-03:** Flash sale dengan timer countdown  
**FR-PROMO-04:** Program loyalitas: poin dari setiap pembelian  
**FR-PROMO-05:** Diskon otomatis untuk pembelian minimal (minimum order)  

### 8.8 Modul Laporan & Analitik

**FR-RPT-01:** Dashboard penjualan harian, mingguan, bulanan  
**FR-RPT-02:** Top 10 produk terlaris  
**FR-RPT-03:** Grafik revenue vs. target  
**FR-RPT-04:** Laporan stok (barang masuk, keluar, sisa)  
**FR-RPT-05:** Customer acquisition report (sumber registrasi)  
**FR-RPT-06:** Export laporan ke format CSV/PDF  
**FR-RPT-07:** Korelasi antara Trend Score dengan penjualan aktual  

---

## 9. Desain UI/UX & First Impression

### 9.1 Prinsip Desain

Pembeli yang pertama kali mengunjungi website harus merasakan kesan **modern, bersih, dan terpercaya** dalam 3 detik pertama. Berikut prinsip desain yang diterapkan:

- **Visual-first** — foto produk berkualitas tinggi mendominasi layout
- **Trust Signals** — badge BPOM, Halal, rating bintang, jumlah pembeli, testimoni
- **Social Proof** — integrasi konten TikTok/Instagram yang menampilkan produk
- **Speed** — target Largest Contentful Paint (LCP) < 2.5 detik
- **Mobile-first** — desain responsif dengan prioritas tampilan mobile

### 9.2 Halaman & Komponen Utama

#### Landing Page / Homepage
- **Hero Section** — Full-width banner animasi dengan produk trending, CTA "Belanja Sekarang"
- **Trending Bar** — Scrolling ticker yang menampilkan keyword sedang trending
- **Trending Products Grid** — 8 produk dengan TrendScore tertinggi, dilengkapi badge viral
- **Category Showcase** — Navigasi visual ke kategori utama (Skincare, Haircare, Herbal, Supplement)
- **Viral on TikTok Section** — Carousel produk yang viral di TikTok (thumbnail video TikTok)
- **Flash Sale Section** — Produk diskon dengan countdown timer
- **Best Seller Section** — Produk terlaris sepanjang waktu
- **Brand Trust Section** — Logo BPOM, Halal MUI, penghargaan, jumlah customer
- **Testimonial Section** — Review bintang dan foto pembeli nyata
- **Footer** — Informasi toko, kebijakan, media sosial

#### Product Listing Page
- Filter sidebar: kategori, harga, brand, label (Halal, Herbal, BPOM)
- Sort by: Trending, Terlaris, Terbaru, Harga
- Product card: gambar, nama, harga, rating, badge trend, tombol add-to-cart

#### Product Detail Page
- Image gallery (zoom on hover)
- Trend indicator: "📈 Pencarian naik 40% minggu ini"
- Varian selector (jika ada)
- Quantity selector
- Add to cart + Buy Now button
- Tab: Deskripsi, Komposisi, Review (rating distribution chart)
- Related products (trending di kategori yang sama)

#### Cart & Checkout
- Mini cart (slide-in panel)
- Checkout multi-step: Alamat → Pengiriman → Pembayaran → Konfirmasi
- Order summary sticky di sidebar
- Midtrans Snap popup untuk pembayaran

#### Customer Dashboard
- Riwayat order, status real-time
- Wishlist
- Profil & alamat tersimpan
- Poin loyalitas

### 9.3 Palet Warna & Tipografi

| Elemen | Nilai |
|--------|-------|
| Primary Color | `#C94F7C` (Dusty Rose — identik kosmetik) |
| Secondary Color | `#2D7D6F` (Deep Teal — herbal/natural) |
| Accent | `#F5A623` (Warm Amber — highlight/badge) |
| Background | `#FAFAF8` (Off-white) |
| Text Primary | `#1A1A1A` |
| Font Heading | **Playfair Display** (elegant, premium) |
| Font Body | **Inter** (clean, readable) |

---

## 10. Sistem Autentikasi & OAuth

### 10.1 Flow OAuth Google

```
1. User klik "Login dengan Google"
2. Frontend redirect ke: GET /api/v1/auth/google
3. Backend redirect ke Google OAuth consent screen
4. Google redirect ke: GET /api/v1/auth/google/callback?code=XXX
5. Backend tukar code dengan access token ke Google
6. Backend ambil profil user (email, nama, foto) dari Google
7. Backend cek apakah email sudah ada di DB:
   - Ada → gunakan akun existing
   - Belum ada → buat akun baru (auto-register)
8. Backend generate JWT access token + refresh token
9. Backend redirect frontend dengan token di URL param atau cookie
10. Frontend simpan token di memory / httpOnly cookie
```

### 10.2 JWT Structure

```json
// Access Token Payload
{
  "sub": "user_id",
  "email": "user@email.com",
  "role": "customer",
  "iat": 1714000000,
  "exp": 1714003600
}
```

### 10.3 Refresh Token Flow
- Refresh token disimpan di database (revocable)
- Setiap refresh menghasilkan token baru (rotation)
- httpOnly cookie untuk keamanan

---

## 11. Sistem Transaksi & Pembayaran

### 11.1 Checkout Flow

```
1. Customer review cart
2. Input / pilih alamat pengiriman
3. Pilih ekspedisi → sistem kalkulasi ongkir via RajaOngkir
4. Input voucher (opsional)
5. Review order summary (subtotal + ongkir - diskon = total)
6. Klik "Bayar Sekarang"
7. Backend create Order (status: PENDING_PAYMENT)
8. Backend create Midtrans transaction → dapat snap_token
9. Frontend tampilkan Midtrans Snap popup
10. Customer selesaikan pembayaran
11. Midtrans kirim webhook ke: POST /api/v1/payments/webhook
12. Backend verifikasi signature webhook
13. Backend update Order status → PAID
14. Backend kurangi stok produk
15. Backend kirim email konfirmasi ke customer
16. Customer redirect ke halaman sukses
```

### 11.2 Webhook Security
- Verifikasi `ServerKey` Midtrans sebelum memproses webhook
- Idempotency check: setiap `order_id` hanya diproses sekali
- Logging semua webhook event

### 11.3 Refund Flow
- Customer request refund via dashboard
- Admin review dan approve/reject
- Jika approved, proses refund via Midtrans Refund API
- Status order → `REFUNDED`

---

## 12. Database Schema (ERD Overview)

### Tabel Utama

```sql
-- users
users (id, email, name, phone, avatar_url, provider, provider_id, role, 
       loyalty_points, created_at, updated_at)

-- products
products (id, sku, name, slug, description, brand, category_id, base_price, 
          sale_price, weight_gram, is_active, trend_score, trend_badge,
          total_sold, created_at, updated_at)

-- product_variants
product_variants (id, product_id, variant_name, variant_value, price_modifier, 
                  stock, sku_variant)

-- product_images
product_images (id, product_id, url, is_primary, sort_order)

-- product_tags (untuk mapping ke trend keywords)
product_tags (id, product_id, tag)

-- categories
categories (id, name, slug, parent_id, icon_url, is_active)

-- inventory_logs
inventory_logs (id, product_id, variant_id, type [IN/OUT/ADJUSTMENT], 
                quantity, note, reference_id, created_at, created_by)

-- orders
orders (id, order_number, user_id, status, subtotal, shipping_cost, 
        discount_amount, total_amount, payment_method, payment_status,
        shipping_address_id, courier, tracking_number, notes,
        paid_at, shipped_at, delivered_at, created_at)

-- order_items
order_items (id, order_id, product_id, variant_id, product_name, 
             variant_name, price, quantity, subtotal)

-- payments
payments (id, order_id, midtrans_order_id, snap_token, payment_method,
          amount, status, raw_response, paid_at, created_at)

-- addresses
addresses (id, user_id, label, recipient_name, phone, province, city, 
           district, postal_code, address_detail, is_default)

-- vouchers
vouchers (id, code, type [PERCENTAGE/FIXED], value, min_order, max_discount,
          usage_limit, used_count, valid_from, valid_until, is_active)

-- reviews
reviews (id, product_id, user_id, order_id, rating, comment, 
         images, is_verified_purchase, created_at)

-- trend_data
trend_data (id, keyword, source [GOOGLE/TIKTOK], score, metadata_json,
            recorded_at)

-- product_trend_mapping
product_trend_mapping (id, product_id, trend_id, relevance_score)
```

---

## 13. API Endpoints

### 13.1 Auth Endpoints
```
POST   /api/v1/auth/register
POST   /api/v1/auth/login
POST   /api/v1/auth/refresh
POST   /api/v1/auth/logout
GET    /api/v1/auth/google
GET    /api/v1/auth/google/callback
GET    /api/v1/auth/facebook
GET    /api/v1/auth/facebook/callback
POST   /api/v1/auth/forgot-password
POST   /api/v1/auth/reset-password
```

### 13.2 Product Endpoints
```
GET    /api/v1/products                  (list, filter, sort)
GET    /api/v1/products/:slug            (detail)
GET    /api/v1/products/trending         (top trending)
GET    /api/v1/products/flash-sale       (flash sale aktif)
GET    /api/v1/categories
POST   /api/v1/admin/products            (create) [Admin]
PUT    /api/v1/admin/products/:id        (update) [Admin]
DELETE /api/v1/admin/products/:id        (delete) [Admin]
```

### 13.3 Order Endpoints
```
POST   /api/v1/orders                    (create order)
GET    /api/v1/orders                    (my orders)
GET    /api/v1/orders/:order_number      (order detail)
POST   /api/v1/orders/:id/cancel         (cancel order)
GET    /api/v1/admin/orders              (all orders) [Admin]
PUT    /api/v1/admin/orders/:id/status   (update status) [Admin]
```

### 13.4 Payment Endpoints
```
POST   /api/v1/payments/initiate         (get snap token)
POST   /api/v1/payments/webhook          (midtrans webhook)
GET    /api/v1/payments/:order_id        (payment status)
```

### 13.5 Trend Endpoints
```
GET    /api/v1/trends                    (trending keywords)
GET    /api/v1/trends/products           (produk by trend)
GET    /api/v1/admin/trends/dashboard    (trend analytics) [Admin]
```

### 13.6 Shipping Endpoints
```
GET    /api/v1/shipping/provinces
GET    /api/v1/shipping/cities?province_id=X
POST   /api/v1/shipping/calculate        (hitung ongkir)
```

### 13.7 Admin Report Endpoints
```
GET    /api/v1/admin/reports/sales       (penjualan)
GET    /api/v1/admin/reports/inventory   (stok)
GET    /api/v1/admin/reports/customers   (pelanggan)
GET    /api/v1/admin/dashboard/summary   (ringkasan KPI)
```

---

## 14. Non-Functional Requirements

### 14.1 Performa
| Metrik | Target |
|--------|--------|
| API Response Time (P95) | < 200ms |
| Page Load Time (LCP) | < 2.5 detik |
| Time to Interactive | < 3.5 detik |
| Uptime | 99.5% |
| Concurrent Users | ≥ 500 concurrent |

### 14.2 Skalabilitas
- Backend stateless, dapat di-scale horizontal
- Redis untuk caching trend data dan session
- Database connection pooling dengan GORM
- Gambar disimpan di object storage (S3/Cloudinary), bukan server lokal

### 14.3 Kompatibilitas Browser
- Chrome 90+, Firefox 88+, Safari 14+, Edge 90+
- Mobile: Android Chrome, iOS Safari

### 14.4 Aksesibilitas
- WCAG 2.1 Level AA
- Alt text pada semua gambar produk
- Keyboard navigable
- Kontras warna minimal 4.5:1

---

## 15. Keamanan

### 15.1 Autentikasi & Otorisasi
- JWT dengan signature RS256 atau HS256
- Refresh token rotation (single-use)
- Role-based access control (RBAC) di setiap endpoint
- Middleware auth wajib di semua endpoint admin

### 15.2 Data Protection
- Password di-hash dengan bcrypt (cost factor ≥ 12)
- Data sensitif (nomor kartu, data pribadi) tidak disimpan di sistem (delegasi ke Midtrans)
- HTTPS wajib untuk semua endpoint
- CORS dikonfigurasi hanya untuk origin yang diizinkan

### 15.3 API Security
- Rate limiting: 60 request/menit per IP untuk publik, 300 untuk user login
- Request validation ketat di semua endpoint (Gin binding + custom validator)
- SQL injection prevention via GORM parameterized query
- XSS prevention via input sanitization
- CSRF token untuk form sensitif

### 15.4 Payment Security
- Verifikasi Midtrans webhook signature wajib
- Order amount di-verify di backend, tidak dari frontend
- Idempotency key untuk mencegah double payment

### 15.5 Logging & Monitoring
- Semua request dicatat dengan request ID
- Error logging dengan stack trace (Zap logger)
- Alert otomatis jika error rate > 1% atau latency > 500ms

---

## 16. Roadmap & Milestone

### Phase 1 — MVP (Bulan 1–2)
- [ ] Setup project structure (React Vite + Golang Gin)
- [ ] Database schema & migrasi
- [ ] Auth module (email + OAuth Google)
- [ ] Manajemen produk & kategori (CRUD)
- [ ] Storefront: Homepage, Product List, Product Detail
- [ ] Cart & Basic Checkout
- [ ] Integrasi Midtrans (sandbox)
- [ ] Order management dasar

### Phase 2 — Trend Engine & ERP Core (Bulan 3–4)
- [ ] Integrasi Google Trends API
- [ ] Integrasi TikTok Research API
- [ ] Trend Score engine & cron job
- [ ] Trending products section di homepage
- [ ] Inventori management lengkap
- [ ] Dashboard admin: penjualan & laporan
- [ ] OAuth Facebook
- [ ] Integrasi RajaOngkir

### Phase 3 — Engagement & Optimasi (Bulan 5–6)
- [ ] Flash sale & voucher system
- [ ] Review & rating produk
- [ ] Loyalty points system
- [ ] Notifikasi email (order konfirmasi, pengiriman)
- [ ] Advanced analytics & trend correlation
- [ ] Performance optimization (caching, image CDN)
- [ ] Mobile responsiveness audit & polish
- [ ] Load testing & security audit

### Phase 4 — Scale (Bulan 7+)
- [ ] Marketplace integration (Tokopedia/Shopee sync)
- [ ] PWA (Progressive Web App)
- [ ] Push notification
- [ ] AI-powered product recommendation
- [ ] Multi-seller/multi-toko

---

## 17. Risiko & Mitigasi

| Risiko | Dampak | Probabilitas | Mitigasi |
|--------|--------|--------------|----------|
| TikTok Research API akses terbatas/berbayar | Tinggi | Sedang | Siapkan fallback scraping publik atau gunakan data internal sebagai alternatif |
| Google Trends API rate limit | Sedang | Tinggi | Implementasi caching Redis agresif (6 jam), exponential backoff |
| Midtrans downtime | Tinggi | Rendah | Siapkan fallback ke Xendit, monitoring webhook |
| Performa lambat saat produk trending | Tinggi | Sedang | Redis caching untuk halaman trending, CDN untuk aset statis |
| Data stok tidak sinkron | Tinggi | Sedang | Optimistic locking pada stok, transaksional DB |
| Fraud / pembayaran palsu | Tinggi | Sedang | Verifikasi webhook signature, server-side amount validation |

---

## 18. Glosarium

| Istilah | Definisi |
|---------|----------|
| ERP | Enterprise Resource Planning — sistem terintegrasi untuk manajemen bisnis |
| Trend Score | Skor numerik yang merepresentasikan tingkat popularitas produk saat ini |
| OAuth | Protokol autentikasi delegasi yang memungkinkan login via akun pihak ketiga |
| JWT | JSON Web Token — standar token untuk autentikasi stateless |
| Snap Token | Token sementara dari Midtrans untuk membuka payment popup |
| Webhook | HTTP callback yang dikirim otomatis oleh layanan eksternal saat event terjadi |
| SKU | Stock Keeping Unit — kode unik identifikasi produk/varian |
| LCP | Largest Contentful Paint — metrik performa loading halaman web |
| GORM | Go Object Relational Mapper — library ORM untuk Golang |
| Idempotency | Sifat operasi yang menghasilkan hasil sama meskipun dieksekusi berkali-kali |

---

*Dokumen ini bersifat living document dan akan diperbarui sesuai perkembangan proyek.*

**Prepared by:** Product Team  
**Last Updated:** Mei 2026  
**Next Review:** Juni 2026
