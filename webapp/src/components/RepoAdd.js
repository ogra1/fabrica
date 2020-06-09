import React, {Component} from 'react';
import {Button, Card, Form, Input, Row, Select} from "@canonical/react-components";
import {T} from "./Utils";

class RepoAdd extends Component {
    render() {
        let data = this.props.keys.map(k => {
            return {
                label: k.name, value: k.id
            }
        })
        data.unshift({label:'', value:''})

        return (
            <Row>
                <Card>
                    <Form>
                        <Input onChange={this.props.onChange} type="text" id="repo" placeholder="https://github.com/ogra1/fabrica.git" label={T('git-repo')} value={this.props.repo}/>
                        <Input onChange={this.props.onChangeBranch} type="text" id="branch" placeholder={T('branch')} label={T('git-branch')} value={this.props.branch}/>
                        <Select onChange={this.props.onChangeKeyId} label={T('repo-key')} name="keyId" defaultValue={this.props.keyId} options={data}/>
                        <Button onClick={this.props.onClick} appearance="positive">{T('add')}</Button>
                        <Button onClick={this.props.onCancel} appearance="neutral">{T('cancel')}</Button>
                    </Form>
                </Card>
            </Row>
        );
    }
}

export default RepoAdd;