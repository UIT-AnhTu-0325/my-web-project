export interface Room {
  id: number;
  title: string;
  description: string;
  price_per_night: number;
  max_occupancy: number;
  room_type: string;
  amenities: string[];
  images: string[];
  room_number: string;
  is_available: boolean;
  created_at: string;
  updated_at: string;
}

export interface Product {
  id: number;
  name: string;
  description: string;
  price: number;
  stock_quantity: number;
  images: string[];
  category: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}
