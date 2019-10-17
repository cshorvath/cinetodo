import React from "react";
import {connect} from "react-redux";
import {MovieItem} from "./MovieItem";
import {removeMovie, toggleMovie} from "../actions/actions";
import Table from "react-bootstrap/Table";
import Row from "react-bootstrap/Row";
import {faEye} from "@fortawesome/free-solid-svg-icons/faEye";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faStar} from "@fortawesome/free-solid-svg-icons";

const MovieTable = ({title, icon, items, removeMovie, toggleMovie}) => {
    if (items.length)
        return <Row>
            <h4><FontAwesomeIcon icon={icon} /> {title}</h4>
            <Table hover>
                <colgroup>
                    <col className="col-md-auto"/>
                    <col className="col-md-auto"/>
                    <col className="col-md-auto"/>
                    <col className="col-md-1"/>
                    <col className="col-md-1"/>
                </colgroup>
                <tbody>
                {items.sort(compareItem).map(item =>
                    <MovieItem
                        item={item}
                        onRemove={() => removeMovie(item.id)}
                        onToggle={() => toggleMovie(item.id)}
                    />
                )}
                </tbody>
            </Table>
        </Row>;
    return null;
};

const MovieList = ({items, removeMovie, toggleMovie}) => {
    const unseen = items.filter(i => !i.seen);
    const seen = items.filter(i => i.seen);
    return <>
        <MovieTable items={unseen} icon={faStar} title="Új filmek" removeMovie={removeMovie} toggleMovie={toggleMovie}/>
        <MovieTable items={seen} icon={faEye} title="Látott filmek" removeMovie={removeMovie}
                    toggleMovie={toggleMovie}/>
    </>
};

function mapStateToProps(state) {
    return {
        items: [...state.items.values()]
    }
}

function compareItem(a, b) {
    return a.title.localeCompare(b.title)
}

export default connect(mapStateToProps, {removeMovie, toggleMovie})(MovieList);