
import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Heart, Plus, Star, Clock } from "lucide-react";
import { Movie } from "@/types/movie";

interface MovieCardProps {
  movie: Movie;
  isInWatchList: boolean;
  onAddToWatchList: (movieId: number) => void;
  onRemoveFromWatchList: (movieId: number) => void;
}

export const MovieCard = ({
  movie,
  isInWatchList,
  onAddToWatchList,
  onRemoveFromWatchList,
}: MovieCardProps) => {
  const [isHovered, setIsHovered] = useState(false);

  const handleWatchListToggle = () => {
    if (isInWatchList) {
      onRemoveFromWatchList(movie.id);
    } else {
      onAddToWatchList(movie.id);
    }
  };

  return (
    <Card 
      className="group relative overflow-hidden bg-gray-800/50 border-gray-700 hover:border-purple-500/50 transition-all duration-300 hover:scale-105 hover:shadow-2xl hover:shadow-purple-500/20"
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
      <div className="relative">
        <img
          src={movie.poster}
          alt={movie.title}
          className="w-full h-80 object-cover transition-transform duration-300 group-hover:scale-110"
        />
        
        {/* Overlay */}
        <div className={`absolute inset-0 bg-gradient-to-t from-black/80 via-transparent to-transparent transition-opacity duration-300 ${isHovered ? 'opacity-100' : 'opacity-0'}`} />
        
        {/* Action Button */}
        <Button
          onClick={handleWatchListToggle}
          size="sm"
          className={`absolute top-3 right-3 transition-all duration-300 ${
            isInWatchList 
              ? 'bg-red-500 hover:bg-red-600' 
              : 'bg-purple-500 hover:bg-purple-600'
          } ${isHovered ? 'opacity-100 scale-100' : 'opacity-0 scale-75'}`}
        >
          {isInWatchList ? (
            <Heart className="w-4 h-4 fill-current" />
          ) : (
            <Plus className="w-4 h-4" />
          )}
        </Button>

        {/* Rating Badge */}
        <Badge className="absolute top-3 left-3 bg-yellow-500/90 text-black font-semibold">
          <Star className="w-3 h-3 mr-1 fill-current" />
          {movie.rating}
        </Badge>
      </div>

      <CardContent className="p-4 space-y-3">
        <div>
          <h3 className="font-bold text-lg text-white line-clamp-1 group-hover:text-purple-300 transition-colors">
            {movie.title}
          </h3>
          <p className="text-gray-400 text-sm flex items-center gap-1">
            <Clock className="w-3 h-3" />
            {movie.year} • {movie.genre}
          </p>
        </div>
        
        <p className="text-gray-300 text-sm line-clamp-3 leading-relaxed">
          {movie.description}
        </p>
      </CardContent>
    </Card>
  );
};
