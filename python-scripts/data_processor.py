#!/usr/bin/env python3
"""
Data Processor for Hotel E-commerce System

This Python script handles data processing tasks for the hotel booking system.
It provides utilities for order analytics, data validation, and report generation.

Features:
- Order analytics and statistics
- Sales reporting
- Data validation utilities
- CSV export functionality
- Room occupancy calculations

Learning Goals:
- Data analysis with Python
- Working with JSON and CSV files
- Statistics and analytics
- File I/O operations
- Data validation techniques
"""

import json
import csv
import os
from datetime import datetime, timedelta
from collections import defaultdict, Counter
import statistics

class DataProcessor:
    """Main class for data processing operations"""
    
    def __init__(self, data_directory='../backend/data'):
        self.data_directory = data_directory
        self.orders_file = os.path.join(data_directory, 'orders.json')
        self.rooms_file = os.path.join(data_directory, 'rooms.json')
        self.products_file = os.path.join(data_directory, 'products.json')
    
    def load_json_file(self, filepath):
        """
        Load data from a JSON file
        
        Args:
            filepath (str): Path to the JSON file
            
        Returns:
            list: Data from the file, or empty list if file doesn't exist
        """
        try:
            if os.path.exists(filepath):
                with open(filepath, 'r') as f:
                    return json.load(f)
            return []
        except Exception as e:
            print(f"Error loading {filepath}: {str(e)}")
            return []
    
    def save_json_file(self, filepath, data):
        """
        Save data to a JSON file
        
        Args:
            filepath (str): Path to save the file
            data (list): Data to save
        """
        try:
            os.makedirs(os.path.dirname(filepath), exist_ok=True)
            with open(filepath, 'w') as f:
                json.dump(data, f, indent=2, default=str)
            print(f"Data saved to {filepath}")
        except Exception as e:
            print(f"Error saving {filepath}: {str(e)}")
    
    def validate_order_data(self, order):
        """
        Validate order data structure
        
        Args:
            order (dict): Order data to validate
            
        Returns:
            tuple: (is_valid, error_messages)
        """
        errors = []
        required_fields = ['order_number', 'customer_name', 'customer_phone', 'total_amount', 'items']
        
        # Check required fields
        for field in required_fields:
            if field not in order or not order[field]:
                errors.append(f"Missing required field: {field}")
        
        # Validate phone number format
        if 'customer_phone' in order:
            phone = order['customer_phone']
            if not (phone.startswith('+') or phone.startswith('0')) or len(phone) < 10:
                errors.append("Invalid phone number format")
        
        # Validate total amount
        if 'total_amount' in order:
            try:
                amount = float(order['total_amount'])
                if amount <= 0:
                    errors.append("Total amount must be positive")
            except (ValueError, TypeError):
                errors.append("Total amount must be a valid number")
        
        # Validate items
        if 'items' in order and isinstance(order['items'], list):
            if len(order['items']) == 0:
                errors.append("Order must contain at least one item")
            
            for i, item in enumerate(order['items']):
                if 'item_name' not in item or not item['item_name']:
                    errors.append(f"Item {i+1}: Missing item name")
                if 'quantity' not in item or item['quantity'] <= 0:
                    errors.append(f"Item {i+1}: Invalid quantity")
                if 'total_price' not in item or item['total_price'] <= 0:
                    errors.append(f"Item {i+1}: Invalid price")
        
        return len(errors) == 0, errors
    
    def generate_order_analytics(self, start_date=None, end_date=None):
        """
        Generate analytics for orders within a date range
        
        Args:
            start_date (str): Start date in YYYY-MM-DD format
            end_date (str): End date in YYYY-MM-DD format
            
        Returns:
            dict: Analytics data
        """
        orders = self.load_json_file(self.orders_file)
        
        if not orders:
            return {'error': 'No orders found'}
        
        # Filter by date range if specified
        if start_date or end_date:
            filtered_orders = []
            for order in orders:
                order_date = order.get('created_at', '')[:10]  # Extract date part
                if start_date and order_date < start_date:
                    continue
                if end_date and order_date > end_date:
                    continue
                filtered_orders.append(order)
            orders = filtered_orders
        
        if not orders:
            return {'error': 'No orders found in the specified date range'}
        
        # Calculate analytics
        total_orders = len(orders)
        total_revenue = sum(float(order.get('total_amount', 0)) for order in orders)
        average_order_value = total_revenue / total_orders if total_orders > 0 else 0
        
        # Order status distribution
        status_counts = Counter(order.get('status', 'unknown') for order in orders)
        
        # Revenue by day
        daily_revenue = defaultdict(float)
        for order in orders:
            date = order.get('created_at', '')[:10]
            daily_revenue[date] += float(order.get('total_amount', 0))
        
        # Popular items
        item_counts = defaultdict(int)
        item_revenue = defaultdict(float)
        
        for order in orders:
            for item in order.get('items', []):
                item_name = item.get('item_name', 'Unknown')
                quantity = item.get('quantity', 1)
                price = float(item.get('total_price', 0))
                
                item_counts[item_name] += quantity
                item_revenue[item_name] += price
        
        # Top selling items
        top_items_by_quantity = sorted(item_counts.items(), key=lambda x: x[1], reverse=True)[:10]
        top_items_by_revenue = sorted(item_revenue.items(), key=lambda x: x[1], reverse=True)[:10]
        
        # Room vs Product revenue
        room_revenue = 0
        product_revenue = 0
        
        for order in orders:
            for item in order.get('items', []):
                price = float(item.get('total_price', 0))
                if item.get('item_type') == 'room':
                    room_revenue += price
                else:
                    product_revenue += price
        
        return {
            'summary': {
                'total_orders': total_orders,
                'total_revenue': round(total_revenue, 2),
                'average_order_value': round(average_order_value, 2),
                'room_revenue': round(room_revenue, 2),
                'product_revenue': round(product_revenue, 2)
            },
            'order_status_distribution': dict(status_counts),
            'daily_revenue': dict(daily_revenue),
            'top_items_by_quantity': top_items_by_quantity,
            'top_items_by_revenue': top_items_by_revenue,
            'date_range': {
                'start': start_date or 'all time',
                'end': end_date or 'all time'
            }
        }
    
    def calculate_room_occupancy(self, start_date, end_date):
        """
        Calculate room occupancy rate for a date range
        
        Args:
            start_date (str): Start date in YYYY-MM-DD format
            end_date (str): End date in YYYY-MM-DD format
            
        Returns:
            dict: Occupancy data
        """
        rooms = self.load_json_file(self.rooms_file)
        orders = self.load_json_file(self.orders_file)
        
        if not rooms:
            return {'error': 'No rooms data found'}
        
        total_rooms = len(rooms)
        
        # Calculate total room-nights in the period
        start = datetime.strptime(start_date, '%Y-%m-%d')
        end = datetime.strptime(end_date, '%Y-%m-%d')
        total_days = (end - start).days + 1
        total_room_nights = total_rooms * total_days
        
        # Count booked room-nights
        booked_room_nights = 0
        room_bookings = defaultdict(list)
        
        for order in orders:
            if order.get('status') in ['confirmed', 'completed']:
                for item in order.get('items', []):
                    if item.get('item_type') == 'room':
                        check_in = item.get('check_in_date')
                        check_out = item.get('check_out_date')
                        nights = item.get('nights', 1)
                        
                        if check_in and check_out:
                            # Check if booking overlaps with our date range
                            booking_start = max(datetime.strptime(check_in, '%Y-%m-%d'), start)
                            booking_end = min(datetime.strptime(check_out, '%Y-%m-%d'), end)
                            
                            if booking_start <= booking_end:
                                overlap_nights = (booking_end - booking_start).days
                                booked_room_nights += overlap_nights
                                
                                room_bookings[item.get('item_name', 'Unknown')].append({
                                    'check_in': check_in,
                                    'check_out': check_out,
                                    'nights': nights,
                                    'order_number': order.get('order_number')
                                })
        
        occupancy_rate = (booked_room_nights / total_room_nights * 100) if total_room_nights > 0 else 0
        
        return {
            'period': {
                'start_date': start_date,
                'end_date': end_date,
                'total_days': total_days
            },
            'rooms': {
                'total_rooms': total_rooms,
                'total_room_nights_available': total_room_nights,
                'booked_room_nights': booked_room_nights,
                'occupancy_rate': round(occupancy_rate, 2)
            },
            'room_bookings': dict(room_bookings)
        }
    
    def export_orders_to_csv(self, filename='orders_export.csv', start_date=None, end_date=None):
        """
        Export orders to CSV file
        
        Args:
            filename (str): Output filename
            start_date (str): Start date filter
            end_date (str): End date filter
        """
        orders = self.load_json_file(self.orders_file)
        
        # Filter by date if specified
        if start_date or end_date:
            filtered_orders = []
            for order in orders:
                order_date = order.get('created_at', '')[:10]
                if start_date and order_date < start_date:
                    continue
                if end_date and order_date > end_date:
                    continue
                filtered_orders.append(order)
            orders = filtered_orders
        
        try:
            with open(filename, 'w', newline='', encoding='utf-8') as csvfile:
                fieldnames = [
                    'order_number', 'customer_name', 'customer_phone', 'customer_email',
                    'total_amount', 'status', 'created_at', 'item_name', 'item_type',
                    'quantity', 'unit_price', 'check_in_date', 'check_out_date', 'nights'
                ]
                
                writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
                writer.writeheader()
                
                for order in orders:
                    for item in order.get('items', []):
                        row = {
                            'order_number': order.get('order_number', ''),
                            'customer_name': order.get('customer_name', ''),
                            'customer_phone': order.get('customer_phone', ''),
                            'customer_email': order.get('customer_email', ''),
                            'total_amount': order.get('total_amount', 0),
                            'status': order.get('status', ''),
                            'created_at': order.get('created_at', ''),
                            'item_name': item.get('item_name', ''),
                            'item_type': item.get('item_type', ''),
                            'quantity': item.get('quantity', 1),
                            'unit_price': item.get('unit_price', 0),
                            'check_in_date': item.get('check_in_date', ''),
                            'check_out_date': item.get('check_out_date', ''),
                            'nights': item.get('nights', '')
                        }
                        writer.writerow(row)
            
            print(f"Orders exported to {filename}")
            return True
            
        except Exception as e:
            print(f"Error exporting to CSV: {str(e)}")
            return False

def main():
    """Main function for command-line usage"""
    print("🔄 Hotel E-commerce Data Processor")
    print("=" * 50)
    
    processor = DataProcessor()
    
    while True:
        print("\nAvailable operations:")
        print("1. Generate order analytics")
        print("2. Calculate room occupancy")
        print("3. Export orders to CSV")
        print("4. Validate order data")
        print("5. Exit")
        
        choice = input("\nSelect an option (1-5): ").strip()
        
        if choice == '1':
            start_date = input("Start date (YYYY-MM-DD, or press Enter for all time): ").strip() or None
            end_date = input("End date (YYYY-MM-DD, or press Enter for all time): ").strip() or None
            
            analytics = processor.generate_order_analytics(start_date, end_date)
            print("\n📊 Order Analytics:")
            print(json.dumps(analytics, indent=2))
        
        elif choice == '2':
            start_date = input("Start date (YYYY-MM-DD): ").strip()
            end_date = input("End date (YYYY-MM-DD): ").strip()
            
            if start_date and end_date:
                occupancy = processor.calculate_room_occupancy(start_date, end_date)
                print("\n🏨 Room Occupancy Analysis:")
                print(json.dumps(occupancy, indent=2))
            else:
                print("❌ Both start and end dates are required for occupancy calculation")
        
        elif choice == '3':
            filename = input("CSV filename (default: orders_export.csv): ").strip() or 'orders_export.csv'
            start_date = input("Start date (YYYY-MM-DD, or press Enter for all): ").strip() or None
            end_date = input("End date (YYYY-MM-DD, or press Enter for all): ").strip() or None
            
            success = processor.export_orders_to_csv(filename, start_date, end_date)
            if success:
                print(f"✅ Orders exported to {filename}")
            else:
                print("❌ Export failed")
        
        elif choice == '4':
            # Sample order validation
            sample_order = {
                'order_number': 'ORD-001',
                'customer_name': 'John Doe',
                'customer_phone': '+1234567890',
                'total_amount': 199.99,
                'items': [
                    {
                        'item_name': 'Deluxe Room',
                        'item_type': 'room',
                        'quantity': 1,
                        'total_price': 149.99
                    }
                ]
            }
            
            is_valid, errors = processor.validate_order_data(sample_order)
            print(f"\n✅ Sample order validation: {'Valid' if is_valid else 'Invalid'}")
            if errors:
                for error in errors:
                    print(f"❌ {error}")
        
        elif choice == '5':
            print("👋 Goodbye!")
            break
        
        else:
            print("❌ Invalid option. Please try again.")

if __name__ == '__main__':
    main()
