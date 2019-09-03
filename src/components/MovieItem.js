import React from "react";


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

const MovieItem = props => {
    return (<React.Fragment>
        {props.item.title} {getOriginalTitle(props.item)} {getYear(props.item)}
    </React.Fragment>)
};

export default MovieItem;