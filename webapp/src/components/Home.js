import React, {Component} from 'react';
import RepoList from "./RepoList";
import BuildList from "./BuildList";
import api from "./api";
import {T, formatError} from "./Utils";
import {Row, Notification} from '@canonical/react-components'
import ImageList from "./ImageList";
import ConnectionList from "./ConnectionList";

class Home extends Component {
    constructor(props) {
        super(props)
        this.state = {
            ready: true,
            connectReady: true,
            images: [{name: 'fabrica-bionic', available: true}, {name: 'fabrica-xenial', available: false}],
            connections: [{name: 'lxd', available: true}, {name: 'system-observe', available: false}],
            repos: [
                {id:'aaa', name:'test', repo:'github.com/TestCompany/test', keyId:'a123', branch:'master', hash:'abcdef', created:'2020-05-14T19:01:34Z', modified:'2020-05-14T19:01:34Z'}
            ],
            builds: [
                {id:'bbb', name:'test', repo:'github.com/TestCompany/test', branch:'master', status:'in-progress', duration: 222, created:'2020-05-14T19:30:34Z'}
            ],
            keys: [
                {id:'a123', name:'example', username:'example'}
            ]
        }
    }

    componentDidMount() {
        this.getDataRepos()
        this.getDataBuilds()
        this.getDataConnections()
        this.getDataImages()
        this.getDataKeys()
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

    pollConnections = () => {
        // Polls every 2s
        if (!this.state.connectReady) {
            setTimeout(this.getDataConnections.bind(this), 2000);
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

    getDataConnections() {
        api.connectionList().then(response => {
            let ready = true
            response.data.records.map(r => {
                ready = ready && r.available
                return ready
            })

            this.setState({connections: response.data.records, connectReady: ready})
        })
        .catch(e => {
            console.log(formatError(e.response.data))
            this.setState({error: formatError(e.response.data), message: ''});
        })
        .finally( ()=> {
            this.pollConnections()
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

    getDataKeys() {
        api.keysList().then(response => {
            this.setState({keys: response.data.records})
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

    handleRepoDelete = (repoId, deleteBuilds) => {
        api.repoDelete(repoId, deleteBuilds).then(response => {
            this.getDataRepos()
            this.getDataBuilds()
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
                        <Row>
                            <Notification type="negative" status={T('error') + ':'}>
                                {this.state.error}
                            </Notification>
                        </Row>
                        : ''
                }
                {
                    this.state.connectReady ?
                        '' :
                        <ConnectionList connections={this.state.connections} />
                }
                {
                    this.state.ready ?
                        '' :
                        <ImageList images={this.state.images} />
                }
                <RepoList records={this.state.repos} keys={this.state.keys} onBuild={this.handleBuildClick} onCreate={this.handleRepoCreateClick} onDelete={this.handleRepoDelete} />
                <BuildList records={this.state.builds} onDelete={this.handleBuildDelete}/>
            </div>
        );
    }
}

export default Home;