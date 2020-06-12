import React, {Component} from 'react';
import {formatError, T} from "./Utils";
import {Row} from '@canonical/react-components'
import api from "./api";
import KeysList from "./KeysList";

class Settings extends Component {
    constructor(props) {
        super(props)
        this.state = {
            keys: [{id:'a123', name: 'first', created:'2020-06-09T19:01:34Z'}]
        }
    }

    componentDidMount() {
        this.getDataKeys()
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

    handleCreate = () => {
        this.getDataKeys()
    }

    render() {
        return (
            <Row>
                <h2>{T('settings')}</h2>

                <KeysList records={this.state.keys} onCreate={this.handleCreate} />
            </Row>
        );
    }
}

export default Settings;