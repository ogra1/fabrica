import React, {Component} from 'react';
import RepoList from "./RepoList";
import BuildList from "./BuildList";
import api from "./api";
import {formatError} from "./Utils";

class Home extends Component {
    constructor(props) {
        super(props)
        this.state = {
            repos: [
                {id:'aaa', name:'test', repo:'github.com/TestCompany/test', hash:'abcdef', created:'2020-05-14T19:01:34Z', modified:'2020-05-14T19:01:34Z'}
            ],
            builds: [
                {id:'bbb', name:'test', repo:'github.com/TestCompany/test', status:'complete', duration: 30, created:'2020-05-14T19:30:34Z'}
            ],
        }
    }

    componentDidMount() {
        this.getDataRepos()
        this.getDataBuilds()
    }

    getDataRepos() {
        api.repoList().then(response => {
            this.setState({repos: response.data.records})
        })
        .catch(e => {
            console.log(formatError(e.response.data))
            this.setState({error: formatError(e.response.data), message: ''});
        })
    }

    getDataBuilds() {
        api.buildList().then(response => {
            this.setState({builds: response.data.records})
        })
        .catch(e => {
            console.log(formatError(e.response.data))
            this.setState({error: formatError(e.response.data), message: ''});
        })
    }

    handleRepoCreateClick = () => {
        this.getDataRepos()
    }

    handleBuildClick = (e) => {
        e.preventDefault()
        let repoId = e.target.getAttribute('data-key')

        api.build(repoId).then(response => {
            this.getDataBuilds()
        })
        .catch(e => {
            console.log(formatError(e.response.data))
            this.setState({error: formatError(e.response.data), message: ''});
        })
    }

    render() {
        return (
            <div>
                <RepoList records={this.state.repos} onBuild={this.handleBuildClick} onCreate={this.handleRepoCreateClick}/>
                <BuildList records={this.state.builds}/>
            </div>
        );
    }
}

export default Home;