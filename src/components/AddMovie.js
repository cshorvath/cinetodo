import React from "react";
import {connect} from "react-redux";
import {requestMovies} from "../actions/actions";
import {AsyncTypeahead} from "react-bootstrap-typeahead"


const AddMovie = props => {
    return <AsyncTypeahead
        autofocus
        options={props.searchResult.results}
        isLoading={props.isLoading}
        labelKey="title"
        id="search"
        minLength={2}
        onSearch={props.requestMovies}
        placeholder="Keresés..."
        renderMenuItemChildren={option => ResultItem(option)}
    />
};

function getFragment(resultItem) {
    if (resultItem.release_date)
        return <span>({resultItem.release_date.split("-")[0]})</span>;
    return null;
}

const ResultItem = resultItem => {
    return (<React.Fragment>
        {resultItem.title}{getFragment(resultItem)}
    </React.Fragment>)
};

function mapStateToProps(state) {
    return {
        searchResult: state.searchResult,
        isLoading: state.isSearchLoading,
        error: state.error
    }
}

export default connect(mapStateToProps, {requestMovies})(AddMovie);