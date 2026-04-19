export const useMedia = () => {
  const { request } = useApi(); const { getAuthHeaders } = useAuth()

  return {
    listMedia: () => request('/api/media', { headers: getAuthHeaders() }),
    getMedia: (id: string) => request(`/api/media/${id}`, { headers: getAuthHeaders() }),
    searchTMDB: (q: string) => request(`/api/tmdb/search?q=${encodeURIComponent(q)}`, { headers: getAuthHeaders() }),
    getPopularTMDB: () => request('/api/tmdb/popular', { headers: getAuthHeaders() }),
    syncTMDBMovie: (tmdbId: string | number) => request(`/api/tmdb/sync/${tmdbId}`, {
      method: 'POST',
      headers: getAuthHeaders()
    })
  }
}
