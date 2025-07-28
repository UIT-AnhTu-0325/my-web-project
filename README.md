# Hotel Booking & E-commerce System

A full-stack application for hotel room booking and product sales with multi-language architecture.

## Architecture

- **Frontend**: Next.js with TypeScript
- **Backend**: Go (Golang) REST API
- **Database**: PostgreSQL
- **Python Scripts**: Email service and data processing

## Features

### Customer Features
- 📱 Phone-based authentication with OTP
- 🏨 Browse and book hotel rooms
- 🛍️ Product catalog and shopping cart
- 📝 Order management and history
- 📧 Email notifications

### Admin Features
- 📊 Admin dashboard
- 🏨 Room management
- 📦 Product inventory management
- 📋 Order processing
- 👥 Customer management

## Project Structure

```
├── frontend/           # Next.js TypeScript application
├── backend/            # Go REST API server
├── python-scripts/     # Python utilities and services
└── docs/              # Documentation
```

## Getting Started

### Prerequisites
- Node.js 18+
- Go 1.21+
- Python 3.9+
- PostgreSQL 14+

### Installation

1. **Database Setup**
   - Create a PostgreSQL database named: `hotel_ecommerce`

2. **Backend Setup**
```bash
cd backend
# Copy .env.example to .env and configure database connection
go mod tidy
go run cmd/migrate/main.go  # Run database migrations
go run cmd/server/main.go   # Start API server
```
2. **Frontend Setup**
```bash
cd frontend
npm install
npm run dev
```

3. **Python Scripts Setup**
npm run dev
```

3. **Backend Setup**
```bash
cd backend
go mod init hotel-backend
# Copy .env.example to .env and configure database connection
go run cmd/server/main.go
```

4. **Python Scripts Setup**
```bash
cd python-scripts
pip install -r requirements.txt
python email_service.py
```

## API Endpoints

### Authentication
- `POST /api/auth/send-otp` - Send OTP to phone
- `POST /api/auth/verify-otp` - Verify OTP and login
- `POST /api/auth/logout` - Logout user

### Rooms
- `GET /api/rooms` - Get all rooms
- `GET /api/rooms/:id` - Get room details
- `POST /api/rooms/book` - Book a room

### Products
- `GET /api/products` - Get all products
- `GET /api/products/:id` - Get product details

### Cart & Orders
- `GET /api/cart` - Get cart items
- `POST /api/cart/add` - Add item to cart
- `POST /api/orders` - Create order
- `GET /api/orders` - Get user orders

### Admin
- `GET /api/admin/orders` - Get all orders
- `PUT /api/admin/orders/:id` - Update order status
- `POST /api/admin/rooms` - Add new room
- `POST /api/admin/products` - Add new product

## Development

### Running in Development
```bash
# Terminal 1 - Frontend
cd frontend && npm run dev

# Terminal 2 - Backend
cd backend && go run cmd/server/main.go

# Terminal 3 - Python Services
cd python-scripts && python email_service.py
```

## Learning Goals

This project is designed to help learn:
- **Go**: REST API development, JSON handling, middleware
- **Python**: Email services, data processing, file operations
- **TypeScript**: React components, API integration, type safety
- **Full-stack**: Communication between different technologies

## License

MIT License
