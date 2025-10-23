import { useAuth } from '@/contexts/AuthContext';
import { Button } from '@/components/ui/button';
import { movieService, UserMovie } from '@/services/movieService';
import { Film, Trash2, Eye, EyeOff } from 'lucide-react';
import { toast } from '@/hooks/use-toast';

interface WatchListProps {
  userMovies: UserMovie[];
  onMoviesChanged: () => void;
}

export const WatchList = ({ userMovies, onMoviesChanged }: WatchListProps) => {
  const { token } = useAuth();

  const handleDelete = async (id: number) => {
    if (!token) return;
    try {
      await movieService.deleteMovie(id, token);
      toast({ title: 'Removed', description: 'Movie removed from your list' });
      onMoviesChanged();
    } catch (error) {
      toast({ title: 'Error', description: 'Failed to remove movie', variant: 'destructive' });
    }
  };

  const handleToggleSeen = async (id: number, currentSeen: boolean) => {
    if (!token) return;
    try {
      await movieService.toggleSeen(id, !currentSeen, token);
      onMoviesChanged();
    } catch (error) {
      toast({ title: 'Error', description: 'Failed to update movie', variant: 'destructive' });
    }
  };

  if (userMovies.length === 0) {
    return (
      <div className="text-center py-16">
        <div className="bg-gray-800/50 rounded-full w-24 h-24 flex items-center justify-center mx-auto mb-6">
          <Film className="w-12 h-12 text-gray-400" />
        </div>
        <h3 className="text-2xl font-bold text-white mb-2">Your list is empty</h3>
        <p className="text-gray-400 max-w-md mx-auto">
          Search for movies and add them to your list to keep track of what you want to watch.
        </p>
      </div>
    );
  }

  return (
    <div className="space-y-8">
      <div>
        <h2 className="text-2xl font-bold text-white mb-2">My List</h2>
        <p className="text-gray-400">
          {userMovies.length} movie{userMovies.length !== 1 ? 's' : ''} in your collection
        </p>
      </div>
      
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
        {userMovies.map((movie) => (
          <div key={movie.id} className="bg-gray-800 rounded-lg overflow-hidden animate-fade-in">
            <div className="w-full h-80 bg-gray-700 flex items-center justify-center">
              <Film className="w-16 h-16 text-gray-500" />
            </div>
            <div className="p-4 space-y-3">
              <h3 className="text-white font-bold text-lg line-clamp-1">{movie.title}</h3>
              <p className="text-gray-400 text-sm">{movie.year}</p>
              <p className="text-gray-300 text-sm line-clamp-2">{movie.director}</p>
              <div className="flex items-center justify-end gap-2">
                <Button
                  onClick={() => handleToggleSeen(movie.id, movie.seen)}
                  size="sm"
                  variant={movie.seen ? "default" : "outline"}
                  className="bg-green-600 hover:bg-green-700"
                >
                  {movie.seen ? <Eye className="w-4 h-4" /> : <EyeOff className="w-4 h-4" />}
                </Button>
                <Button
                  onClick={() => handleDelete(movie.id)}
                  size="sm"
                  variant="destructive"
                >
                  <Trash2 className="w-4 h-4" />
                </Button>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};
