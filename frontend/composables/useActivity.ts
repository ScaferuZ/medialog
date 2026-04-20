export const useActivity = () => {
  const { request } = useApi()

  return {
    getLatestActivity: () => request('/api/activity/latest')
  }
}
