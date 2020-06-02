import React, {Component} from 'react';
import Header from "./components/Header";
import Home from "./components/Home";
import {parseRoute} from "./components/Utils";
import BuildLog from "./components/BuildLog";
import Footer from "./components/Footer";

class App extends Component {
    render() {
        const r = parseRoute()

        return (
            <div>
                <Header/>

                {r.section===''? <Home/> : ''}
                {r.section==='builds'? <BuildLog buildId={r.sectionId} /> : ''}

                <Footer />
            </div>
        );
    }
}

export default App;
