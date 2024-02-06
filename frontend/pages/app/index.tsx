import React, { useState, useRef, useContext, useEffect } from 'react';
import Chat from '../../components/chat';
import { WebsocketContext } from '../../modules/websocket';
import { useRouter } from 'next/router';
import { API_URL } from '../../constants';
import autosize from 'autosize';
import { AuthContext } from '../../modules/auth';

export type Message = {
  message: string;
  client_id: string;
  username: string;
  room_id: string;
  type: 'send' | 'receive';
};

const Index = () => {
  const [messages, setMessages] = useState<Array<Message>>([]);
  const textarea = useRef<HTMLTextAreaElement>(null);
  const { conn } = useContext(WebsocketContext);
  const [users, setUsers] = useState<Array<{ username: string }>>([]);
  const { user, authenticated } = useContext(AuthContext);
  const router = useRouter();

  useEffect(() => {
    if (conn === null) {
      router.push('/');
      return;
    }

    const roomId = conn.url.split('/')[7];
    async function getUsers() {
      try {
        const res = await fetch(`${API_URL}/rooms/clients/${roomId}`, {
          method: 'GET',
          headers: { 'Content-Type': 'application/json', 'Authorization': user?.access_token },
        });
        const data = await res.json();

        setUsers(data);
      } catch (e) {
        console.error(e);
      }
    }
    getUsers();
  }, []);

  useEffect(() => {
    if (textarea.current) {
      autosize(textarea.current);
    }

    if (conn === null) {
      router.push('/');
      return;
    }

    conn.onmessage = (message) => {
      const m: Message = JSON.parse(message.data);
      if (m.message === 'New user has joined.') {
        setUsers(prevUsers => [...prevUsers, { username: m.username }]);
      }

      m.type = user?.username === m.username ? 'send' : 'receive';
      setMessages(prevMessages => [...prevMessages, m]);
    };

    conn.onclose = () => {};
    conn.onerror = () => {};
    conn.onopen = () => {};
  }, [textarea, conn, user]);

  const sendMessage = () => {
    if (!textarea.current?.value) return;
    if (conn === null) {
      router.push('/');
      return;
    }

    conn.send(textarea.current.value);
    textarea.current.value = '';
  };

  const handleLeave = () => {
    if (conn === null) {
      router.push('/');
      return;
    }

    conn.send("The user has left from the room.");
    router.push('/');
  };

  return (
      <>
        <div className='flex flex-col w-full'>
          <div className='p-4 md:mx-6 mb-14'>
            <Chat data={messages} />
          </div>
          <div className='fixed bottom-0 mt-4 w-full'>
            <div className='flex md:flex-row px-4 py-2 bg-grey md:mx-4 rounded-md'>
              <div className='flex w-full mr-4 rounded-md border border-blue'>
              <textarea
                  ref={textarea}
                  placeholder='type your message here'
                  className='w-full h-10 p-2 rounded-md focus:outline-none'
                  style={{resize: 'none'}}
              />
              </div>
              <div className='flex items-center'>
                <button
                    className='p-2 rounded-md bg-blue text-white'
                    onClick={sendMessage}
                    disabled={!authenticated}
                >
                  Send
                </button>
              </div>
              <div className='flex items-center'>
                <button
                    className='ml-2 p-2 rounded-md bg-red text-white'
                    onClick={handleLeave}
                >
                  Leave
                </button>
              </div>
            </div>
          </div>
        </div>
      </>
  );
}

export default Index;
