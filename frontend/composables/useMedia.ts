export const useMedia = () => {
  const { request } = useApi(); const { getAuthHeaders } = useAuth()
  return {
    listMedia: () => request('/media', { headers: getAuthHeaders() }),
    getMedia: (id: string) => request(`/media/${id}`, { headers: getAuthHeaders() }),
    searchTMDB: (q: string) => request(`/tmdb/search?q=${q}`, { headers: getAuthHeaders() })
  }
}
