import { useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { movieService, Movie, UserMovie } from '@/services/movieService';
import { Search, Plus, Check, Film } from 'lucide-react';
import { toast } from '@/hooks/use-toast';

interface MovieSearchProps {
  userMovies: UserMovie[];
  onMovieAdded: () => void;
}

export const MovieSearch = ({ userMovies, onMovieAdded }: MovieSearchProps) => {
  const [query, setQuery] = useState('');
  const [results, setResults] = useState<Movie[]>([]);
  const [isSearching, setIsSearching] = useState(false);
  const { token } = useAuth();

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!query.trim() || !token) return;

    setIsSearching(true);
    try {
      const movies = await movieService.searchMovies(query, token);
      setResults(movies);
    } catch (error) {
      toast({ title: 'Error', description: 'Failed to search movies', variant: 'destructive' });
    } finally {
      setIsSearching(false);
    }
  };

  const handleAddMovie = async (movie: Movie) => {
    if (!token) return;
    try {
      await movieService.addMovie(movie.id, token);
      toast({ title: 'Added!', description: `${movie.title} added to your list` });
      onMovieAdded();
    } catch (error) {
      toast({ title: 'Error', description: 'Failed to add movie', variant: 'destructive' });
    }
  };

  const isInList = (movieId: number) => userMovies.some(m => m.id === movieId);

  return (
    <div className="space-y-8">
      <form onSubmit={handleSearch} className="flex gap-2">
        <Input
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Search for movies..."
          className="bg-gray-800 border-gray-700 text-white"
        />
        <Button type="submit" disabled={isSearching} className="bg-purple-600 hover:bg-purple-700">
          <Search className="w-4 h-4" />
        </Button>
      </form>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
        {results.map((movie) => (
          <div key={movie.id} className="bg-gray-800 rounded-lg overflow-hidden animate-fade-in">
            <div className="w-full h-80 bg-gray-700 flex items-center justify-center">
              <Film className="w-16 h-16 text-gray-500" />
            </div>
            <div className="p-4 space-y-3">
              <h3 className="text-white font-bold text-lg line-clamp-1">{movie.title}</h3>
              <p className="text-gray-400 text-sm">{movie.year}</p>
              <p className="text-gray-300 text-sm line-clamp-2">{movie.director}</p>
              <div className="flex items-center justify-end">
                <Button
                  onClick={() => handleAddMovie(movie)}
                  disabled={isInList(movie.id)}
                  size="sm"
                  className="bg-purple-600 hover:bg-purple-700"
                >
                  {isInList(movie.id) ? <Check className="w-4 h-4" /> : <Plus className="w-4 h-4" />}
                </Button>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};
