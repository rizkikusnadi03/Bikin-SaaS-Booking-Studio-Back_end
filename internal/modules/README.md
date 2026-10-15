# Business Modules Directory

This directory contains business modules following the Modular Monolith architecture:

- `auth/`: Authentication & user session management (Chapter 5)
- `studio/`: Studio profile management for studio owners (Chapter 8)
- `package/`: Photo package management per studio (Chapter 9)
- `portfolio/`: Studio portfolio gallery management (Chapter 13)
- `booking/`: Booking handling & owner booking management (Chapters 11, 12)
- `public/`: Public API endpoints for visitor booking portal (Chapter 10)

Each module adheres to clear boundaries:
- `dto.go`: Request & response data transfer objects
- `handler.go`: HTTP controller layer
- `service.go`: Business logic layer
- `route.go`: Route definitions and middleware bindings
