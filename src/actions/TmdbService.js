import axios from "axios";
import CONFIG from "../config";

class TmdbService {

    URL = "https://api.themoviedb.org/3/";

    constructor(apiKey, languages) {
        this._apiKey = apiKey;
        this._languages = languages;
    }

    async search(query) {
        const response = await this.apiRequest(
            "search/movie",
            {
                query,
                language: this._languages.join(","),
            });
        return response.data
    }

    async apiRequest(path, params) {
        return axios.get(this.URL + path,
            {
                params: {...params, apiKey: this._apiKey}
            });
    }
}

export default new TmdbService(CONFIG.TMDB_API_KEY, CONFIG.TMDB_LANGUAGES);