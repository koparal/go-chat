import React, { useState, useEffect, useContext } from 'react';
import { API_URL } from '../constants';
import { WEBSOCKET_URL } from '../constants';
import { v4 as uuidv4 } from 'uuid';
import { AuthContext } from '../modules/auth';
import { WebsocketContext } from '../modules/websocket';
import { useRouter } from 'next/router';

const Index = () => {
    const [rooms, setRooms] = useState<{ id: string; name: string }[]>([]);
    const [roomName, setRoomName] = useState('');
    const [topics, setTopics] = useState<{ id: string; name: string }[]>([]);
    const [selectedTopic, setSelectedTopic] = useState('');
    const { user, authenticated, setAuthenticated, setUser } = useContext(AuthContext);
    const { setConn } = useContext(WebsocketContext);
    const router = useRouter();
    const accessToken = user.access_token

    const getRooms = async () => {
        try {
            const res = await fetch(`${API_URL}/rooms/list`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': accessToken
                },
            });

            const data = await res.json();
            if (res.ok) {
                setRooms(data);
            }
        } catch (err) {
            console.log(err);
        }
    };

    const getTopics = async () => {
        try {
            const res = await fetch(`${API_URL}/topics`, {
                method: 'GET',
                headers: { 'Content-Type': 'application/json', 'Authorization': accessToken },
            });

            const data = await res.json();
            if (res.ok) {
                setTopics(data);
            }
        } catch (err) {
            console.log(err);
        }
    };

    useEffect(() => {
        getRooms();
        getTopics();
    }, []);

    const submitHandler = async (e: React.SyntheticEvent) => {
        e.preventDefault();

        try {
            setRoomName('');
            const res = await fetch(`${API_URL}/rooms/create`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json', 'Authorization': accessToken },
                credentials: 'include',
                body: JSON.stringify({
                    id: uuidv4(),
                    name: selectedTopic,
                }),
            });

            if (res.ok) {
                getRooms();
            }
        } catch (err) {
            console.log(err);
        }
    };

    const joinRoom = (roomId: string) => {
        const ws = new WebSocket(
            `${WEBSOCKET_URL}/rooms/join/${roomId}?userId=${user.id}&username=${user.username}`
        );
        if (ws.OPEN) {
            setConn(ws);
            router.push('/app');
            return;
        }
    };

    const manageTopic = () => {
        router.push('/topic');
    };


    const logout = async () => {
        try {
            const res = await fetch(`${API_URL}/logout`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json', 'Authorization': user.access_token },
            });
            const data = await res.json();
            if (res.ok) {
                setUser('');
                setAuthenticated(false);
                localStorage.removeItem('user_info');
                router.push('/login');
            } else {
                throw new Error(data.error);
            }
        } catch (err) {
            alert(err);
        }
    }

    const login = () => {
        router.push('/login');
    };

    return (
        <>
            <div className='my-8 px-4 md:mx-32 w-full h-full'>
                <div className='flex mt-3'>
                    <select
                        value={selectedTopic}
                        onChange={(e) => setSelectedTopic(e.target.value)}
                        className='border border-grey p-2 rounded-md focus:outline-none focus:border-blue'
                    >
                        <option value=''>Select a topic</option>
                        {topics && topics.length > 0 && topics.map((topic) => (
                            <option key={topic.id} value={topic.name}>
                                {topic.name}
                            </option>
                        ))}
                    </select>
                    <button
                        disabled={!selectedTopic}
                        className={`bg-blue border text-white rounded-md p-2 md:ml-4 ${
                            !selectedTopic ? 'opacity-50 cursor-not-allowed' : ''
                        }`}
                        onClick={submitHandler}
                    >
                        Create Chat Room
                    </button>
                    <button
                        disabled={!user.is_admin}
                        className={`bg-red border text-white rounded-md p-2 md:ml-2 ${
                            !user.is_admin ? 'opacity-50 cursor-not-allowed' : ''
                        }`}
                        onClick={manageTopic}
                    >
                        Manage Topics
                    </button>
                    {authenticated ? (
                        <button
                            className='ml-2 p-2 rounded-md bg-grey text-white'
                            onClick={logout}
                        >
                            Logout
                        </button>
                    ) : (
                        <button
                            className='ml-2 p-2 rounded-md bg-grey text-white'
                            onClick={login}
                        >
                            Login
                        </button>
                    )}
                </div>
                <div className='mt-6'>
                    <div className='font-bold'>Rooms</div>
                    <div className='grid grid-cols-1 md:grid-cols-5 gap-4 mt-6'>
                        {rooms.length > 0 ? (
                            rooms.map((room, index) => (
                                <div
                                    key={index}
                                    className='border border-blue p-4 flex items-center rounded-md w-full'
                                >
                                    <div className='w-full'>
                                        <div className='text-blue font-bold text-lg'>#{room.name}</div>
                                    </div>
                                    <div className=''>
                                        <button
                                            className='px-4 text-white bg-blue rounded-md'
                                            onClick={() => joinRoom(room.id)}
                                        >
                                            Join
                                        </button>
                                    </div>
                                </div>
                            ))
                        ) : (
                            <div>No rooms available</div>
                        )}
                    </div>
                </div>
            </div>
        </>
    );
};

export default Index;
