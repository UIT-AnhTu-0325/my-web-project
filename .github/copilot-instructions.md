<!-- Use this file to provide workspace-specific custom instructions to Copilot. For more details, visit https://code.visualstudio.com/docs/copilot/copilot-customization#_use-a-githubcopilotinstructionsmd-file -->

# Hotel Booking & E-commerce System - Copilot Instructions

## Project Overview
This is a full-stack hotel booking and e-commerce system with the following architecture:
- **Frontend**: Next.js with TypeScript (React-based)
- **Backend**: Go (Golang) REST API server
- **Database**: PostgreSQL
- **Python Scripts**: Email service and data processing utilities

## Code Style Guidelines

### Go Backend
- Use clean architecture patterns with clear separation of concerns
- Follow standard Go naming conventions (PascalCase for exported, camelCase for unexported)
- Use Gin framework for HTTP routing and middleware
- Implement proper error handling with meaningful error messages
- Use structured logging
- Follow REST API conventions for endpoints
- Use database/sql with lib/pq driver for PostgreSQL
- Implement JWT-based authentication for API security

### TypeScript Frontend
- Use functional components with React hooks
- Implement proper TypeScript typing for all props and state
- Use Tailwind CSS for styling
- Follow component-based architecture
- Implement proper error boundaries and loading states
- Use Next.js App Router for routing
- Implement responsive design principles

### Python Scripts
- Follow PEP 8 style guidelines
- Use type hints where appropriate
- Implement proper error handling and logging
- Use Flask for HTTP APIs
- Document functions with clear docstrings
- Use environment variables for configuration
- Implement data validation for all inputs

## Database Schema
- Users table with phone-based authentication
- Rooms table for hotel inventory
- Products table for e-commerce items
- Orders and OrderItems for purchase management
- CartItems for shopping cart functionality
- RoomBookings for reservation tracking
- OTPs for authentication

## API Endpoints Structure
- `/api/auth/*` - Authentication endpoints
- `/api/rooms/*` - Room management
- `/api/products/*` - Product catalog
- `/api/cart/*` - Shopping cart operations
- `/api/orders/*` - Order management
- `/api/admin/*` - Administrative functions

## Key Features to Implement
1. Phone-based OTP authentication
2. Room booking with date selection
3. Product catalog and shopping cart
4. Order processing and management
5. Admin dashboard for inventory and orders
6. Email notifications via Python service
7. Data analytics and reporting

## Learning Focus Areas
- **Go**: REST API development, database operations, middleware, JWT authentication
- **Python**: Email services, data processing, Flask APIs, CSV/JSON handling
- **TypeScript**: React components, state management, API integration, type safety
- **PostgreSQL**: Relational database design, queries, migrations

## Integration Points
- Go backend serves REST API on port 8080
- Python email service runs on port 8001
- Frontend communicates with Go backend via HTTP
- Go backend calls Python email service for notifications
- All services use environment variables for configuration

## Security Considerations
- Implement proper input validation
- Use parameterized queries to prevent SQL injection
- Implement rate limiting for authentication endpoints
- Secure JWT token handling
- Validate and sanitize all user inputs
- Use HTTPS in production
- Implement proper CORS configuration
