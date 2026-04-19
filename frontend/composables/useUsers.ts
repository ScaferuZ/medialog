export const useUsers = () => {
  const { request } = useApi(); const { getAuthHeaders } = useAuth()
  return {
    getUserProfile: (username: string) => request(`/api/users/${username}`, { headers: getAuthHeaders() }),
    getUserStats: (username: string) => request(`/api/users/${username}/stats`, { headers: getAuthHeaders() }),
    followUser: (username: string) => request(`/api/users/${username}/follow`, { method: 'POST', headers: getAuthHeaders() })
  }
}
