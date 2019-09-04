import React from 'react';
import AddMovie from "./components/AddMovie";
import MovieList from "./components/MovieList";
import Row from "react-bootstrap/Row";
import Col from "react-bootstrap/Col";


function App() {
    return (
        <React.Fragment>
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
        </React.Fragment>
    );
}

export default App;
