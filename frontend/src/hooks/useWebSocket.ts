import { useEffect, useRef, useState } from 'react'
import { WS_URL } from '../lib/utils'

interface WebSocketMessage {
  type: string
  data: Record<string, unknown>
}

export const useWebSocket = () => {
  const [messages, setMessages] = useState<WebSocketMessage[]>([])
  const [connected, setConnected] = useState(false)
  const ws = useRef<WebSocket | null>(null)

  useEffect(() => {
    const connect = () => {
      ws.current = new WebSocket(WS_URL)

      ws.current.onopen = () => {
        setConnected(true)
        console.log('WebSocket connected')
      }

      ws.current.onmessage = (event) => {
        try {
          const message: WebSocketMessage = JSON.parse(event.data)
          setMessages(prev => [...prev, message])
          
          // Show notification
          if (Notification.permission === 'granted') {
            if (message.type === 'new_blog_post') {
              new Notification('New Blog Post!', {
                body: message.data.title,
                icon: '/logo.png'
              })
            } else if (message.type === 'new_comment') {
              new Notification('New Comment!', {
                body: message.data.content.substring(0, 50),
                icon: '/logo.png'
              })
            }
          }
        } catch (err) {
          console.error('Failed to parse WebSocket message:', err)
        }
      }

      ws.current.onclose = () => {
        setConnected(false)
        console.log('WebSocket disconnected')
        // Don't automatically reconnect to avoid spam
        // setTimeout(connect, 5000)
      }

      ws.current.onerror = (error) => {
        console.log('WebSocket connection not available (optional feature)')
      }
    }

    // Request notification permission
    if ('Notification' in window && Notification.permission === 'default') {
      Notification.requestPermission()
    }

    connect()

    return () => {
      if (ws.current) {
        ws.current.close()
      }
    }
  }, [])

  return { messages, connected }
}

