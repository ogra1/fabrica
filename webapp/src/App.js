import React, {Component} from 'react';
import Header from "./components/Header";
import Home from "./components/Home";
import {parseRoute} from "./components/Utils";
import BuildLog from "./components/BuildLog";
//import createHistory from 'history/createBrowserHistory'

//const history = createHistory()

class App extends Component {
    // handleNavigation(location) {
    //     this.setState({ location: location })
    //     window.scrollTo(0, 0)
    // }

    render() {
        const r = parseRoute()

        return (
            <div>
                <Header/>

                {r.section===''? <Home/> : ''}
                {r.section==='builds'? <BuildLog buildId={r.sectionId} /> : ''}

            </div>
        );
    }
}

export default App;
