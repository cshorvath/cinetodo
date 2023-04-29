import React from 'react';
import AddMovie from "./components/AddMovie";
import MovieList from "./components/MovieList";
import Row from "react-bootstrap/Row";
import Col from "react-bootstrap/Col";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faFilm} from "@fortawesome/free-solid-svg-icons/faFilm";
import Container from "react-bootstrap/Container";
import Navbar from "react-bootstrap/Navbar";

function Nav() {
    return (
        <Navbar bg="dark" variant="dark" className="mb-3">
            <div className="mx-auto order-0">
                <Navbar.Brand className="mx-auto"><FontAwesomeIcon icon={faFilm}/> CineTODO</Navbar.Brand>
            </div>
        </Navbar>
    )
}

function App() {
    return <>
        <Nav/>
        <Container>
            <Row className="mb-3">
                <Col>
                    <AddMovie/>
                </Col>
            </Row>
            <Row>
                <Col>
                    <MovieList/>
                </Col>
            </Row>
        </Container>
    </>;
}

export default App;
