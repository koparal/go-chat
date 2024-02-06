import { useState, useContext, useEffect } from 'react'
import { API_URL } from '../../constants'
import { useRouter } from 'next/router'
import { AuthContext, UserInfo } from '../../modules/auth'


const index = () => {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const { authenticated } = useContext(AuthContext)
  const { setAuthenticated, setUser } = useContext(AuthContext)

  const router = useRouter()

  useEffect(() => {
    if (authenticated) {
      router.push('/')
      return
    }
  }, [authenticated])

  const submitHandler = async (e: React.SyntheticEvent) => {
    e.preventDefault()

    try {
      const res = await fetch(`${API_URL}/register`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password }),
      })

      const data = await res.json()
      if (res.ok) {
        const user: UserInfo = {
          username: data.username,
          id: data.id,
          access_token: data.access_token,
          is_admin: data.is_admin,
        }

        localStorage.setItem('user_info', JSON.stringify(user))

        setUser(user);
        setAuthenticated(true);

        return router.push('/')
      } else {
        throw new Error(data.error);
      }
    } catch (err) {
      alert(err)
    }
  }

  const login = () => {
    router.push('/login');
  };

  return (
      <div className='flex items-center justify-center min-w-full min-h-screen'>
        <form className='flex flex-col md:w-1/5'>
          <div className='text-3xl font-bold text-center'>
            <span className='text-blue'>register</span>
          </div>
          <input
              placeholder='username'
              className='p-3 mt-8 rounded-md border-2 border-grey focus:outline-none focus:border-blue'
              value={username}
              onChange={(e) => setUsername(e.target.value)}
          />
          <input
              type='password'
              placeholder='password'
              className='p-3 mt-4 rounded-md border-2 border-grey focus:outline-none focus:border-blue'
              value={password}
              onChange={(e) => setPassword(e.target.value)}
          />
          <button
              className='p-3 mt-6 rounded-md bg-blue font-bold text-white'
              type='submit'
              onClick={submitHandler}
          >
            register
          </button>
          <label
              className='p-3 mt-6 rounded-md bg-grey font-bold text-white'
              onClick={login}
              style={{ textAlign: 'center' }}
          >
            login
          </label>
        </form>
      </div>
  )
}

export default index
