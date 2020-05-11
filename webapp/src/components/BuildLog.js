import React, {Component} from 'react';
import api from "./api";
import {formatError, T} from "./Utils";
import {Row, Code} from '@canonical/react-components'
import DetailsCard from "./DetailsCard";

class BuildLog extends Component {
    constructor(props) {
        super(props)
        this.state = {
            build: {},
        }
    }

    componentDidMount() {
        this.getData()
    }

    poll = () => {
        // Polls every 0.5s
        setTimeout(this.getData.bind(this), 500);
    }

    getData() {
        api.buildGet(this.props.buildId).then(response => {
            this.setState({build: response.data.record})
        })
        .catch(e => {
            console.log(formatError(e.response.data))
            this.setState({error: formatError(e.response.data), message: ''});
        })
        .finally( ()=> {
            this.poll()
        })
    }

    renderLog() {
        if (!this.state.build.logs) {return T('getting-ready')+ '\r\n'}

        return this.state.build.logs.map(l => {
            return l.message + '\r\n'
        })
    }

    render() {
        return (
            <Row>
                <h3>{T('build-log')}</h3>
                <Row>
                    <DetailsCard fields={[
                        {label: T('name'), value: this.state.build.name},
                        {label: T('repo'), value: this.state.build.repo},
                        {label: T('created'), value: this.state.build.created},
                        {label: T('status'), value: this.state.build.status},
                        ]} />
                    <Code>
                        {this.renderLog()}
                    </Code>
                </Row>
            </Row>
        );
    }
}

export default BuildLog;