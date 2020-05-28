import React, {Component} from 'react';
import RepoList from "./RepoList";
import BuildList from "./BuildList";
import api from "./api";
import {T, formatError} from "./Utils";
import {Notification} from '@canonical/react-components'
import ImageList from "./ImageList";

class Home extends Component {
    constructor(props) {
        super(props)
        this.state = {
            ready: true,
            images: [{alias: 'fabrica-bionic', available: true}, {alias: 'fabrica-xenial', available: false}],
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
        this.getDataImages()
    }

    poll = () => {
        // Polls every 30s
        setTimeout(this.getDataBuilds.bind(this), 30000);
    }

    pollImages = () => {
        // Polls every 2s
        if (!this.state.ready) {
            setTimeout(this.getDataImages.bind(this), 2000);
        }
    }

    getDataImages() {
        api.imageList().then(response => {
            let ready = true
            response.data.records.map(r => {
                ready = ready && r.available
                return ready
            })

            this.setState({images: response.data.records, ready: ready})
        })
        .catch(e => {
            console.log(formatError(e.response.data))
            this.setState({error: formatError(e.response.data), message: ''});
        })
        .finally( ()=> {
            this.pollImages()
        })
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
                {
                    this.state.ready ?
                        '' :
                        <ImageList images={this.state.images} />
                }
                <RepoList records={this.state.repos} onBuild={this.handleBuildClick} onCreate={this.handleRepoCreateClick}/>
                <BuildList records={this.state.builds} onDelete={this.handleBuildDelete}/>
            </div>
        );
    }
}

export default Home;