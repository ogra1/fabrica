import React, {Component} from 'react';
import api from "./api";
import {formatError, T} from "./Utils";
import {MainTable, Row} from "@canonical/react-components";
import Build from "./Build";
import BuildStatus from "./BuildStatus";
import BuildActions from "./BuildActions";

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

    render() {

        let data = this.state.builds.map(r => {
            return {
                columns:[
                    {content: r.name, role: 'rowheader'},
                    {content: r.repo},
                    {content: r.created},
                    {content: <BuildStatus status={r.status} />},
                    {content: <BuildActions id={r.id} download={r.download}/>}
                    ],
            }
        })

        return (
            <div>
                <Build onClick={this.handleBuildClick} />
                <Row>
                    <MainTable headers={[
                    {
                        content: T('name')
                    }, {
                        content: T('repo'),
                        className: "col-large"
                    }, {
                        content: T('created'),
                    }, {
                        content: T('status'),
                        className: "u-align--center col-small"
                    }, {
                        content: T('actions'),
                        className: "u-align--center"
                    }]} rows={data} />
                </Row>
            </div>
        );
    }
}

export default BuildList;