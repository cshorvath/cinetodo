import React from "react";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faTrash} from "@fortawesome/free-solid-svg-icons/faTrash";
import Button from "react-bootstrap/Button";
import {faEye} from "@fortawesome/free-solid-svg-icons/faEye";
import ButtonGroup from "react-bootstrap/ButtonGroup";
import {faEyeSlash} from "@fortawesome/free-solid-svg-icons/faEyeSlash";


function getYear(item, parenthesize) {
    if (item.release_date) {
        const year = item.release_date.split("-")[0];
        return parenthesize ? `(${year})` : year;
    }
    return null;
}

function getOriginalTitle(item) {
    if (item.title === item.original_title)
        return null;
    return <span className="text-secondary">{item.original_title}</span>;
}

export const MovieSearchItem = props => {
    return (
        <>
            <span> {props.item.title} {getOriginalTitle(props.item)} </span> {getYear(props.item, true)}
        </>);
};

export const MovieItem = ({item, onRemove, onToggle}) => {
    const clsName = item.seen ? "seen" : null;
    return <tr>
        <td className={clsName}>{item.title}</td>
        <td className={clsName}>{getOriginalTitle(item)}</td>
        <td className={clsName}><span className="text-secondary">{item.director}</span></td>
        <td className={clsName}><span className="text-secondary">{getYear(item)}</span></td>
        <td className="text-right">
            <ButtonGroup>
                <Button size="sm" variant="dark" onClick={onToggle}><FontAwesomeIcon icon={item.seen ? faEyeSlash : faEye}/></Button>
                <Button size="sm" variant="danger" onClick={onRemove}><FontAwesomeIcon icon={faTrash}/></Button>
            </ButtonGroup>
        </td>
    </tr>
};