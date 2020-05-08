import React, {Component} from 'react';
import {Button, Card, Form, Input, Row} from "@canonical/react-components";
import {T, formatError} from "./Utils";
import api from './api'

class Build extends Component {
    constructor(props) {
        super(props)
        this.state = {repo: ''}
    }

    handleRepoChange = (e) => {
        e.preventDefault()
        this.setState({repo: e.target.value})
    }

    handleClick = (e) => {
        e.preventDefault()
        api.build(this.state.repo).then(response => {
            window.location.href = '/builds/' + response.data.message
        })
        .catch(e => {
            console.log(formatError(e.response.data))
            this.setState({error: formatError(e.response.data), message: ''});
        })
    }

    render() {
        return (
            <Row>
                <br />
                <Card title={T('build-title')}>
                    <Form>
                        <Input onChange={this.handleRepoChange} type="text" id="repo" placeholder="https://github.com/ogra1/fabrica.git" label="Git Repo" value={this.state.repo}/>
                        <Button onClick={this.handleClick} appearance="positive">{T('build')}</Button>
                    </Form>
                </Card>
            </Row>
        );
    }
}

export default Build;