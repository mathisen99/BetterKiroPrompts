/**
 * Frontend Log Collector
 * 
 * Captures JavaScript errors, React errors, and API call metrics.
 * Batches logs and sends them to the backend every 5 seconds.
 * Provides colored console output for development.
 */

export type LogLevel = 'debug' | 'info' | 'warn' | 'error'

export interface LogEntry {
  level: LogLevel
  message: string
  stack?: string
  url: string
  component?: string
  userAgent: string
  timestamp: string
  metadata?: Record<string, unknown>
}

// Console colors for different log levels
const levelColors: Record<LogLevel, string> = {
  debug: '#36D7B7', // Cyan
  info: '#3498DB',  // Blue
  warn: '#F39C12',  // Yellow
  error: '#E74C3C', // Red
}

// Only log to console in development
const isDevelopment = import.meta.env.DEV

class LogCollector {
  private buffer: LogEntry[] = []
  private flushInterval: number = 5000 // 5 seconds
  private maxBufferSize: number = 50
  private flushTimer: ReturnType<typeof setInterval> | null = null
  private isInitialized: boolean = false

  constructor() {
    // Defer initialization to avoid issues during SSR or testing
    if (typeof window !== 'undefined') {
      this.initialize()
    }
  }

  private initialize(): void {
    if (this.isInitialized) return
    this.isInitialized = true
    
    this.setupErrorHandlers()
    this.setupFlushInterval()
  }

  /**
   * Set up global error handlers for uncaught errors and unhandled rejections
   */
  private setupErrorHandlers(): void {
    // Handle uncaught JavaScript errors
    window.onerror = (message, source, line, col, error) => {
      const errorMessage = typeof message === 'string' 
        ? message 
        : 'Unknown error'
      
      this.error(
        `${errorMessage} at ${source || 'unknown'}:${line || 0}:${col || 0}`,
        error?.stack,
        'window'
      )
      
      // Return false to allow default error handling
      return false
    }

    // Handle unhandled promise rejections
    window.onunhandledrejection = (event) => {
      const reason = event.reason instanceof Error 
        ? event.reason.message 
        : String(event.reason)
      
      this.error(
        `Unhandled rejection: ${reason}`,
        event.reason instanceof Error ? event.reason.stack : undefined,
        'promise'
      )
    }
  }

  /**
   * Set up periodic flush and beforeunload handler
   */
  private setupFlushInterval(): void {
    this.flushTimer = setInterval(() => this.flush(), this.flushInterval)
    
    // Flush on page unload
    window.addEventListener('beforeunload', () => {
      this.flush()
    })

    // Also flush on visibility change (tab hidden)
    document.addEventListener('visibilitychange', () => {
      if (document.visibilityState === 'hidden') {
        this.flush()
      }
    })
  }

  /**
   * Core logging method
   */
  log(level: LogLevel, message: string, component?: string, stack?: string, metadata?: Record<string, unknown>): void {
    const entry: LogEntry = {
      level,
      message,
      component,
      stack,
      url: typeof window !== 'undefined' ? window.location.href : '',
      userAgent: typeof navigator !== 'undefined' ? navigator.userAgent : '',
      timestamp: new Date().toISOString(),
      metadata,
    }

    this.buffer.push(entry)

    // Console output with colors
    this.logToConsole(entry)

    // Flush if buffer is full
    if (this.buffer.length >= this.maxBufferSize) {
      this.flush()
    }
  }

  /**
   * Output log entry to console with colors (development only)
   */
  private logToConsole(entry: LogEntry): void {
    // Skip console output in production (errors still go to backend)
    if (!isDevelopment) return

    const color = levelColors[entry.level]
    const componentTag = entry.component ? `[${entry.component}]` : '[app]'
    
    const style = `color: ${color}; font-weight: bold`
    const resetStyle = 'color: inherit'
    
    // Format: [LEVEL] [component] message
    console.log(
      `%c[${entry.level.toUpperCase()}]%c ${componentTag} ${entry.message}`,
      style,
      resetStyle
    )

    // Log stack trace for errors
    if (entry.stack && entry.level === 'error') {
      console.log('%cStack trace:', 'color: #888', entry.stack)
    }

    // Log metadata if present
    if (entry.metadata && Object.keys(entry.metadata).length > 0) {
      console.log('%cMetadata:', 'color: #888', entry.metadata)
    }
  }

  /**
   * Log debug message
   */
  debug(message: string, component?: string, metadata?: Record<string, unknown>): void {
    this.log('debug', message, component, undefined, metadata)
  }

  /**
   * Log info message
   */
  info(message: string, component?: string, metadata?: Record<string, unknown>): void {
    this.log('info', message, component, undefined, metadata)
  }

  /**
   * Log warning message
   */
  warn(message: string, component?: string, metadata?: Record<string, unknown>): void {
    this.log('warn', message, component, undefined, metadata)
  }

  /**
   * Log error message
   */
  error(message: string, stack?: string, component?: string, metadata?: Record<string, unknown>): void {
    this.log('error', message, component, stack, metadata)
  }

  /**
   * Log API call with timing information
   */
  logApiCall(method: string, url: string, status: number, durationMs: number): void {
    const level: LogLevel = status >= 400 ? 'error' : 'info'
    const message = `${method} ${url} → ${status} (${durationMs}ms)`
    
    this.log(level, message, 'api', undefined, {
      method,
      url,
      status,
      durationMs,
    })
  }

  /**
   * Log React component error (for ErrorBoundary)
   */
  logReactError(error: Error, errorInfo: { componentStack?: string }): void {
    this.error(
      `React error: ${error.message}`,
      error.stack,
      'react',
      { componentStack: errorInfo.componentStack }
    )
  }

  /**
   * Flush buffered logs to backend
   */
  async flush(): Promise<void> {
    if (this.buffer.length === 0) return

    const logs = [...this.buffer]
    this.buffer = []

    try {
      // Use sendBeacon for reliability during page unload
      const useBeacon = typeof navigator !== 'undefined' && 
                        typeof navigator.sendBeacon === 'function' &&
                        document.visibilityState === 'hidden'

      if (useBeacon) {
        const blob = new Blob([JSON.stringify({ logs })], { type: 'application/json' })
        navigator.sendBeacon('/api/logs/client', blob)
      } else {
        await fetch('/api/logs/client', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ logs }),
        })
      }
    } catch {
      // Re-add failed logs to buffer (up to half max to prevent overflow)
      const logsToRestore = logs.slice(-Math.floor(this.maxBufferSize / 2))
      this.buffer = [...logsToRestore, ...this.buffer].slice(0, this.maxBufferSize)
    }
  }

  /**
   * Clean up resources
   */
  destroy(): void {
    if (this.flushTimer) {
      clearInterval(this.flushTimer)
      this.flushTimer = null
    }
    this.flush()
  }
}

// Export singleton instance
export const logger = new LogCollector()
