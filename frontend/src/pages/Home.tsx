import { useQuery } from '@tanstack/react-query'
import { motion } from 'framer-motion'
import { Calendar, User, MessageCircle } from 'lucide-react'
import { Link } from 'react-router-dom'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../components/ui/Card'
import { API_BASE_URL } from '../lib/utils'
import { useWebSocket } from '../hooks/useWebSocket'

interface BlogPost {
  id: string
  title: string
  content: string
  author_id: string
  created_at: string
  updated_at: string
}

export const Home = () => {
  const { connected } = useWebSocket()
  
  const { data: posts, isLoading } = useQuery<BlogPost[]>({
    queryKey: ['posts'],
    queryFn: async () => {
      const response = await fetch(`${API_BASE_URL}/blogposts`)
      if (!response.ok) throw new Error('Failed to fetch posts')
      return response.json()
    },
  })

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-glow text-4xl gradient-text">Loading...</div>
      </div>
    )
  }

  return (
    <div className="container mx-auto px-4 py-8">
      {/* Hero Section */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="text-center mb-16 relative"
      >
        <div className="absolute inset-0 gradient-bg opacity-10 blur-3xl -z-10" />
        <h1 className="text-6xl font-bold mb-4 gradient-text">
          Welcome to ModernBlog
        </h1>
        <p className="text-xl text-muted-foreground mb-8">
          A beautiful, modern blogging platform with real-time updates
        </p>
        <div className="flex items-center justify-center space-x-2">
          <div className={`h-2 w-2 rounded-full ${connected ? 'bg-green-500 animate-pulse' : 'bg-red-500'}`} />
          <span className="text-sm text-muted-foreground">
            {connected ? 'Live updates enabled' : 'Connecting...'}
          </span>
        </div>
      </motion.div>

      {/* Blog Posts Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {posts?.map((post, index) => (
          <motion.div
            key={post.id}
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: index * 0.1 }}
          >
            <Link to={`/post/${post.id}`}>
              <Card className="h-full hover:shadow-2xl transition-all duration-300 hover:-translate-y-2">
                <CardHeader>
                  <CardTitle className="line-clamp-2">{post.title}</CardTitle>
                  <CardDescription className="flex items-center space-x-4 text-xs">
                    <span className="flex items-center">
                      <Calendar className="mr-1 h-3 w-3" />
                      {new Date(post.created_at).toLocaleDateString()}
                    </span>
                    <span className="flex items-center">
                      <User className="mr-1 h-3 w-3" />
                      {post.author_id.substring(0, 8)}
                    </span>
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <p className="text-muted-foreground line-clamp-3">
                    {post.content}
                  </p>
                  <div className="mt-4 flex items-center text-sm text-muted-foreground">
                    <MessageCircle className="mr-1 h-4 w-4" />
                    <span>View comments</span>
                  </div>
                </CardContent>
              </Card>
            </Link>
          </motion.div>
        ))}
      </div>

      {!posts || posts.length === 0 && (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          className="text-center py-16"
        >
          <p className="text-2xl text-muted-foreground">
            No posts yet. Be the first to create one!
          </p>
        </motion.div>
      )}
    </div>
  )
}

