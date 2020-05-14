import React, {Component} from 'react';
import api from "./api";
import {formatError, T} from "./Utils";
import RepoAdd from "./RepoAdd";
import {MainTable, Row, Button} from "@canonical/react-components";

class RepoList extends Component {
    constructor(props) {
        super(props)
        this.state = {
            showAdd: false,
            repo: '',
        }
    }

    handleAddClick = (e) => {
        e.preventDefault()
        this.setState({showAdd: true})
    }

    handleCancelClick = (e) => {
        e.preventDefault()
        this.setState({showAdd: false})
    }

    handleRepoChange = (e) => {
        e.preventDefault()
        this.setState({repo: e.target.value})
    }

    handleRepoCreate = (e) => {
        e.preventDefault()
        api.repoCreate(this.state.repo).then(response => {
            this.props.onCreate()
            this.setState({error:'', showAdd: false})
        })
        .catch(e => {
            console.log(formatError(e.response.data))
            this.setState({error: formatError(e.response.data), message: ''});
        })
    }

    render() {
        let data = this.props.records.map(r => {
            return {
                columns:[
                    {content: r.name, role: 'rowheader'},
                    {content: r.repo},
                    {content: r.hash},
                    {content: r.created},
                    {content: r.modified},
                    {content: <Button data-key={r.id} onClick={this.props.onBuild}>{T('build')}</Button>}
                    ],
            }
        })

        return (
            <section>
                <Row>
                    <div>
                        <h3 className="u-float-left">{T('repo-list')}</h3>
                        <Button onClick={this.handleAddClick} className="u-float-right">
                            {T('add-repo')}
                        </Button>
                    </div>
                    {this.state.showAdd ?
                        <RepoAdd onClick={this.handleRepoCreate} onCancel={this.handleCancelClick} onChange={this.handleRepoChange} repo={this.state.repo}/>
                        :
                        ''
                    }
                    <MainTable headers={[
                    {
                        content: T('name')
                    }, {
                        content: T('repo'),
                        className: "col-large"
                    }, {
                        content: T('last-commit'),
                    }, {
                        content: T('created'),
                    }, {
                        content: T('modified'),
                    }, {
                        content: T('actions'),
                    }
                    ]} rows={data} />
                </Row>
            </section>
        );
    }
}

export default RepoList;