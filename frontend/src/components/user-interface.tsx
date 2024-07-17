'use client'

import axios from "axios"
import { useEffect, useState } from "react"
import CardComponent from "./card-component"

interface User {
    id: number
    name: string
    email: string
}

export default function UserInterface() {
    const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8000'
    const [users, setUsers] = useState<User[]>([])
    const [newUser, setNewUser] = useState({name: '', email: ''})
    const [updateUser, setUpdateUser] = useState({ id: '', name: '', email: '' });

    useEffect(()=>{
        const fetchData = async () => {
            try{
                const response = await axios.get(`${apiUrl}/api/go/users`)
                setUsers(response.data.data.reverse())
            } catch (e) {
                console.error('Error fetching data: ',e)
            }
        }
        fetchData()
    }, [apiUrl])

    const createUser = async (e: React.FormEvent<HTMLFormElement>) =>{
        e.preventDefault()

        try{
            const response = await axios.post(`${apiUrl}/api/go/users`, newUser)
            setUsers([response.data.data, ...users])
            setNewUser({name: '', email: ''})
        } catch (e){
            console.error('Error create user', e)
        }
    }

  // Delete a user
  const deleteUser = async (userId: number) => {
    try {
      await axios.delete(`${apiUrl}/api/go/users/${userId}`);
      setUsers(users.filter((user) => user.id !== userId));
    } catch (error) {
      console.error('Error deleting user:', error);
    }
  }
  return (
    <div>
        <form onSubmit={createUser} className="flex flex-col gap-1 mb-5">
            <input placeholder="Name" value={newUser.name} onChange={(e)=>setNewUser({...newUser, name: e.target.value})} className="my-2 rounded p-2"/>
            <input placeholder="Email" value={newUser.email} onChange={(e)=>setNewUser({...newUser, email: e.target.value})} className="my-2 rounded p-2"/>
            <button type="submit" className="bg-white p-1 rounded">Create User</button>
        </form>
        <div className="flex flex-col gap-2">
            {users.map(user => (
                <div key={user.id}>
                    <CardComponent card={user} />
                    <button onClick={()=>deleteUser(user.id)} className="text-red-500">
                        Delete User
                    </button>
                </div>
            ))}
        </div>
    </div>
  )
}
