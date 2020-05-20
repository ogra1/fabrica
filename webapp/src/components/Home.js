import React, {Component} from 'react';
import RepoList from "./RepoList";
import BuildList from "./BuildList";
import api from "./api";
import {T, formatError} from "./Utils";
import {Notification} from '@canonical/react-components'

class Home extends Component {
    constructor(props) {
        super(props)
        this.state = {
            repos: [
                {id:'aaa', name:'test', repo:'github.com/TestCompany/test', hash:'abcdef', created:'2020-05-14T19:01:34Z', modified:'2020-05-14T19:01:34Z'}
            ],
            builds: [
                {id:'bbb', name:'test', repo:'github.com/TestCompany/test', status:'complete', duration: 222, created:'2020-05-14T19:30:34Z'}
            ],
        }
    }

    componentDidMount() {
        this.getDataRepos()
        this.getDataBuilds()
    }

    poll = () => {
        // Polls every 30s
        setTimeout(this.getDataBuilds.bind(this), 30000);
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
        .finally( ()=> {
            this.poll()
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

    handleBuildDelete = (buildId) => {
        api.buildDelete(buildId).then(response => {
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
                {
                    this.state.error ?
                        <Notification type="negative" status={T('error') + ':'}>
                            {this.state.error}
                        </Notification>
                        : ''
                }
                <RepoList records={this.state.repos} onBuild={this.handleBuildClick} onCreate={this.handleRepoCreateClick}/>
                <BuildList records={this.state.builds} onDelete={this.handleBuildDelete}/>
            </div>
        );
    }
}

export default Home;