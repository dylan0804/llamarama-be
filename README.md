# Chat

A real-time chat application built with Go, featuring WebSocket connections, Redis session management, and PostgreSQL database storage.

## 🚀 Features

- Real-time messaging using WebSocket connections
- Secure user authentication and session management with Redis
- Persistent message storage using PostgreSQL (Supabase)
- Concurrent connection handling using Go's goroutines
- RESTful API endpoints for room management

## 🛠️ Technology Stack

- **Backend**: Go (Gin framework)
- **Real-time Communication**: WebSocket
- **Session Management**: Redis
- **Database**: PostgreSQL (Supabase)
- **Authentication**: JWT tokens

## 🏗️ Architecture

- Uses goroutines for handling multiple WebSocket connections concurrently
- Implements room-based chat system
- Maintains user sessions using Redis for better scalability
- Stores chat history and user data in PostgreSQL

## 🔧 Setup

### Prerequisites

- Go 1.19 or higher
- Redis server
- PostgreSQL database (or Supabase account)

### Environment Variables

Create a `.env` file in the root directory:
