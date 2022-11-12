package tmdb

import (
	"cinetodoapi/model"
	"log"
	"os"
	"strconv"
	"strings"

	tmdb "github.com/cyruzin/golang-tmdb"
)

type TmdbClient struct {
	client   *tmdb.Client
	language string
}

func NewClient(apiKey string, language string) *TmdbClient {
	client, err := tmdb.Init(apiKey)
	if err != nil {
		log.Fatal(err)
		panic("Cannot initialize TMDB client")
	}
	return &TmdbClient{client, language}
}

var Instance *TmdbClient = NewClient(os.Getenv("TMDB_API_KEY"), os.Getenv("TMDB_LANGUAGE"))

func (t *TmdbClient) SearchMovies(query string) ([]model.Movie, error) {
	res, err := t.client.GetSearchMovies(query, map[string]string{"language": t.language})
	if err != nil {
		return nil, err
	}
	ret := make([]model.Movie, 0, len(res.Results))
	for _, elem := range res.Results {
		ret = append(ret, model.Movie{
			ID:            elem.ID,
			Title:         elem.Title,
			OriginalTitle: elem.OriginalTitle,
			Year:          parseYear(elem.ReleaseDate),
			Director:      "",
		})
	}
	return ret, nil
}

func (t *TmdbClient) GetMovie(id int) (*model.Movie, error) {
	res, err := t.client.GetMovieDetails(id, map[string]string{})
	if err != nil {
		return nil, err
	}
	director, err := t.getDirector(id)
	if err != nil {
		return nil, err
	}
	return &model.Movie{
		ID:            res.ID,
		Title:         res.Title,
		OriginalTitle: res.OriginalTitle,
		Year:          parseYear(res.ReleaseDate),
		Director:      director,
	}, nil
}

func (t *TmdbClient) getDirector(id int) (name string, err error) {
	res, err := t.client.GetMovieCredits(id, map[string]string{})
	if err != nil {
		return
	}
	name = ""
	for _, elem := range res.Crew {
		if strings.ToLower(elem.Job) == "director" {
			name = elem.Name
		}
	}
	return
}

func parseYear(releaseDate string) uint16 {
	parts := strings.Split(releaseDate, "-")
	if len(parts) > 0 {
		ret, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0
		}
		return uint16(ret)
	}
	return 0
}
