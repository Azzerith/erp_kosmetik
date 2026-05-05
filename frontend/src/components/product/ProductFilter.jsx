import React, { useState } from 'react';
import { FaChevronDown, FaChevronUp } from 'react-icons/fa';
import { categories } from '../../services/dummyData';

const ProductFilter = ({ onFilter }) => {
  const [selectedCategories, setSelectedCategories] = useState([]);
  const [priceRange, setPriceRange] = useState({ min: 0, max: 500000 });
  const [selectedBrands, setSelectedBrands] = useState([]);
  const [certifications, setCertifications] = useState([]);
  const [expandedSections, setExpandedSections] = useState({
    categories: true,
    price: true,
    brands: true,
    certifications: true,
  });

  const brands = ['GlowLab', 'HerbalIndo', 'DewyLab', 'SunShield', 'Lush Beauty', 'HairGro', 'BodyLove'];
  
  const certificationOptions = [
    { id: 'bpom', label: 'BPOM Certified' },
    { id: 'halal', label: 'Halal Certified' },
    { id: 'vegan', label: 'Vegan' },
    { id: 'herbal', label: 'Herbal' },
  ];

  const toggleSection = (section) => {
    setExpandedSections(prev => ({ ...prev, [section]: !prev[section] }));
  };

  const handleCategoryChange = (categoryId) => {
    const updated = selectedCategories.includes(categoryId)
      ? selectedCategories.filter(id => id !== categoryId)
      : [...selectedCategories, categoryId];
    setSelectedCategories(updated);
    applyFilters(updated, priceRange, selectedBrands, certifications);
  };

  const handlePriceChange = (type, value) => {
    const newRange = { ...priceRange, [type]: parseInt(value) || 0 };
    setPriceRange(newRange);
    applyFilters(selectedCategories, newRange, selectedBrands, certifications);
  };

  const handleBrandChange = (brand) => {
    const updated = selectedBrands.includes(brand)
      ? selectedBrands.filter(b => b !== brand)
      : [...selectedBrands, brand];
    setSelectedBrands(updated);
    applyFilters(selectedCategories, priceRange, updated, certifications);
  };

  const handleCertificationChange = (certId) => {
    const updated = certifications.includes(certId)
      ? certifications.filter(c => c !== certId)
      : [...certifications, certId];
    setCertifications(updated);
    applyFilters(selectedCategories, priceRange, selectedBrands, updated);
  };

  const applyFilters = (categories, price, brands, certs) => {
    // In a real app, this would filter from API
    // For dummy data, we'll simulate filtering
    console.log('Applying filters:', { categories, price, brands, certs });
  };

  const resetFilters = () => {
    setSelectedCategories([]);
    setPriceRange({ min: 0, max: 500000 });
    setSelectedBrands([]);
    setCertifications([]);
    applyFilters([], { min: 0, max: 500000 }, [], []);
  };

  return (
    <div className="bg-white rounded-xl shadow-sm p-6 sticky top-24">
      <div className="flex justify-between items-center mb-4">
        <h3 className="font-semibold text-lg">Filter</h3>
        <button onClick={resetFilters} className="text-sm text-primary hover:underline">
          Reset
        </button>
      </div>

      {/* Categories */}
      <div className="border-b pb-4 mb-4">
        <button
          onClick={() => toggleSection('categories')}
          className="flex justify-between items-center w-full font-medium mb-2"
        >
          <span>Kategori</span>
          {expandedSections.categories ? <FaChevronUp /> : <FaChevronDown />}
        </button>
        {expandedSections.categories && (
          <div className="space-y-2 mt-2">
            {categories.map(cat => (
              <label key={cat.id} className="flex items-center gap-2 text-sm">
                <input
                  type="checkbox"
                  checked={selectedCategories.includes(cat.id)}
                  onChange={() => handleCategoryChange(cat.id)}
                  className="rounded text-primary focus:ring-primary"
                />
                <span>{cat.name}</span>
                <span className="text-gray-400 text-xs">({cat.count})</span>
              </label>
            ))}
          </div>
        )}
      </div>

      {/* Price Range */}
      <div className="border-b pb-4 mb-4">
        <button
          onClick={() => toggleSection('price')}
          className="flex justify-between items-center w-full font-medium mb-2"
        >
          <span>Harga</span>
          {expandedSections.price ? <FaChevronUp /> : <FaChevronDown />}
        </button>
        {expandedSections.price && (
          <div className="space-y-3 mt-2">
            <div className="flex gap-3">
              <div className="flex-1">
                <label className="text-xs text-gray-500">Min</label>
                <input
                  type="number"
                  value={priceRange.min}
                  onChange={(e) => handlePriceChange('min', e.target.value)}
                  className="w-full px-3 py-2 border rounded-lg text-sm"
                  placeholder="0"
                />
              </div>
              <div className="flex-1">
                <label className="text-xs text-gray-500">Max</label>
                <input
                  type="number"
                  value={priceRange.max}
                  onChange={(e) => handlePriceChange('max', e.target.value)}
                  className="w-full px-3 py-2 border rounded-lg text-sm"
                  placeholder="500000"
                />
              </div>
            </div>
            <input
              type="range"
              min="0"
              max="500000"
              value={priceRange.max}
              onChange={(e) => handlePriceChange('max', e.target.value)}
              className="w-full"
            />
          </div>
        )}
      </div>

      {/* Brands */}
      <div className="border-b pb-4 mb-4">
        <button
          onClick={() => toggleSection('brands')}
          className="flex justify-between items-center w-full font-medium mb-2"
        >
          <span>Brand</span>
          {expandedSections.brands ? <FaChevronUp /> : <FaChevronDown />}
        </button>
        {expandedSections.brands && (
          <div className="space-y-2 mt-2 max-h-48 overflow-y-auto">
            {brands.map(brand => (
              <label key={brand} className="flex items-center gap-2 text-sm">
                <input
                  type="checkbox"
                  checked={selectedBrands.includes(brand)}
                  onChange={() => handleBrandChange(brand)}
                  className="rounded text-primary focus:ring-primary"
                />
                <span>{brand}</span>
              </label>
            ))}
          </div>
        )}
      </div>

      {/* Certifications */}
      <div className="border-b pb-4 mb-4">
        <button
          onClick={() => toggleSection('certifications')}
          className="flex justify-between items-center w-full font-medium mb-2"
        >
          <span>Sertifikasi</span>
          {expandedSections.certifications ? <FaChevronUp /> : <FaChevronDown />}
        </button>
        {expandedSections.certifications && (
          <div className="space-y-2 mt-2">
            {certificationOptions.map(cert => (
              <label key={cert.id} className="flex items-center gap-2 text-sm">
                <input
                  type="checkbox"
                  checked={certifications.includes(cert.id)}
                  onChange={() => handleCertificationChange(cert.id)}
                  className="rounded text-primary focus:ring-primary"
                />
                <span>{cert.label}</span>
              </label>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

export default ProductFilter;