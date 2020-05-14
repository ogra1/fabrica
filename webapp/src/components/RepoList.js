import React, {Component} from 'react';
import api from "./api";
import {formatError, T} from "./Utils";
import RepoAdd from "./RepoAdd";
import {MainTable, Row, Button} from "@canonical/react-components";

class RepoList extends Component {
    constructor(props) {
        super(props)
        this.state = {
            records: [
                //{id:'aaa', name:'test', repo:'github.com/TestCompany/test', hash:'abcdef', created:'2020-05-14T19:01:34Z', modified:'2020-05-14T19:01:34Z'}
                ],
            showAdd: false,
        }
    }

    getData() {
        api.repoList().then(response => {
            this.setState({records: response.data.records})
        })
        .catch(e => {
            console.log(formatError(e.response.data))
            this.setState({error: formatError(e.response.data), message: ''});
        })
    }

    componentDidMount() {
        this.getData()
    }

    handleAddClick = (e) => {
        e.preventDefault()
        this.setState({showAdd: true})
    }

    handleCancelClick = (e) => {
        e.preventDefault()
        this.setState({showAdd: false})
    }

    handleRepoAddClick = (e) => {
        e.preventDefault()
        this.getData()
    }

    handleBuildClick = (e) => {
        e.preventDefault()
        let repoId = e.target.getAttribute('data-key')

        api.build(repoId).then(response => {
            this.props.onClick()
        })
        .catch(e => {
            console.log(formatError(e.response.data))
            this.setState({error: formatError(e.response.data), message: ''});
        })
    }

    render() {
        let data = this.state.records.map(r => {
            return {
                columns:[
                    {content: r.name, role: 'rowheader'},
                    {content: r.repo},
                    {content: r.hash},
                    {content: r.created},
                    {content: r.modified},
                    {content: <Button data-key={r.id} onClick={this.handleBuildClick}>{T('build')}</Button>}
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
                        <RepoAdd onClick={this.handleRepoAddClick} onCancel={this.handleCancelClick}/>
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