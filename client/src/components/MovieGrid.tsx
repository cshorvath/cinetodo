
import { MovieCard } from "@/components/MovieCard";
import { sampleMovies } from "@/data/sampleMovies";

interface MovieGridProps {
  watchList: number[];
  onAddToWatchList: (movieId: number) => void;
  onRemoveFromWatchList: (movieId: number) => void;
}

export const MovieGrid = ({
  watchList,
  onAddToWatchList,
  onRemoveFromWatchList,
}: MovieGridProps) => {
  return (
    <div className="space-y-8">
      <div>
        <h2 className="text-2xl font-bold text-white mb-2">Trending Now</h2>
        <p className="text-gray-400">Discover the most popular movies right now</p>
      </div>
      
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
        {sampleMovies.map((movie) => (
          <div key={movie.id} className="animate-fade-in">
            <MovieCard
              movie={movie}
              isInWatchList={watchList.includes(movie.id)}
              onAddToWatchList={onAddToWatchList}
              onRemoveFromWatchList={onRemoveFromWatchList}
            />
          </div>
        ))}
      </div>
    </div>
  );
};
