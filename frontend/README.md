# ModernBlog Frontend 🎨✨

A stunning, modern blog frontend built with the latest web technologies.

## 🚀 Tech Stack

- **React 18** - Modern React with hooks
- **TypeScript** - Type safety and better DX
- **Vite** - Lightning-fast build tool
- **Tailwind CSS** - Utility-first CSS framework
- **Framer Motion** - Smooth animations
- **React Router** - Client-side routing
- **TanStack Query** - Powerful data fetching
- **Lucide React** - Beautiful icons

## ✨ Features

- 🌗 **Dark/Light Mode** - Smooth theme switching
- ⚡ **Real-time Updates** - WebSocket integration for live notifications
- 📱 **Fully Responsive** - Mobile-first design
- 🎨 **Modern UI** - Glassmorphism, gradients, smooth animations
- 🔐 **JWT Authentication** - Secure login/register flow
- 💬 **Comments System** - Nested comments with real-time updates
- 🎭 **Smooth Animations** - Framer Motion powered transitions
- 🔔 **Push Notifications** - Desktop notifications for new content
- 👥 **User Profiles** - View and manage user profiles
- 👑 **Admin Dashboard** - User management for admins

## 🏗️ Installation

1. **Navigate to frontend directory**:
   ```bash
   cd frontend
   ```

2. **Install dependencies**:
   ```bash
   npm install
   ```

3. **Create environment file**:
   Create a `.env` file in the `frontend` directory:
   ```env
   VITE_API_BASE_URL=http://localhost:8080
   VITE_WS_URL=ws://localhost:8080/ws
   ```

4. **Start the development server**:
   ```bash
   npm run dev
   ```

5. **Open your browser**:
   Navigate to `http://localhost:5173`

## 🛠️ Development

### Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run lint` - Run ESLint

### Project Structure

```
frontend/
├── src/
│   ├── components/       # Reusable UI components
│   │   ├── ui/          # Base UI components (Button, Card, Input)
│   │   └── Navbar.tsx   # Navigation bar
│   ├── contexts/        # React contexts
│   │   ├── AuthContext.tsx   # Authentication state
│   │   └── ThemeContext.tsx  # Theme management
│   ├── hooks/           # Custom React hooks
│   │   └── useWebSocket.ts   # WebSocket hook
│   ├── pages/           # Page components
│   │   ├── Home.tsx     # Homepage with blog list
│   │   ├── Login.tsx    # Login page
│   │   └── Register.tsx # Registration page
│   ├── lib/            # Utilities
│   │   └── utils.ts    # Helper functions
│   ├── App.tsx         # Main app component
│   ├── main.tsx        # Entry point
│   └── index.css       # Global styles
├── public/             # Static assets
├── .env               # Environment variables
├── tailwind.config.js # Tailwind configuration
├── vite.config.ts     # Vite configuration
└── package.json       # Dependencies

```

## 🎨 Design System

### Colors

The app uses a modern color palette with support for dark and light themes:
- **Primary**: Purple gradient (262 83% 58%)
- **Background**: Dynamic based on theme
- **Accent**: Subtle hover states
- **Destructive**: Error states

### Typography

- Clean, modern fonts
- Gradient text for headings
- Proper text hierarchy

### Components

All components are built with:
- Consistent sizing and spacing
- Smooth transitions
- Hover effects
- Accessibility in mind

## 🔌 API Integration

The frontend connects to the Go backend API:

### Endpoints Used

- `POST /auth/register` - User registration
- `POST /auth/login` - User login
- `GET /blogposts` - Fetch all blog posts
- `GET /blogposts/{id}` - Fetch single post
- `POST /blogposts` - Create post (authenticated)
- `WS /ws` - WebSocket connection

### Authentication Flow

1. User logs in or registers
2. JWT token received and stored in localStorage
3. Token sent in Authorization header for protected requests
4. Auto-redirect to login on 401 responses

## 🌐 WebSocket Features

Real-time updates for:
- 📝 New blog posts
- 💬 New comments
- 🔔 Desktop notifications (with permission)

The WebSocket connection automatically reconnects if disconnected.

## 📱 Responsive Design

Breakpoints:
- Mobile: < 768px
- Tablet: 768px - 1024px
- Desktop: > 1024px

All components are mobile-first and fully responsive.

## 🚀 Production Build

1. **Build the app**:
   ```bash
   npm run build
   ```

2. **Preview locally**:
   ```bash
   npm run preview
   ```

3. **Deploy**:
   The `dist/` folder contains the production build. Deploy to:
   - Vercel
   - Netlify
   - GitHub Pages
   - Any static hosting service

### Environment Variables for Production

Update your production environment variables:
```env
VITE_API_BASE_URL=https://your-api-domain.com
VITE_WS_URL=wss://your-api-domain.com/ws
```

## 🎯 Features to Build Next

This is a foundation with the core pages. You can extend with:

- 📝 **Blog post detail page** - Full post view with comments
- ✍️ **Create/Edit post page** - Rich text editor
- 👤 **Profile page** - User profile with edit capabilities
- 👑 **Admin dashboard** - User management panel
- 🔍 **Search functionality** - Search posts
- 📊 **Analytics** - Post views, user stats
- 🖼️ **Image uploads** - Cover images for posts
- 📱 **Social sharing** - Share buttons
- 🌐 **Social login** - OAuth2 integration

## 💡 Tips

### Dark Mode

The app defaults to dark mode. Users can toggle with the button in the navbar.

### Notifications

Grant notification permissions for desktop alerts when new content is posted.

### Performance

- Images are lazy-loaded
- Code splitting for routes
- Optimized bundle size

## 🐛 Troubleshooting

### WebSocket not connecting

Check that the Go backend is running on port 8080 and WebSocket endpoint is accessible.

### API requests failing

Ensure the backend is running and `VITE_API_BASE_URL` is correct.

### Build errors

Delete `node_modules` and reinstall:
```bash
rm -rf node_modules
npm install
```

## 📄 License

MIT License - feel free to use this for your projects!

## 🙏 Credits

Built with love using modern web technologies.

---

**Happy coding! 🚀✨**
