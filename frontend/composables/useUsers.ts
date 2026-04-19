export const useUsers = () => {
  const { request } = useApi(); const { getAuthHeaders } = useAuth()
  return {
    getUserProfile: (username: string) => request(`/users/${username}`, { headers: getAuthHeaders() }),
    getUserStats: (username: string) => request(`/users/${username}/stats`, { headers: getAuthHeaders() }),
    followUser: (username: string) => request(`/users/${username}/follow`, { method: 'POST', headers: getAuthHeaders() })
  }
}
