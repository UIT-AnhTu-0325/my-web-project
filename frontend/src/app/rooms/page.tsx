'use client';

import { useEffect, useState } from 'react';
import { apiClient } from '@/lib/api-client';
import { Room } from '@/types/api';
import RoomCard from '@/components/RoomCard';
import LoadingSpinner from '@/components/LoadingSpinner';

export default function RoomsPage() {
  const [rooms, setRooms] = useState<Room[]>([]);
  const [rawRooms, setRawRooms] = useState<unknown>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchRooms = async () => {
      try {
        const response = await apiClient.get('/api/rooms');
        setRawRooms(response.data);
        if (Array.isArray(response.data)) {
          setRooms(response.data);
        } else if (Array.isArray(response.data?.rooms)) {
          setRooms(response.data.rooms);
        } else {
          setRooms([]);
        }
        console.log('Rooms API response:', response.data);
      } catch (error) {
        console.error('Error fetching rooms:', error);
      } finally {
        setLoading(false);
      }
    };
    fetchRooms();
  }, []);

  if (loading) {
    return <LoadingSpinner />;
  }

  return (
    <div className="py-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center mb-12">
          <h1 className="text-4xl font-bold text-gray-900 mb-4">
            Our Rooms
          </h1>
          <p className="text-lg text-gray-600 max-w-3xl mx-auto">
            Choose from our selection of comfortable and luxurious rooms. 
            Each room is designed to provide you with the best possible experience during your stay.
          </p>
        </div>

        {Array.isArray(rooms) ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
            {rooms.map((room) => (
              <RoomCard key={room.id} room={room} />
            ))}
          </div>
        ) : (
          <div className="text-center py-12">
            <p className="text-red-500 text-lg">Rooms data is not an array. Check API response.</p>
            <pre className="text-xs text-gray-500 overflow-x-auto bg-gray-100 p-2 rounded mt-2">{JSON.stringify(rawRooms, null, 2)}</pre>
          </div>
        )}

        {Array.isArray(rooms) && rooms.length === 0 && !loading && (
          <div className="text-center py-12">
            <p className="text-gray-500 text-lg">No rooms available at the moment.</p>
          </div>
        )}
      </div>
    </div>
  );
}
