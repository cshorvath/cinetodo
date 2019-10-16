import React from "react";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faTrash} from "@fortawesome/free-solid-svg-icons/faTrash";


function getYear(item) {
    if (item.release_date)
        return <span className="badge badge-secondary"> {item.release_date.split("-")[0]}</span>;
    return null;
}

function getOriginalTitle(item) {
    if (item.title === item.original_title)
        return null;
    return <span className="text-secondary">{item.original_title}</span>;
}

export const MovieItem = props => {
    return (
        <>
            <span
                className={props.item.seen ? "seen" : null} onClick={props.onClick}> {props.item.title} {getOriginalTitle(props.item)} </span>{getYear(props.item)}
        </>)
};

export const MovieItemWithControls = props => {
    return (
        <>
            <MovieItem item={props.item} onClick={props.onToggle}/>
            <FontAwesomeIcon className="fa-pull-right" icon={faTrash} onClick={props.onRemove}/>
        </>
    )
};
