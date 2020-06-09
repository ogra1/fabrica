import React, {Component} from 'react';
import api from "./api";
import {formatError, T} from "./Utils";
import RepoAdd from "./RepoAdd";
import {MainTable, Row, Button} from "@canonical/react-components";
import RepoActions from "./RepoActions";
import RepoDelete from "./RepoDelete";

class RepoList extends Component {
    constructor(props) {
        super(props)
        this.state = {
            showAdd: false,
            repo: '',
            branch: 'master',
            keyId: '',
            showDelete: false,
            delete: {deleteBuilds:false},
        }
    }

    handleAddClick = (e) => {
        e.preventDefault()
        this.setState({showAdd: true})
    }

    handleCancelClick = (e) => {
        e.preventDefault()
        this.setState({showAdd: false, showDelete: false, delete: {deleteBuilds:false}})
    }

    handleDeleteBuildsClick = (e) => {
        e.preventDefault()
        let del = this.state.delete
        del.deleteBuilds = !del.deleteBuilds
        this.setState({delete: del})
    }

    handleDeleteClick = (e) => {
        e.preventDefault()
        let id = e.target.getAttribute('data-key')

        let rr = this.props.records.filter(r => {
            return r.id === id
        })

        let del = this.state.delete
        del.id = id
        del.repo = rr[0].repo
        this.setState({showDelete: true, delete: del})
    }

    handleDeleteDo = (e) => {
        e.preventDefault()

        this.props.onDelete(this.state.delete.id, this.state.delete.deleteBuilds)
        this.setState({showDelete: false, delete: {deleteBuilds:false}})
    }

    handleRepoChange = (e) => {
        e.preventDefault()
        this.setState({repo: e.target.value})
    }

    handleBranchChange = (e) => {
        e.preventDefault()
        this.setState({branch: e.target.value})
    }

    handleKeyIdChange = (e) => {
        e.preventDefault()
        this.setState({keyId: e.target.value})
    }

    handleRepoCreate = (e) => {
        e.preventDefault()
        api.repoCreate(this.state.repo, this.state.branch, this.state.keyId).then(response => {
            this.props.onCreate()
            this.setState({error:'', showAdd: false, repo:''})
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
                    {content: r.branch},
                    {content: r.hash},
                    {content: r.created},
                    {content: r.modified},
                    {content: <RepoActions id={r.id} onBuild={this.props.onBuild} onDelete={this.handleDeleteClick} />}
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
                        <RepoAdd onClick={this.handleRepoCreate} onCancel={this.handleCancelClick}
                                 onChange={this.handleRepoChange} onChangeBranch={this.handleBranchChange} onChangeKeyId={this.handleKeyIdChange}
                                 repo={this.state.repo} branch={this.state.branch} keyId={this.state.keyId} keys={this.props.keys} />
                        :
                        ''
                    }
                    {this.state.showDelete ?
                        <RepoDelete onCancel={this.handleCancelClick} onConfirm={this.handleDeleteDo} onDeleteBuilds={this.handleDeleteBuildsClick} deleteBuilds={this.state.delete.deleteBuilds} message={this.state.delete.repo} />
                        : ''
                    }
                    <MainTable headers={[
                    {
                        content: T('name')
                    }, {
                            content: T('repo'),
                            className: "col-medium"
                    }, {
                        content: T('branch'),
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