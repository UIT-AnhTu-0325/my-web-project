'use client';

import { useEffect, useState } from 'react';
import { apiClient } from '@/lib/api-client';
import { Product } from '@/types/api';
import ProductCard from '@/components/ProductCard';
import LoadingSpinner from '@/components/LoadingSpinner';

export default function ProductsPage() {
  const [products, setProducts] = useState<Product[]>([]);
  const [rawProducts, setRawProducts] = useState<unknown>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchProducts = async () => {
      try {
        const response = await apiClient.get('/api/products');
        setRawProducts(response.data);
        if (Array.isArray(response.data)) {
          setProducts(response.data);
        } else if (Array.isArray(response.data?.products)) {
          setProducts(response.data.products);
        } else {
          setProducts([]);
        }
        console.log('Products API response:', response.data);
      } catch (error) {
        console.error('Error fetching products:', error);
      } finally {
        setLoading(false);
      }
    };
    fetchProducts();
  }, []);

  if (loading) {
    return <LoadingSpinner />;
  }

  return (
    <div className="py-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center mb-12">
          <h1 className="text-4xl font-bold text-gray-900 mb-4">
            Our Products
          </h1>
          <p className="text-lg text-gray-600 max-w-3xl mx-auto">
            Discover our curated collection of premium products. 
            From essentials to luxury items, we have everything you need for a perfect experience.
          </p>
        </div>

        {Array.isArray(products) ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
            {products.map((product) => (
              <ProductCard key={product.id} product={product} />
            ))}
          </div>
        ) : (
          <div className="text-center py-12">
            <p className="text-red-500 text-lg">Products data is not an array. Check API response.</p>
            <pre className="text-xs text-gray-500 overflow-x-auto bg-gray-100 p-2 rounded mt-2">{JSON.stringify(rawProducts, null, 2)}</pre>
          </div>
        )}

        {Array.isArray(products) && products.length === 0 && !loading && (
          <div className="text-center py-12">
            <p className="text-gray-500 text-lg">No products available at the moment.</p>
          </div>
        )}
      </div>
    </div>
  );
}
