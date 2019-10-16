import {ADD_MOVIE, TOGGLE_MOVIE, ERROR, RECEIVE_MOVIES, REMOVE_MOVIE, REQUEST_MOVIES} from "../actions/actions";

export function rootReducer(state, action) {
    const newState = {
        ...state,
    };
    switch (action.type) {
        case ADD_MOVIE:
            if (!newState.items.has(action.item.id)) {
                newState.items.set(action.item.id, action.item)
            }
            return newState;
        case  REQUEST_MOVIES:
            return {
                ...state,
                isSearchLoading: true
            };
        case RECEIVE_MOVIES: {
            return {
                ...state,
                isSearchLoading: false,
                searchResult: action.response,
                error: null
            }
        }
        case REMOVE_MOVIE: {
            newState.items.delete(action.id);
            return newState;
        }
        case TOGGLE_MOVIE: {
            const movie = newState.items.get(action.id);
            if (movie) {
               movie.seen = !movie.seen;
            }
            return newState;
        }
        case ERROR: {
            return {
                ...state,
                isSearchLoading: false,
                error: action.error
            }
        }
        default:
            return state;
    }
}