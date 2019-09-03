import React from 'react';
import AddMovie from "./components/AddMovie";
import MovieList from "./components/MovieList";
import Col from "react-bootstrap/Col";


function App() {
    return (
        <React.Fragment>
                <Col>
                    <AddMovie/>
                </Col>
                <Col>
                    <MovieList/>
                </Col>
        </React.Fragment>
    );
}

export default App;
