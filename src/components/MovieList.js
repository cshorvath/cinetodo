import React from "react";
import {connect} from "react-redux";
import ListGroup from "react-bootstrap/ListGroup";
import ListGroupItem from "react-bootstrap/ListGroupItem";
import {MovieItemWithControls} from "./MovieItem";
import {removeMovie, toggleMovie} from "../actions/actions";

const MovieList = (props) => {
    return (<ListGroup>
        {props.items.sort(compareItem).map(item =>
            <ListGroupItem key={item.id} variant={item.seen ? "dark" : "default"} action={true} onClick={() => props.toggleMovie(item.id)}>
                <MovieItemWithControls
                    item={item}
                    onRemove={() => props.removeMovie(item.id)}
                />
            </ListGroupItem>
        )}
    </ListGroup>)

};

function mapStateToProps(state) {
    return {
        items: [...state.items.values()]
    }
}

function compareItem(a, b) {
    if(!!a.seen != !!b.seen) {
        return a.seen - b.seen;
    }
    return a.title.localeCompare(b.title)
}

export default connect(mapStateToProps, {removeMovie, toggleMovie})(MovieList);