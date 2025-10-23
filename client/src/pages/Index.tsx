import { useState, useEffect } from "react";
import { useAuth } from "@/contexts/AuthContext";
import { Header } from "@/components/Header";
import { MovieSearch } from "@/components/MovieSearch";
import { WatchList } from "@/components/WatchList";
import { Button } from "@/components/ui/button";
import { Search, List } from "lucide-react";
import { movieService, UserMovie } from "@/services/movieService";

const Index = () => {
  const [activeView, setActiveView] = useState<'search' | 'watchlist'>('search');
  const [userMovies, setUserMovies] = useState<UserMovie[]>([]);
  const { token } = useAuth();

  useEffect(() => {
    loadUserMovies();
  }, []);

  const loadUserMovies = async () => {
    if (!token) return;
    try {
      const movies = await movieService.getUserMovies(token);
      setUserMovies(movies);
    } catch (error) {
      console.error('Failed to load movies:', error);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 via-purple-900 to-black">
      <Header />
      
      {/* Navigation */}
      <div className="container mx-auto px-4 py-6">
        <div className="flex gap-4 mb-8">
          <Button
            onClick={() => setActiveView('search')}
            variant={activeView === 'search' ? 'default' : 'outline'}
            className="flex items-center gap-2 bg-purple-600 hover:bg-purple-700 border-purple-500"
          >
            <Search className="w-4 h-4" />
            Search Movies
          </Button>
          <Button
            onClick={() => setActiveView('watchlist')}
            variant={activeView === 'watchlist' ? 'default' : 'outline'}
            className="flex items-center gap-2 bg-purple-600 hover:bg-purple-700 border-purple-500"
          >
            <List className="w-4 h-4" />
            My List ({userMovies.length})
          </Button>
        </div>

        {activeView === 'search' ? (
          <MovieSearch userMovies={userMovies} onMovieAdded={loadUserMovies} />
        ) : (
          <WatchList userMovies={userMovies} onMoviesChanged={loadUserMovies} />
        )}
      </div>
    </div>
  );
};

export default Index;
