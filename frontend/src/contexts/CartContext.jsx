import React, { createContext, useContext, useState, useEffect } from 'react';
import toast from 'react-hot-toast';

const CartContext = createContext();

export const useCart = () => {
  const context = useContext(CartContext);
  if (!context) {
    throw new Error('useCart must be used within CartProvider');
  }
  return context;
};

export const CartProvider = ({ children }) => {
  const [cartItems, setCartItems] = useState([]);
  const [isCartOpen, setIsCartOpen] = useState(false);

  // Load cart from localStorage
  useEffect(() => {
    const savedCart = localStorage.getItem('cart');
    if (savedCart) {
      setCartItems(JSON.parse(savedCart));
    }
  }, []);

  // Save cart to localStorage
  useEffect(() => {
    localStorage.setItem('cart', JSON.stringify(cartItems));
  }, [cartItems]);

  const addToCart = (product, variant = null, quantity = 1) => {
    setCartItems((prev) => {
      const existingItem = prev.find(
        (item) => item.id === product.id && item.variantId === (variant?.id || null)
      );

      if (existingItem) {
        const updated = prev.map((item) =>
          item.id === product.id && item.variantId === (variant?.id || null)
            ? { ...item, quantity: item.quantity + quantity }
            : item
        );
        toast.success(`Jumlah ${product.name} ditambahkan`);
        return updated;
      }

      toast.success(`${product.name} ditambahkan ke keranjang`);
      return [
        ...prev,
        {
          id: product.id,
          name: product.name,
          price: product.salePrice || product.basePrice,
          image: product.images?.[0],
          variantName: variant?.name,
          variantId: variant?.id || null,
          quantity,
          maxStock: variant?.stock || product.stock,
        },
      ];
    });
  };

  const removeFromCart = (productId, variantId = null) => {
    setCartItems((prev) =>
      prev.filter((item) => !(item.id === productId && item.variantId === variantId))
    );
    toast.success('Item dihapus dari keranjang');
  };

  const updateQuantity = (productId, variantId, newQuantity) => {
    if (newQuantity < 1) {
      removeFromCart(productId, variantId);
      return;
    }

    setCartItems((prev) =>
      prev.map((item) =>
        item.id === productId && item.variantId === variantId
          ? { ...item, quantity: newQuantity }
          : item
      )
    );
  };

  const clearCart = () => {
    setCartItems([]);
    toast.success('Keranjang dikosongkan');
  };

  const getCartTotal = () => {
    return cartItems.reduce((total, item) => total + item.price * item.quantity, 0);
  };

  const getCartCount = () => {
    return cartItems.reduce((count, item) => count + item.quantity, 0);
  };

  return (
    <CartContext.Provider
      value={{
        cartItems,
        isCartOpen,
        setIsCartOpen,
        addToCart,
        removeFromCart,
        updateQuantity,
        clearCart,
        getCartTotal,
        getCartCount,
      }}
    >
      {children}
    </CartContext.Provider>
  );
};