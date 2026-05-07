interface User {
  id: string
  email: string
  username: string
}

export function useAuth() {
  const user = useState<User | null>('auth.user', () => null)
  const config = useRuntimeConfig()
  const base = config.public.apiBase

  async function fetchMe() {
    try {
      const headers = import.meta.server ? useRequestHeaders(['cookie']) : undefined
      const me = await $fetch<User>(`${base}/auth/me`, {
        credentials: 'include',
        headers,
      })
      user.value = me
    } catch {
      user.value = null
    }
  }

  async function login(email: string, password: string) {
    const res = await $fetch<User>(`${base}/auth/login`, {
      method: 'POST',
      body: { email, password },
      credentials: 'include',
    })
    user.value = res
  }

  async function signup(email: string, password: string) {
    const res = await $fetch<User>(`${base}/auth/signup`, {
      method: 'POST',
      body: { email, password },
      credentials: 'include',
    })
    user.value = res
  }

  async function logout() {
    await $fetch(`${base}/auth/logout`, {
      method: 'POST',
      credentials: 'include',
    })
    user.value = null
    await navigateTo('/')
  }

  async function renameUser(username: string) {
    const res = await $fetch<User>(`${base}/auth/me/username`, {
      method: 'PATCH',
      body: { username },
      credentials: 'include',
    })
    user.value = res
    return res
  }

  return { user, fetchMe, login, signup, logout, renameUser }
}
