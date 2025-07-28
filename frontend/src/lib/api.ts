// API configuration and utilities
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export class ApiClient {
  private baseURL: string;
  private userToken: string | null = null;

  constructor() {
    this.baseURL = API_BASE_URL;
    // Get token from localStorage if available
    if (typeof window !== 'undefined') {
      this.userToken = localStorage.getItem('auth_token');
    }
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`;
    
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      ...options.headers,
    };

    // Add auth token if available
    if (this.userToken) {
      headers['Authorization'] = `Bearer ${this.userToken}`;
      headers['X-User-ID'] = this.getUserId() || '1'; // Default user for demo
    }

    const config: RequestInit = {
      ...options,
      headers,
    };

    try {
      const response = await fetch(url, config);
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      return await response.json();
    } catch (error) {
      console.error('API request failed:', error);
      throw error;
    }
  }

  // Authentication methods
  async sendOTP(phoneNumber: string) {
    return this.request('/api/auth/send-otp', {
      method: 'POST',
      body: JSON.stringify({ phone_number: phoneNumber }),
    });
  }

  async verifyOTP(phoneNumber: string, otpCode: string) {
    const response = await this.request<any>('/api/auth/verify-otp', {
      method: 'POST',
      body: JSON.stringify({ phone_number: phoneNumber, otp_code: otpCode }),
    });

    if (response.token) {
      this.userToken = response.token;
      localStorage.setItem('auth_token', response.token);
      localStorage.setItem('user_data', JSON.stringify(response.user));
    }

    return response;
  }

  async logout() {
    await this.request('/api/auth/logout', { method: 'POST' });
    this.userToken = null;
    localStorage.removeItem('auth_token');
    localStorage.removeItem('user_data');
  }

  // User methods
  getUserData() {
    if (typeof window !== 'undefined') {
      const userData = localStorage.getItem('user_data');
      return userData ? JSON.parse(userData) : null;
    }
    return null;
  }

  getUserId() {
    const userData = this.getUserData();
    return userData?.id?.toString() || null;
  }

  isAuthenticated() {
    return !!this.userToken;
  }

  // Rooms methods
  async getRooms() {
    return this.request('/api/rooms');
  }

  async getRoomById(id: string) {
    return this.request(`/api/rooms/${id}`);
  }

  async checkRoomAvailability(roomId: number, checkInDate: string, checkOutDate: string) {
    return this.request('/api/rooms/check-availability', {
      method: 'POST',
      body: JSON.stringify({
        room_id: roomId,
        check_in_date: checkInDate,
        check_out_date: checkOutDate,
      }),
    });
  }

  // Products methods
  async getProducts(category?: string) {
    const url = category ? `/api/products?category=${category}` : '/api/products';
    return this.request(url);
  }

  async getProductById(id: string) {
    return this.request(`/api/products/${id}`);
  }

  async getProductCategories() {
    return this.request('/api/products/categories');
  }

  // Cart methods
  async getCartItems() {
    return this.request('/api/cart');
  }

  async addToCart(item: {
    item_type: string;
    item_id: number;
    quantity: number;
    check_in_date?: string;
    check_out_date?: string;
  }) {
    return this.request('/api/cart/add', {
      method: 'POST',
      body: JSON.stringify(item),
    });
  }

  async removeFromCart(cartItemId: string) {
    return this.request(`/api/cart/${cartItemId}`, {
      method: 'DELETE',
    });
  }

  async clearCart() {
    return this.request('/api/cart/clear', {
      method: 'DELETE',
    });
  }

  // Orders methods
  async getOrders() {
    return this.request('/api/orders');
  }

  async createOrder(orderData: {
    customer_name: string;
    customer_phone: string;
    customer_email?: string;
    notes?: string;
  }) {
    return this.request('/api/orders', {
      method: 'POST',
      body: JSON.stringify(orderData),
    });
  }

  async getOrderById(id: string) {
    return this.request(`/api/orders/${id}`);
  }

  // Admin methods
  async getAdminDashboard() {
    return this.request('/api/admin/dashboard');
  }

  async getAllOrders(status?: string, limit?: number, offset?: number) {
    let url = '/api/admin/orders';
    const params = new URLSearchParams();
    
    if (status) params.append('status', status);
    if (limit) params.append('limit', limit.toString());
    if (offset) params.append('offset', offset.toString());
    
    if (params.toString()) {
      url += `?${params.toString()}`;
    }
    
    return this.request(url);
  }

  async updateOrderStatus(orderId: string, status: string, notes?: string) {
    return this.request(`/api/admin/orders/${orderId}`, {
      method: 'PUT',
      body: JSON.stringify({ status, notes }),
    });
  }
}

// Create a singleton instance
export const apiClient = new ApiClient();

// Types
export interface Room {
  id: number;
  room_number: string;
  room_type: string;
  title: string;
  description: string;
  price_per_night: number;
  max_occupancy: number;
  amenities: string[];
  images: string[];
  is_available: boolean;
  created_at: string;
  updated_at: string;
}

export interface Product {
  id: number;
  name: string;
  description: string;
  price: number;
  category: string;
  stock_quantity: number;
  images: string[];
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface CartItem {
  id: number;
  item_type: string;
  item_id: number;
  item_name: string;
  quantity: number;
  unit_price: number;
  total_price: number;
  images: string[];
  check_in_date?: string;
  check_out_date?: string;
  nights?: number;
}

export interface Order {
  id: number;
  user_id: number;
  order_number: string;
  total_amount: number;
  status: string;
  customer_name: string;
  customer_phone: string;
  customer_email: string;
  notes: string;
  created_at: string;
  updated_at: string;
  items?: OrderItem[];
}

export interface OrderItem {
  id: number;
  order_id: number;
  item_type: string;
  item_id: number;
  item_name: string;
  quantity: number;
  unit_price: number;
  total_price: number;
  check_in_date?: string;
  check_out_date?: string;
  nights?: number;
  created_at: string;
}
