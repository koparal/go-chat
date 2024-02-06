import { useState, createContext, useEffect } from 'react'
import { useRouter } from 'next/router'

export type UserInfo = {
  id: string
  username: string
  access_token: string
  is_admin: boolean
}

export const AuthContext = createContext<{
  authenticated: boolean
  setAuthenticated: (auth: boolean) => void
  user: UserInfo
  setUser: (user: UserInfo) => void
}>({
  authenticated: false,
  user: { username: '', id: '', is_admin: false, access_token: '' },
  setAuthenticated: () => {},
  setUser: () => {},
})

const AuthContextProvider = ({ children }: { children: React.ReactNode }) => {
  const [authenticated, setAuthenticated] = useState(false)
  const [user, setUser] = useState<UserInfo>({ username: '', id: '', is_admin: false, access_token: '' })

  const router = useRouter()

  useEffect(() => {
    const userInfo = localStorage.getItem('user_info')

    if (!userInfo) {
      if (window.location.pathname !== '/register' && window.location.pathname !== '/') {
        router.push('/login');
        return;
      }
    } else {
      const user: UserInfo = JSON.parse(userInfo)
      if (user) {
        setUser({
          username: user.username,
          id: user.id,
          is_admin: user.is_admin,
          access_token: user.access_token,
        })
      }
      setAuthenticated(true)
    }
  }, [authenticated])

  return (
    <AuthContext.Provider
      value={{
        authenticated: authenticated,
        setAuthenticated: setAuthenticated,
        user: user,
        setUser: setUser,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export default AuthContextProvider
