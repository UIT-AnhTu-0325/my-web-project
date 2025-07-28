
import { useState } from 'react';
import { Product } from '@/lib/api';
import Image from 'next/image';

interface ProductCardProps {
  product: Product;
  onAddToCart?: (product: Product) => void;
}

export default function ProductCard({ product, onAddToCart }: ProductCardProps) {
  const [imgSrc, setImgSrc] = useState(
    product.images && product.images.length > 0
      ? `/images/products/${product.images[0]}`
      : '/images/placeholder-product.jpg'
  );
  return (
    <div className="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow">
      <div className="relative h-48">
        {imgSrc ? (
          <Image
            src={imgSrc}
            alt={product.name}
            fill
            className="object-cover"
            onError={() => setImgSrc('/images/placeholder-product.jpg')}
          />
        ) : (
          <div className="w-full h-full bg-gray-200 flex items-center justify-center">
            <span className="text-gray-500 text-4xl">📦</span>
          </div>
        )}
        <div className="absolute top-2 right-2">
          <span className="bg-green-600 text-white px-2 py-1 rounded text-sm">
            {product.category}
          </span>
        </div>
      </div>
      
      <div className="p-4">
        <h3 className="text-lg font-semibold text-gray-900 mb-2">{product.name}</h3>
        <p className="text-gray-600 text-sm mb-3 line-clamp-2">{product.description}</p>
        
        <div className="flex items-center justify-between mb-3">
          <div className="text-2xl font-bold text-green-600">
            ${product.price}
          </div>
          <div className="text-sm text-gray-500">
            Stock: {product.stock_quantity}
          </div>
        </div>

        <button
          onClick={() => onAddToCart?.(product)}
          disabled={product.stock_quantity === 0 || !onAddToCart}
          className={`w-full py-2 px-4 rounded transition-colors ${
            product.stock_quantity === 0 || !onAddToCart
              ? 'bg-gray-300 text-gray-500 cursor-not-allowed'
              : 'bg-blue-600 text-white hover:bg-blue-700'
          }`}
        >
          {product.stock_quantity === 0 ? 'Out of Stock' : 'Add to Cart'}
        </button>
      </div>
    </div>
  );
}
