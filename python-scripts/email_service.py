#!/usr/bin/env python3
"""
Email Service for Hotel E-commerce System

This Python script handles email notifications for the hotel booking system.
It provides a simple HTTP API that the Go backend can call to send emails.

Features:
- Send order confirmation emails to customers
- Send order notification emails to admin
- Template-based email system
- Simple HTTP API using Flask

Learning Goals:
- HTTP server development in Python
- Email handling with smtplib
- Template rendering with Jinja2
- JSON API development
- Environment variable management
"""

import os
import smtplib
import json
from datetime import datetime
from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart
from flask import Flask, request, jsonify
from jinja2 import Template

app = Flask(__name__)

# Email configuration
SMTP_SERVER = os.getenv('SMTP_SERVER', 'smtp.gmail.com')
SMTP_PORT = int(os.getenv('SMTP_PORT', '587'))
EMAIL_USERNAME = os.getenv('EMAIL_USERNAME', '')
EMAIL_PASSWORD = os.getenv('EMAIL_PASSWORD', '')
FROM_EMAIL = os.getenv('FROM_EMAIL', 'noreply@hotel.com')
ADMIN_EMAIL = os.getenv('ADMIN_EMAIL', 'admin@hotel.com')

def send_email(to_email, subject, html_content, text_content=None):
    """
    Send an email using SMTP
    
    Args:
        to_email (str): Recipient email address
        subject (str): Email subject
        html_content (str): HTML content of the email
        text_content (str): Plain text content (optional)
    
    Returns:
        bool: True if email sent successfully, False otherwise
    """
    try:
        # Create message
        msg = MIMEMultipart('alternative')
        msg['Subject'] = subject
        msg['From'] = FROM_EMAIL
        msg['To'] = to_email
        
        # Add text content
        if text_content:
            text_part = MIMEText(text_content, 'plain')
            msg.attach(text_part)
        
        # Add HTML content
        html_part = MIMEText(html_content, 'html')
        msg.attach(html_part)
        
        # Send email
        with smtplib.SMTP(SMTP_SERVER, SMTP_PORT) as server:
            server.starttls()
            if EMAIL_USERNAME and EMAIL_PASSWORD:
                server.login(EMAIL_USERNAME, EMAIL_PASSWORD)
            server.send_message(msg)
        
        print(f"Email sent successfully to {to_email}")
        return True
        
    except Exception as e:
        print(f"Failed to send email to {to_email}: {str(e)}")
        return False

def render_order_confirmation_email(order_data):
    """
    Render order confirmation email template
    
    Args:
        order_data (dict): Order information
        
    Returns:
        tuple: (html_content, text_content)
    """
    html_template = Template("""
    <!DOCTYPE html>
    <html>
    <head>
        <style>
            body { font-family: Arial, sans-serif; margin: 0; padding: 20px; }
            .container { max-width: 600px; margin: 0 auto; }
            .header { background-color: #2c3e50; color: white; padding: 20px; text-align: center; }
            .content { padding: 20px; background-color: #f8f9fa; }
            .order-details { background-color: white; padding: 15px; margin: 15px 0; border-radius: 5px; }
            .item { border-bottom: 1px solid #eee; padding: 10px 0; }
            .total { font-weight: bold; font-size: 18px; color: #2c3e50; }
            .footer { text-align: center; color: #666; margin-top: 20px; }
        </style>
    </head>
    <body>
        <div class="container">
            <div class="header">
                <h1>Order Confirmation</h1>
                <p>Thank you for your booking!</p>
            </div>
            <div class="content">
                <h2>Hello {{ customer_name }}!</h2>
                <p>Your order has been confirmed. Here are the details:</p>
                
                <div class="order-details">
                    <h3>Order Information</h3>
                    <p><strong>Order Number:</strong> {{ order_number }}</p>
                    <p><strong>Order Date:</strong> {{ order_date }}</p>
                    <p><strong>Status:</strong> {{ status }}</p>
                </div>
                
                <div class="order-details">
                    <h3>Items Ordered</h3>
                    {% for item in items %}
                    <div class="item">
                        <strong>{{ item.item_name }}</strong>
                        {% if item.item_type == 'room' %}
                            <br>Check-in: {{ item.check_in_date }}
                            <br>Check-out: {{ item.check_out_date }}
                            <br>Nights: {{ item.nights }}
                        {% endif %}
                        <br>Quantity: {{ item.quantity }}
                        <br>Price: ${{ "%.2f"|format(item.total_price) }}
                    </div>
                    {% endfor %}
                    
                    <div class="total">
                        Total Amount: ${{ "%.2f"|format(total_amount) }}
                    </div>
                </div>
                
                {% if notes %}
                <div class="order-details">
                    <h3>Special Notes</h3>
                    <p>{{ notes }}</p>
                </div>
                {% endif %}
                
                <p>We will contact you soon with further details. If you have any questions, please don't hesitate to contact us.</p>
            </div>
            <div class="footer">
                <p>Best regards,<br>Hotel Management Team</p>
                <p>Contact: {{ customer_phone }} | Email: {{ admin_email }}</p>
            </div>
        </div>
    </body>
    </html>
    """)
    
    text_template = Template("""
    ORDER CONFIRMATION
    
    Hello {{ customer_name }}!
    
    Your order has been confirmed. Here are the details:
    
    Order Number: {{ order_number }}
    Order Date: {{ order_date }}
    Status: {{ status }}
    
    Items Ordered:
    {% for item in items %}
    - {{ item.item_name }}
      {% if item.item_type == 'room' %}Check-in: {{ item.check_in_date }}, Check-out: {{ item.check_out_date }}, Nights: {{ item.nights }}{% endif %}
      Quantity: {{ item.quantity }}, Price: ${{ "%.2f"|format(item.total_price) }}
    {% endfor %}
    
    Total Amount: ${{ "%.2f"|format(total_amount) }}
    
    {% if notes %}Special Notes: {{ notes }}{% endif %}
    
    We will contact you soon with further details.
    
    Best regards,
    Hotel Management Team
    """)
    
    context = {
        **order_data,
        'admin_email': ADMIN_EMAIL,
        'order_date': datetime.now().strftime('%Y-%m-%d %H:%M:%S')
    }
    
    html_content = html_template.render(**context)
    text_content = text_template.render(**context)
    
    return html_content, text_content

def render_admin_notification_email(order_data):
    """
    Render admin notification email template
    
    Args:
        order_data (dict): Order information
        
    Returns:
        tuple: (html_content, text_content)
    """
    html_template = Template("""
    <!DOCTYPE html>
    <html>
    <head>
        <style>
            body { font-family: Arial, sans-serif; margin: 0; padding: 20px; }
            .container { max-width: 600px; margin: 0 auto; }
            .header { background-color: #e74c3c; color: white; padding: 20px; text-align: center; }
            .content { padding: 20px; background-color: #f8f9fa; }
            .order-details { background-color: white; padding: 15px; margin: 15px 0; border-radius: 5px; }
            .item { border-bottom: 1px solid #eee; padding: 10px 0; }
            .total { font-weight: bold; font-size: 18px; color: #e74c3c; }
            .urgent { background-color: #fff3cd; border: 1px solid #ffeaa7; padding: 10px; border-radius: 5px; }
        </style>
    </head>
    <body>
        <div class="container">
            <div class="header">
                <h1>🚨 New Order Received</h1>
                <p>Action Required</p>
            </div>
            <div class="content">
                <div class="urgent">
                    <strong>⚠️ A new order has been placed and requires your attention!</strong>
                </div>
                
                <div class="order-details">
                    <h3>Order Information</h3>
                    <p><strong>Order Number:</strong> {{ order_number }}</p>
                    <p><strong>Order Date:</strong> {{ order_date }}</p>
                    <p><strong>Status:</strong> {{ status }}</p>
                    <p><strong>Total Amount:</strong> ${{ "%.2f"|format(total_amount) }}</p>
                </div>
                
                <div class="order-details">
                    <h3>Customer Information</h3>
                    <p><strong>Name:</strong> {{ customer_name }}</p>
                    <p><strong>Phone:</strong> {{ customer_phone }}</p>
                    <p><strong>Email:</strong> {{ customer_email }}</p>
                </div>
                
                <div class="order-details">
                    <h3>Items Ordered</h3>
                    {% for item in items %}
                    <div class="item">
                        <strong>{{ item.item_name }}</strong> ({{ item.item_type }})
                        {% if item.item_type == 'room' %}
                            <br>📅 Check-in: {{ item.check_in_date }}
                            <br>📅 Check-out: {{ item.check_out_date }}
                            <br>🌙 Nights: {{ item.nights }}
                        {% endif %}
                        <br>📦 Quantity: {{ item.quantity }}
                        <br>💰 Price: ${{ "%.2f"|format(item.total_price) }}
                    </div>
                    {% endfor %}
                    
                    <div class="total">
                        💵 Total Amount: ${{ "%.2f"|format(total_amount) }}
                    </div>
                </div>
                
                {% if notes %}
                <div class="order-details">
                    <h3>Special Notes</h3>
                    <p>{{ notes }}</p>
                </div>
                {% endif %}
                
                <p><strong>Please process this order as soon as possible.</strong></p>
            </div>
        </div>
    </body>
    </html>
    """)
    
    context = {
        **order_data,
        'order_date': datetime.now().strftime('%Y-%m-%d %H:%M:%S')
    }
    
    html_content = html_template.render(**context)
    text_content = f"New order received: {order_data['order_number']} from {order_data['customer_name']} - ${order_data['total_amount']:.2f}"
    
    return html_content, text_content

@app.route('/health', methods=['GET'])
def health_check():
    """Health check endpoint"""
    return jsonify({
        'status': 'ok',
        'message': 'Email service is running',
        'timestamp': datetime.now().isoformat()
    })

@app.route('/send-order-confirmation', methods=['POST'])
def send_order_confirmation():
    """
    Send order confirmation email to customer
    
    Expected JSON payload:
    {
        "customer_email": "customer@example.com",
        "customer_name": "John Doe",
        "customer_phone": "+1234567890",
        "order_number": "ORD-001",
        "total_amount": 199.99,
        "status": "confirmed",
        "items": [...],
        "notes": "Special requests..."
    }
    """
    try:
        order_data = request.get_json()
        
        if not order_data or not order_data.get('customer_email'):
            return jsonify({'error': 'Customer email is required'}), 400
        
        # Render email content
        html_content, text_content = render_order_confirmation_email(order_data)
        
        # Send email
        subject = f"Order Confirmation - {order_data['order_number']}"
        success = send_email(
            order_data['customer_email'],
            subject,
            html_content,
            text_content
        )
        
        if success:
            return jsonify({'message': 'Order confirmation email sent successfully'})
        else:
            return jsonify({'error': 'Failed to send email'}), 500
            
    except Exception as e:
        return jsonify({'error': f'Error processing request: {str(e)}'}), 500

@app.route('/send-admin-notification', methods=['POST'])
def send_admin_notification():
    """
    Send new order notification to admin
    
    Expected JSON payload: Same as order confirmation
    """
    try:
        order_data = request.get_json()
        
        if not order_data:
            return jsonify({'error': 'Order data is required'}), 400
        
        # Render email content
        html_content, text_content = render_admin_notification_email(order_data)
        
        # Send email to admin
        subject = f"🚨 New Order: {order_data['order_number']} - ${order_data['total_amount']:.2f}"
        success = send_email(
            ADMIN_EMAIL,
            subject,
            html_content,
            text_content
        )
        
        if success:
            return jsonify({'message': 'Admin notification email sent successfully'})
        else:
            return jsonify({'error': 'Failed to send admin notification'}), 500
            
    except Exception as e:
        return jsonify({'error': f'Error processing request: {str(e)}'}), 500

if __name__ == '__main__':
    print("🚀 Starting Email Service for Hotel E-commerce System")
    print(f"📧 Admin Email: {ADMIN_EMAIL}")
    print(f"📤 From Email: {FROM_EMAIL}")
    print(f"🌐 Server starting on http://localhost:8001")
    
    # Run the Flask app
    app.run(host='0.0.0.0', port=8001, debug=True)
