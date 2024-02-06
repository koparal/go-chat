import { useState, useContext } from 'react'
import { API_URL } from '../../constants'
import { useRouter } from 'next/router'
import { AuthContext } from '../../modules/auth'

const Index = () => {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const { setAuthenticated, setUser } = useContext(AuthContext)
  const router = useRouter()

  const submitHandler = async (e: React.SyntheticEvent) => {
    e.preventDefault()

    try {
      const res = await fetch(`${API_URL}/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password }),
      });

      const data = await res.json();
      if (res.ok) {
        const user = {
          username: data.username,
          id: data.id,
          access_token: data.access_token,
          is_admin: data.is_admin
        };

        localStorage.setItem('user_info', JSON.stringify(user));

        setUser(user);
        setAuthenticated(true);

        router.push('/');
      } else {
        throw new Error(data.error);
      }
    } catch (err) {
      alert(err);
    }
  }

  const register = () => {
    router.push('/register');
  };

  return (
      <div className='flex items-center justify-center min-w-full min-h-screen'>
        <form className='flex flex-col md:w-1/5'>
          <div className='text-3xl font-bold text-center'>
            <span className='text-blue'>login</span>
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
            login
          </button>
          <label
              className='p-3 mt-6 rounded-md bg-grey font-bold text-white'
              onClick={register}
              style={{ textAlign: 'center' }}
          >
            register
          </label>
        </form>
      </div>
  )
}

export default Index
