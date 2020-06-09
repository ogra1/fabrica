import React, {Component} from 'react';
import {formatError, T} from "./Utils";
import {Button, MainTable} from "@canonical/react-components";
import KeysAdd from "./KeysAdd";
import api from "./api";

class KeysList extends Component {
    constructor(props) {
        super(props)
        this.state = {
            showAdd: false,
            key: {},
            name: '',
            username: '',
            data: '',
            password: '',
        }
    }

    handleCancelClick = (e) => {
        e.preventDefault()
        this.setState({showAdd: false})
    }

    handleAddClick = (e) => {
        e.preventDefault()
        this.setState({showAdd: true})
    }

    handleOnChange = (field, value) => {
        let key = this.state.key
        key[field] = value
        this.setState({key: key})
    }

    handleCreate = (e) => {
        e.preventDefault()
        let key = {name: this.state.name, username:this.state.username, data: this.state.data, password: this.state.password}

        api.keysCreate(key).then(response => {
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
                    {content: r.created},
                    {content: ''}
                ],
            }
        })

        return (
            <section>
                <div>
                    <h3 className="u-float-left">{T('key-list')}</h3>
                    <Button onClick={this.handleAddClick} className="u-float-right">
                        {T('add-key')}
                    </Button>
                </div>

                {this.state.showAdd ?
                    <KeysAdd onClick={this.handleCreate} onCancel={this.handleCancelClick}
                             onChange={this.handleOnChange}
                             name={this.state.name} username={this.state.username} data={this.state.data} password={this.state.password} />
                    :
                    ''
                }

                <MainTable headers={[
                    {
                        content: T('name'),
                        className: "col-large"
                    }, {
                        content: T('created'),
                    }, {
                        content: T('actions'),
                    }
                ]} rows={data} />
            </section>
        );
    }
}

export default KeysList;