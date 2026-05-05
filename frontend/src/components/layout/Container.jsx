import React from 'react';

const Container = ({ children, className = '', size = 'default', noPadding = false }) => {
  const sizeClasses = {
    sm: 'max-w-5xl',
    default: 'max-w-7xl',
    lg: 'max-w-[90rem]',
    full: 'max-w-full',
  };

  return (
    <div
      className={`
        mx-auto
        ${sizeClasses[size]}
        ${!noPadding && 'px-4 sm:px-6 lg:px-8'}
        ${className}
      `}
    >
      {children}
    </div>
  );
};

export default Container;