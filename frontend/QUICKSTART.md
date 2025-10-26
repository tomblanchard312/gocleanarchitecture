# 🚀 Quick Start Guide

## Get Up and Running in 3 Minutes!

### Step 1: Install Dependencies
```bash
cd frontend
npm install
```

### Step 2: Create Environment File
Create a `.env` file:
```bash
# Copy the example
cp .env.example .env
```

The `.env` should contain:
```env
VITE_API_BASE_URL=http://localhost:8080
VITE_WS_URL=ws://localhost:8080/ws
```

### Step 3: Start the Backend
In a separate terminal, start the Go API:
```bash
cd ..
go run cmd/main.go
```

### Step 4: Start the Frontend
```bash
npm run dev
```

### Step 5: Open Your Browser
Navigate to: **http://localhost:5173**

## 🎨 What You'll See

- **Beautiful landing page** with blog posts
- **Dark mode** by default (toggle with button in navbar)
- **Smooth animations** and modern UI
- **Real-time indicator** showing WebSocket connection status

## 🔑 Try These Features

1. **Register an Account**:
   - Click "Sign Up" in the navbar
   - Fill in the form with your details
   - Auto-login after registration

2. **Login**:
   - Use your email or username
   - Password must be at least 8 characters

3. **Create a Post** (after login):
   - Click "Write" in the navbar
   - (Note: Create post page needs to be added - see README)

4. **Real-time Updates**:
   - Open two browser windows
   - Create a post in one
   - Watch it appear instantly in the other!

## 🎯 Next Steps

Extend the frontend with:
- Blog post detail page
- Create/Edit post page with rich text editor
- User profile page
- Admin dashboard
- Comments section
- Search functionality

See `README.md` for full documentation!

## 💡 Tips

- **Dark Mode**: Enabled by default - toggle anytime
- **Notifications**: Grant permission for desktop alerts
- **API Issues**: Make sure backend is running on port 8080
- **Hot Reload**: Changes auto-refresh in dev mode

Enjoy building! 🚀✨

