import React, {Component} from 'react';
import {Button, Card, Form, Input, Row} from "@canonical/react-components";
import {T} from "./Utils";

class RepoAdd extends Component {
    render() {
        return (
            <Row>
                <Card>
                    <Form>
                        <Input onChange={this.props.onChange} type="text" id="repo" placeholder="https://github.com/ogra1/fabrica.git" label="Git Repo" value={this.props.repo}/>
                        <Input onChange={this.props.onChangeBranch} type="text" id="branch" placeholder={T('branch')} label="Git Branch" value={this.props.branch}/>
                        <Button onClick={this.props.onClick} appearance="positive">{T('add')}</Button>
                        <Button onClick={this.props.onCancel} appearance="neutral">{T('cancel')}</Button>
                    </Form>
                </Card>
            </Row>
        );
    }
}

export default RepoAdd;