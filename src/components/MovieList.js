import React from "react";
import {connect} from "react-redux";
import MovieItem from "./MovieItem";
import ListGroup from "react-bootstrap/ListGroup";
import ListGroupItem from "react-bootstrap/ListGroupItem";

const MovieList = (props) => {
    return (<ListGroup>
        {props.items.map(item =>
            <ListGroupItem><MovieItem item={item}/></ListGroupItem>
        )}
    </ListGroup>)

};

function mapStateToProps(state) {
    return {
        items: [...state.items.values()]
    }
}

export default connect(mapStateToProps)(MovieList);