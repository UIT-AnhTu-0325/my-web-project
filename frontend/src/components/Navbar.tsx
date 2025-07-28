'use client';

import { useState } from 'react';
import Link from 'next/link';
import { apiClient } from '@/lib/api';

export default function Navbar() {
  const [user, setUser] = useState(apiClient.getUserData());
  const [isMenuOpen, setIsMenuOpen] = useState(false);

  const handleLogout = async () => {
    try {
      await apiClient.logout();
      setUser(null);
      window.location.href = '/';
    } catch (error) {
      console.error('Logout failed:', error);
    }
  };

  return (
    <nav className="bg-white shadow-lg">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between h-16">
          <div className="flex items-center">
            <Link href="/" className="flex-shrink-0">
              <h1 className="text-2xl font-bold text-blue-600">🏨 HotelCommerce</h1>
            </Link>
          </div>

          {/* Desktop Menu */}
          <div className="hidden md:flex items-center space-x-8">
            <Link href="/rooms" className="text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium">
              Rooms
            </Link>
            <Link href="/products" className="text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium">
              Products
            </Link>
            <Link href="/cart" className="text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium">
              Cart
            </Link>
            
            {user ? (
              <div className="flex items-center space-x-4">
                <Link href="/orders" className="text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium">
                  My Orders
                </Link>
                {user.is_admin && (
                  <Link href="/admin" className="text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium">
                    Admin
                  </Link>
                )}
                <div className="flex items-center space-x-2">
                  <span className="text-sm text-gray-600">Hello, {user.name}</span>
                  <button
                    onClick={handleLogout}
                    className="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded text-sm"
                  >
                    Logout
                  </button>
                </div>
              </div>
            ) : (
              <Link href="/auth/login" className="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-md text-sm font-medium">
                Login
              </Link>
            )}
          </div>

          {/* Mobile menu button */}
          <div className="md:hidden flex items-center">
            <button
              onClick={() => setIsMenuOpen(!isMenuOpen)}
              className="text-gray-700 hover:text-blue-600 focus:outline-none focus:text-blue-600"
            >
              <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
              </svg>
            </button>
          </div>
        </div>

        {/* Mobile Menu */}
        {isMenuOpen && (
          <div className="md:hidden">
            <div className="px-2 pt-2 pb-3 space-y-1 sm:px-3">
              <Link href="/rooms" className="block text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-base font-medium">
                Rooms
              </Link>
              <Link href="/products" className="block text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-base font-medium">
                Products
              </Link>
              <Link href="/cart" className="block text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-base font-medium">
                Cart
              </Link>
              {user ? (
                <>
                  <Link href="/orders" className="block text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-base font-medium">
                    My Orders
                  </Link>
                  {user.is_admin && (
                    <Link href="/admin" className="block text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-base font-medium">
                      Admin
                    </Link>
                  )}
                  <div className="px-3 py-2">
                    <span className="text-sm text-gray-600">Hello, {user.name}</span>
                    <button
                      onClick={handleLogout}
                      className="ml-2 bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded text-sm"
                    >
                      Logout
                    </button>
                  </div>
                </>
              ) : (
                <Link href="/auth/login" className="block bg-blue-500 hover:bg-blue-600 text-white px-3 py-2 rounded-md text-base font-medium">
                  Login
                </Link>
              )}
            </div>
          </div>
        )}
      </div>
    </nav>
  );
}
