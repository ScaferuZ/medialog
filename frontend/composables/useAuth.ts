interface User {
  id: string
  username: string
  email: string
  display_name?: string
  avatar_url?: string
  bio?: string
  is_public: boolean
  created_at: string
  updated_at: string
}

interface AuthResponse {
  user: User
  access_token: string
  refresh_token: string
  expires_in: number
}

interface LoginRequest {
  username: string
  password: string
}

interface RegisterRequest {
  username: string
  email: string
  password: string
}

export const useAuth = () => {
  const { request } = useApi()
  const config = useRuntimeConfig()
  
  // State
  const user = useState<User | null>('auth:user', () => null)
  const token = useCookie('auth:token', {
    maxAge: 60 * 60 * 24 * 7, // 7 days
    httpOnly: false, // Client-side accessible
    sameSite: 'strict'
  })
  const refreshToken = useCookie('auth:refresh', {
    maxAge: 60 * 60 * 24 * 30, // 30 days
    sameSite: 'strict'
  })
  
  const isAuthenticated = computed(() => !!user.value && !!token.value)
  const isLoading = useState('auth:loading', () => false)
  const error = useState<string | null>('auth:error', () => null)

  // Set auth header for all requests
  const setAuthHeader = (authToken: string) => {
    // $fetch will use this header
    return {
      headers: {
        'Authorization': `Bearer ${authToken}`
      }
    }
  }

  // Login
  const login = async (data: LoginRequest) => {
    isLoading.value = true
    error.value = null
    
    try {
      const response = await request<AuthResponse>('/api/auth/login', {
        method: 'POST',
        body: data
      })
      
      // Store tokens and user
      token.value = response.access_token
      refreshToken.value = response.refresh_token
      user.value = response.user
      
      return response
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Login failed'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  // Register
  const register = async (data: RegisterRequest) => {
    isLoading.value = true
    error.value = null
    
    try {
      const response = await request<AuthResponse>('/api/auth/register', {
        method: 'POST',
        body: data
      })
      
      // Store tokens and user
      token.value = response.access_token
      refreshToken.value = response.refresh_token
      user.value = response.user
      
      return response
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Registration failed'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  // Logout
  const logout = () => {
    user.value = null
    token.value = null
    refreshToken.value = null
    error.value = null
  }

  // Fetch current user
  const fetchUser = async () => {
    if (!token.value) return null
    
    try {
      const response = await request<{ user: User }>('/api/auth/me', {
        headers: {
          'Authorization': `Bearer ${token.value}`
        }
      })
      user.value = response.user
      return response.user
    } catch (err) {
      // Token might be expired, try to refresh
      await tryRefreshToken()
      return null
    }
  }

  // Try to refresh token
  const tryRefreshToken = async () => {
    if (!refreshToken.value) {
      logout()
      return false
    }
    
    try {
      const response = await request<AuthResponse>('/api/auth/refresh', {
        method: 'POST',
        body: { refresh_token: refreshToken.value }
      })
      
      token.value = response.access_token
      refreshToken.value = response.refresh_token
      user.value = response.user
      return true
    } catch {
      logout()
      return false
    }
  }

  // Initialize auth state on app mount
  const initAuth = async () => {
    if (token.value) {
      await fetchUser()
    }
  }

  // Make token available for API calls
  const getAuthHeaders = () => {
    if (!token.value) return {}
    return {
      'Authorization': `Bearer ${token.value}`
    }
  }

  return {
    user,
    token,
    isAuthenticated,
    isLoading,
    error,
    login,
    register,
    logout,
    fetchUser,
    tryRefreshToken,
    initAuth,
    getAuthHeaders
  }
}
