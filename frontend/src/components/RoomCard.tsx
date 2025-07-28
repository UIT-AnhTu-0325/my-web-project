
import { useState } from 'react';
import { Room } from '@/lib/api';
import Image from 'next/image';
import Link from 'next/link';

interface RoomCardProps {
  room: Room;
}

export default function RoomCard({ room }: RoomCardProps) {
  const [imgSrc, setImgSrc] = useState(
    room.images && room.images.length > 0
      ? `/images/rooms/${room.images[0]}`
      : '/images/placeholder-room.jpg'
  );
  return (
    <div className="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow">
      <div className="relative h-48">
        {imgSrc ? (
          <Image
            src={imgSrc}
            alt={room.title}
            fill
            className="object-cover"
            onError={() => setImgSrc('/images/placeholder-room.jpg')}
          />
        ) : (
          <div className="w-full h-full bg-gray-200 flex items-center justify-center">
            <span className="text-gray-500">🏨</span>
          </div>
        )}
        <div className="absolute top-2 right-2">
          <span className="bg-blue-600 text-white px-2 py-1 rounded text-sm">
            {room.room_type}
          </span>
        </div>
      </div>
      
      <div className="p-4">
        <h3 className="text-lg font-semibold text-gray-900 mb-2">{room.title}</h3>
        <p className="text-gray-600 text-sm mb-3 line-clamp-2">{room.description}</p>
        
        <div className="flex items-center justify-between mb-3">
          <div className="text-2xl font-bold text-blue-600">
            ${room.price_per_night}
            <span className="text-sm font-normal text-gray-500">/night</span>
          </div>
          <div className="text-sm text-gray-500">
            Max {room.max_occupancy} guests
          </div>
        </div>

        {room.amenities && room.amenities.length > 0 && (
          <div className="mb-3">
            <div className="flex flex-wrap gap-1">
              {room.amenities.slice(0, 3).map((amenity, index) => (
                <span 
                  key={index}
                  className="bg-gray-100 text-gray-700 px-2 py-1 rounded text-xs"
                >
                  {amenity}
                </span>
              ))}
              {room.amenities.length > 3 && (
                <span className="text-gray-500 text-xs">
                  +{room.amenities.length - 3} more
                </span>
              )}
            </div>
          </div>
        )}

        <div className="flex space-x-2">
          <Link 
            href={`/rooms/${room.id}`}
            className="flex-1 bg-blue-600 text-white text-center py-2 px-4 rounded hover:bg-blue-700 transition-colors"
          >
            View Details
          </Link>
          <Link 
            href={`/rooms/${room.id}/book`}
            className="flex-1 bg-green-600 text-white text-center py-2 px-4 rounded hover:bg-green-700 transition-colors"
          >
            Book Now
          </Link>
        </div>
      </div>
    </div>
  );
}
