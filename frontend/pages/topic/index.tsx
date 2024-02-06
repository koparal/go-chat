import React, {useState, useEffect, useContext} from 'react';
import { API_URL } from '../../constants';
import { useRouter } from 'next/router'
import {AuthContext} from "../../modules/auth";

const Index = () => {
    const [topics, setTopics] = useState([]);
    const [newTopicName, setNewTopicName] = useState('');
    const [parentTopicId, setParentTopicId] = useState('');
    const [updateTopicName, setUpdateTopicName] = useState('');
    const [updateTopicId, setUpdateTopicId] = useState('');
    const [updateParentTopicId, setUpdateParentTopicId] = useState('');
    const { user } = useContext(AuthContext);
    const router = useRouter()
    useEffect(() => {
        getTopics();
    }, []);

    const getTopics = async () => {
        try {
            const res = await fetch(`${API_URL}/topics`, {
                method: 'GET',
                headers: { 'Content-Type': 'application/json', 'Authorization': user.access_token },
            });
            const data = await res.json();
            if (res.ok) {
                setTopics(data);
            }
        } catch (err) {
            console.log(err);
        }
    }

    const createTopic = async () => {
        try {
            if (newTopicName.trim() === '') {
                console.error('Topic name cannot be empty.');
                return;
            }

            const res = await fetch(`${API_URL}/topics`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json', 'Authorization': user.access_token },
                body: JSON.stringify({ name: newTopicName, parent_id: parentTopicId }),
            });
            const data = await res.json();
            if (res.ok) {
                getTopics();
                setNewTopicName('');
                setParentTopicId('');
            }else {
                throw new Error(data.error);
            }
        } catch (err) {
            alert(err);
        }
    }

    const updateTopic = async () => {
        try {
            const res = await fetch(`${API_URL}/topics/${updateTopicId}`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json', 'Authorization': user.access_token },
                body: JSON.stringify({ name: updateTopicName, parent_id: updateParentTopicId }),
            });
            const data = await res.json();
            if (res.ok) {
                getTopics();
                setUpdateTopicName('');
                setUpdateTopicId('');
                setUpdateParentTopicId('');
            } else {
                throw new Error(data.error);
            }
        } catch (err) {
            alert(err);
        }
    }

    const deleteTopic = async (topicId) => {
        try {
            const res = await fetch(`${API_URL}/topics/delete/${topicId}`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json', 'Authorization': user.access_token },
            });
            if (res.ok) {
                getTopics();
            } else {
                console.error('Failed to delete topic:', res.statusText);
            }
        } catch (err) {
            console.error('Error while deleting topic:', err);
        }
    }

    const back = () => {
        router.push('/')
    }

    const isAdmin = user.is_admin;
    useEffect(() => {
        if (!isAdmin) {
            router.push('/');
        }
    }, [isAdmin]);

    return (
        <div className="container mx-auto p-8">
            <div className='flex items-center'>
                <button
                    className='p-2 mb-2 rounded-md bg-red text-white'
                    onClick={back}
                >
                    Back
                </button>
            </div>
            <h1 className="text-3xl font-bold mb-4">Manage Topics</h1>
            <div className="mb-8">
                <h2 className="text-xl font-bold mb-2">Create New Topic</h2>
                <div className="flex">
                    <input
                        type='text'
                        className="border border-gray-300 rounded-md py-2 px-4 mr-2"
                        placeholder='Topic name'
                        value={newTopicName}
                        onChange={(e) => setNewTopicName(e.target.value)}
                    />
                    <input
                        type='text'
                        className="border border-gray-300 rounded-md py-2 px-4 mr-2"
                        placeholder='Parent topic ID (optional)'
                        value={parentTopicId}
                        onChange={(e) => setParentTopicId(e.target.value)}
                    />
                    <button className="bg-blue-500 hover:bg-blue-600 text-black py-2 px-4 rounded-md"
                            onClick={createTopic}>Create
                    </button>
                </div>
            </div>
            <div className="mb-8">
                <h2 className="text-xl font-bold mb-2">Update Existing Topic</h2>
                <div className="flex">
                    <input
                        type='text'
                        className="border border-gray-300 rounded-md py-2 px-4 mr-2"
                        placeholder='Topic ID'
                        value={updateTopicId}
                        onChange={(e) => setUpdateTopicId(e.target.value)}
                    />
                    <input
                        type='text'
                        className="border border-gray-300 rounded-md py-2 px-4 mr-2"
                        placeholder='Parent topic ID (optional)'
                        value={updateParentTopicId}
                        onChange={(e) => setUpdateParentTopicId(e.target.value)}
                    />
                    <input
                        type='text'
                        className="border border-gray-300 rounded-md py-2 px-4 mr-2"
                        placeholder='Topic name'
                        value={updateTopicName}
                        onChange={(e) => setUpdateTopicName(e.target.value)}
                    />

                    <button className="bg-blue-500 hover:bg-blue-600 text-black py-2 px-4 rounded-md"
                            onClick={updateTopic}>Update
                    </button>
                </div>
            </div>
            <div>
                <h2 className="text-xl font-bold mb-2">Topics List</h2>
                <ul>
                    {topics && topics.length > 0 ? (
                        topics.map((topic) => (
                            <li key={topic.id}
                                className="flex items-center justify-between border-b border-gray-300 py-2">
                                <span>{`ID: ${topic.id}, Parent ID: ${topic.parent_id}, Name: ${topic.name}`}</span>
                                <button className="text-red-500 hover:text-red-700"
                                        onClick={() => deleteTopic(topic.id)}>Delete
                                </button>
                            </li>
                        ))
                    ) : (
                        <li className="text-gray-500">No topics available</li>
                    )}
                </ul>
            </div>
        </div>
    )
}

export default Index;
