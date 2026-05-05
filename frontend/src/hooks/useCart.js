import { useCart } from '../contexts/CartContext';

// Re-export untuk convenience
export { useCart };

// Custom hook untuk cart calculations
export const useCartTotals = () => {
  const { cartItems, getCartTotal } = useCart();
  
  const subtotal = getCartTotal();
  const shippingCost = subtotal > 150000 ? 0 : 20000;
  const discount = 0;
  const tax = 0;
  const total = subtotal + shippingCost + tax - discount;
  
  const itemCount = cartItems.reduce((sum, item) => sum + item.quantity, 0);
  
  const isEligibleForFreeShipping = subtotal >= 150000;
  
  return {
    subtotal,
    shippingCost,
    discount,
    tax,
    total,
    itemCount,
    isEligibleForFreeShipping,
    savingFromShipping: isEligibleForFreeShipping ? 20000 : 0,
  };
};

// Hook untuk cart actions
export const useCartActions = () => {
  const { addToCart, removeFromCart, updateQuantity, clearCart } = useCart();
  
  const addItem = (product, variant = null, quantity = 1) => {
    // Check stock availability
    const stock = variant?.stock || product.stock;
    if (quantity > stock) {
      console.warn('Stock tidak mencukupi');
      return false;
    }
    addToCart(product, variant, quantity);
    return true;
  };
  
  const updateItemQuantity = (productId, variantId, newQuantity, maxStock) => {
    if (newQuantity > maxStock) {
      console.warn('Stock tidak mencukupi');
      return false;
    }
    updateQuantity(productId, variantId, newQuantity);
    return true;
  };
  
  return {
    addItem,
    removeItem: removeFromCart,
    updateItemQuantity,
    clearCart,
  };
};