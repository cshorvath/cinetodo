import React from "react";
import {connect} from "react-redux";
import {requestMovies} from "../actions/actions";
import {AsyncTypeahead} from "react-bootstrap-typeahead"


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
        renderMenuItemChildren={option => ResultItem(option)}
    />
};

function getYear(resultItem) {
    if (resultItem.release_date)
        return <span className="badge badge-secondary"> {resultItem.release_date.split("-")[0]}</span>;
    return null;
}

function getOriginalTitle(resultItem) {
    if (resultItem.title === resultItem.original_title)
        return null;
    return <span className="text-secondary">{resultItem.original_title}</span>;
}

const ResultItem = resultItem => {
    return (<React.Fragment>
        {resultItem.title} {getOriginalTitle(resultItem)} {getYear(resultItem)}
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