import React, {Component} from 'react';
import api from "./api";
import {formatError, T} from "./Utils";
import {MainTable, Link} from "@canonical/react-components";
import Build from "./Build";

class BuildList extends Component {
    constructor(props) {
        super(props)
        this.state = {
            builds: [],
            expanded: false,
            expandedContent: null,
        }
    }

    getData() {
        api.buildList().then(response => {
            this.setState({builds: response.data.records})
        })
        .catch(e => {
            console.log(formatError(e.response.data))
            this.setState({error: formatError(e.response.data), message: ''});
        })
    }

    componentDidMount() {
        this.getData()
    }

    handleBuildClick = (e) => {
        e.preventDefault()
        this.getData()
    }

    handleShowClick = (e) => {
        this.setState({expandedRow: 1})
    }

    render() {

        let data = this.state.builds.map(r => {
            return {
                columns:[
                    {content: r.name, role: 'rowheader'},
                    {content: r.repo},
                    {content: r.created},
                    {content: <Link href={'/builds/'+r.id}>{T('show')}</Link>}
                    ],
            }
        })

        return (
            <div>
                <Build onClick={this.handleBuildClick} />
                <MainTable headers={[
                {
                    content: T('name')
                }, {
                    content: T('repo')
                }, {
                    content: T('created')
                }, {
                    content: ''
                }]} rows={data} />
            </div>
        );
    }
}

export default BuildList;