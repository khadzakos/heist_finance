import React from 'react';

interface LoadingSpinnerProps {
  size?: 'small' | 'medium' | 'large';
  color?: string;
}

const LoadingSpinner: React.FC<LoadingSpinnerProps> = ({
  size = 'medium',
  color = '#3b82f6' // blue-500 default
}) => {
  // Determine spinner size
  const sizeMap = {
    small: 'w-4 h-4',
    medium: 'w-8 h-8',
    large: 'w-12 h-12'
  };
  
  const sizeClass = sizeMap[size];
  
  return (
    <div className="flex justify-center items-center">
      <div 
        className={`${sizeClass} border-4 rounded-full animate-spin`}
        style={{
          borderColor: `${color} transparent transparent transparent`,
          borderTopColor: color
        }}
      ></div>
    </div>
  );
};

export default LoadingSpinner; 