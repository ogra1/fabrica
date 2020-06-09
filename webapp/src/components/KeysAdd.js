import React, {Component} from 'react';
import {Button, Card, Form, Input, Row} from "@canonical/react-components";
import {T} from "./Utils";

class KeysAdd extends Component {
    onChangeName = (e) => {
        e.preventDefault()
        this.props.onChange('name', e.target.value)
    }
    onChangeUser = (e) => {
        e.preventDefault()
        this.props.onChange('username', e.target.value)
    }
    onChangeFile = (e) => {
        e.preventDefault()

        let reader = new FileReader();
        let file = e.target.files[0];

        reader.onload = (upload) => {
            this.props.onChange('data', upload.target.result.split(',')[1])
        }

        reader.readAsDataURL(file);
    }
    onChangePassword = (e) => {
        e.preventDefault()
        this.props.onChange('password', e.target.value)
    }

    render() {
        console.log(this.props.record)
        return (
            <Row>
                <Card>
                    <Form>
                        <Input onChange={this.onChangeName} type="text" id="name" placeholder={T('key-name-help')} label={T('key-name')} value={this.props.record.name}/>
                        <Input onChange={this.onChangeUser} type="text" id="username" placeholder={T('username-repo-help')} label={T('username-repo')} value={this.props.record.username}/>
                        <Input onChange={this.onChangeFile} type="file" id="privateKey" placeholder={T('file')} label={T('private-key')} />
                        <Input onChange={this.onChangePassword} type="password" id="password" placeholder={T('private-key-password-help')} label={T('private-key-password')} value={this.props.record.password}/>
                        <Button onClick={this.props.onClick} appearance="positive">{T('add')}</Button>
                        <Button onClick={this.props.onCancel} appearance="neutral">{T('cancel')}</Button>
                    </Form>
                </Card>
            </Row>
        );
    }
}

export default KeysAdd;