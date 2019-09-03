import TmdbService from "./TmdbService";

export const ADD_MOVIE = "ADD_MOVIE";
export const REMOVE_MOVIE = "REMOVE_MOVIE";
export const REQUEST_MOVIES = "REQUEST_MOVIES";
export const RECEIVE_MOVIES = "RECEIVE_MOVIES";
export const ERROR = "ERROR";


export function addMovie(item) {
    return {
        type: ADD_MOVIE,
        item,
    }
}

export function removeMovie(id) {
    return {
        type: REMOVE_MOVIE,
        id
    }
}


function fetchMovies(query) {
    return dispatch => {
        if (!query.length) {
            dispatch(receiveMovies({}));
            return;
        }
        TmdbService.search(query)
            .then(
                value => dispatch(receiveMovies(value)),
                error => dispatch(moviesFetchError(error))
            );
    }
}

export function requestMovies(query) {
    return dispatch => {
        dispatch({type: REQUEST_MOVIES});
        dispatch(fetchMovies(query));
    }
}

export function receiveMovies(response) {
    return {
        type: RECEIVE_MOVIES,
        response
    }
}

export function moviesFetchError(error) {
    return {
        type: ERROR,
        error
    }
}