import {ERROR, RECEIVE_MOVIES, REQUEST_MOVIES} from "../actions/actions";

export function rootReducer(state, action) {
    switch (action.type) {
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