import React, {Component} from 'react';
import {T} from "./Utils";
import {MainTable, Row, Modal} from "@canonical/react-components";
import BuildStatus from "./BuildStatus";
import BuildActions from "./BuildActions";
import moment from "moment";

class BuildList extends Component {
    constructor(props) {
        super(props)
        this.state = {
            confirmDelete: false,
            delete: {},
        }
    }

    handleCancelDelete = (e) => {
        e.preventDefault()
        this.setState({confirmDelete: false, delete: {}})
    }

    handleConfirmDelete = (e) => {
        e.preventDefault()
        let id = e.target.getAttribute('data-key')
        let buildIds = this.props.records.filter(rec => {
            return rec.id === id
        })

        if (buildIds.length>0) {
            this.setState({confirmDelete: true, delete: buildIds[0]})
        }
    }

    handleDoDelete = (e) => {
        e.preventDefault()

        this.props.onDelete(this.state.delete.id)
        this.setState({confirmDelete: false, delete: {}})
    }

    renderConfirmDelete() {
        return (
                <Modal close={this.handleCancelDelete} title={T('confirm-delete')}>
                    <p>
                        {T('confirm-delete-message') + this.state.delete.name + ' (' + this.state.delete.created  +')'}
                    </p>
                    <hr />
                    <div className="u-align--right">
                        <button onClick={this.handleCancelDelete} className="u-no-margin--bottom">
                            {T('cancel')}
                        </button>
                        <button className="p-button--negative u-no-margin--bottom" onClick={this.handleDoDelete} >
                            {T('delete')}
                        </button>
                    </div>
                </Modal>
        )
    }

    render() {
        let data = this.props.records.map(r => {
            let dur =  moment.duration(r.duration,'seconds').minutes() + ' minutes'
            if (r.duration < 120) {
                dur = moment.duration(r.duration,'seconds').seconds() + ' seconds'
            }

            return {
                columns:[
                    {content: r.name, role: 'rowheader'},
                    {content: r.repo},
                    {content: r.branch},
                    {content: r.created},
                    {content: <BuildStatus status={r.status} />},
                    {content: dur, className: "col-medium u-align--left"},
                    {content: <BuildActions id={r.id} download={r.download} onConfirmDelete={this.handleConfirmDelete}/>, className: "col-medium u-align--left"}
                    ],
            }
        })

        return (
            <section>

                {this.state.confirmDelete ? this.renderConfirmDelete() : ''}

                <Row>
                    <h3>{T('build-requests')}</h3>
                    <MainTable headers={[
                    {
                        content: T('name')
                    }, {
                        content: T('repo'), className: "col-large"
                    }, {
                        content: T('branch'),
                    }, {
                        content: T('created'),
                    }, {
                        content: T('status'), className: "u-align--center col-small"
                    }, {
                        content: T('duration'), className: "col u-align--left",
                    }, {
                        content: T('actions'), className: "col u-align--left"
                    }]} rows={data} />
                </Row>
            </section>
        );
    }
}

export default BuildList;