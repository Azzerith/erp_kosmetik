import React, { createContext, useContext, useState, useEffect, useCallback } from 'react';
import toast from 'react-hot-toast';

const AuthContext = createContext();

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider');
  }
  return context;
};

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [token, setToken] = useState(null);

  // Load user from localStorage on mount
  useEffect(() => {
    const storedUser = localStorage.getItem('user');
    const storedToken = localStorage.getItem('token');
    
    if (storedUser && storedToken) {
      setUser(JSON.parse(storedUser));
      setToken(storedToken);
    }
    setLoading(false);
  }, []);

  // Login function
  const login = useCallback(async (email, password, remember = false) => {
    setLoading(true);
    try {
      // Simulasi API call
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      // Dummy login - in real app, call actual API
      if (email === 'demo@example.com' && password === 'demo123') {
        const userData = {
          id: 1,
          name: 'Demo User',
          email: email,
          role: 'customer',
          avatar: 'https://randomuser.me/api/portraits/men/1.jpg',
        };
        const authToken = 'dummy-jwt-token-' + Date.now();
        
        setUser(userData);
        setToken(authToken);
        
        if (remember) {
          localStorage.setItem('user', JSON.stringify(userData));
          localStorage.setItem('token', authToken);
        }
        
        toast.success('Login berhasil!');
        return { success: true };
      }
      
      throw new Error('Email atau password salah');
    } catch (error) {
      toast.error(error.message);
      return { success: false, error: error.message };
    } finally {
      setLoading(false);
    }
  }, []);

  // Register function
  const register = useCallback(async (userData) => {
    setLoading(true);
    try {
      // Simulasi API call
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      const newUser = {
        id: Date.now(),
        name: userData.name,
        email: userData.email,
        phone: userData.phone,
        role: 'customer',
        avatar: null,
      };
      
      // In real app, this would call registration API
      toast.success('Pendaftaran berhasil! Silakan login.');
      return { success: true };
    } catch (error) {
      toast.error(error.message);
      return { success: false, error: error.message };
    } finally {
      setLoading(false);
    }
  }, []);

  // Google Login
  const loginWithGoogle = useCallback(async () => {
    setLoading(true);
    try {
      // Simulasi OAuth flow
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      const userData = {
        id: 2,
        name: 'Google User',
        email: 'google.user@gmail.com',
        role: 'customer',
        avatar: 'https://randomuser.me/api/portraits/women/2.jpg',
        provider: 'google',
      };
      const authToken = 'google-jwt-token-' + Date.now();
      
      setUser(userData);
      setToken(authToken);
      localStorage.setItem('user', JSON.stringify(userData));
      localStorage.setItem('token', authToken);
      
      toast.success('Login dengan Google berhasil!');
      return { success: true };
    } catch (error) {
      toast.error(error.message);
      return { success: false };
    } finally {
      setLoading(false);
    }
  }, []);

  // Facebook Login
  const loginWithFacebook = useCallback(async () => {
    setLoading(true);
    try {
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      const userData = {
        id: 3,
        name: 'Facebook User',
        email: 'fb.user@facebook.com',
        role: 'customer',
        avatar: 'https://randomuser.me/api/portraits/men/3.jpg',
        provider: 'facebook',
      };
      const authToken = 'facebook-jwt-token-' + Date.now();
      
      setUser(userData);
      setToken(authToken);
      localStorage.setItem('user', JSON.stringify(userData));
      localStorage.setItem('token', authToken);
      
      toast.success('Login dengan Facebook berhasil!');
      return { success: true };
    } catch (error) {
      toast.error(error.message);
      return { success: false };
    } finally {
      setLoading(false);
    }
  }, []);

  // Logout function
  const logout = useCallback(() => {
    setUser(null);
    setToken(null);
    localStorage.removeItem('user');
    localStorage.removeItem('token');
    toast.success('Berhasil logout');
  }, []);

  // Update profile
  const updateProfile = useCallback(async (updates) => {
    setLoading(true);
    try {
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      setUser(prev => ({ ...prev, ...updates }));
      localStorage.setItem('user', JSON.stringify({ ...user, ...updates }));
      
      toast.success('Profil berhasil diperbarui');
      return { success: true };
    } catch (error) {
      toast.error(error.message);
      return { success: false };
    } finally {
      setLoading(false);
    }
  }, [user]);

  // Change password
  const changePassword = useCallback(async (oldPassword, newPassword) => {
    setLoading(true);
    try {
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      toast.success('Password berhasil diubah');
      return { success: true };
    } catch (error) {
      toast.error(error.message);
      return { success: false };
    } finally {
      setLoading(false);
    }
  }, []);

  const value = {
    user,
    loading,
    token,
    isAuthenticated: !!user,
    isAdmin: user?.role === 'admin' || user?.role === 'super_admin',
    login,
    register,
    loginWithGoogle,
    loginWithFacebook,
    logout,
    updateProfile,
    changePassword,
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};