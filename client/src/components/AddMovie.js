import React from "react";
import {connect} from "react-redux";
import {addMovieWithDetails, requestMovies} from "../actions/actions";
import {AsyncTypeahead} from "react-bootstrap-typeahead"
import {MovieSearchItem} from "./MovieItem";


const AddMovie = props => {
    return <AsyncTypeahead
        autofocus
        options={props.searchResult.results}
        isLoading={props.isLoading}
        filterBy={(option) => true}
        labelKey="title"
        id="search"
        minLength={2}
        onSearch={props.requestMovies}
        placeholder="Keresés..."
        renderMenuItemChildren={item => (<MovieSearchItem  item={item}/>)}
        onChange={items => {
            if (items.length) {
                props.addMovieWithDetails(items[0]);
            }
        }}
    />
};


function mapStateToProps(state) {
    return {
        searchResult: state.searchResult,
        isLoading: state.isSearchLoading,
        error: state.error
    }
}

export default connect(mapStateToProps, {requestMovies, addMovieWithDetails})(AddMovie);