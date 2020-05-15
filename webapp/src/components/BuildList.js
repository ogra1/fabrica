import React, {Component} from 'react';
import {T} from "./Utils";
import {MainTable, Row} from "@canonical/react-components";
import BuildStatus from "./BuildStatus";
import BuildActions from "./BuildActions";
import moment from "moment";

class BuildList extends Component {
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
                    {content: r.created},
                    {content: <BuildStatus status={r.status} />},
                    {content: dur, className: "col-medium u-align--right"},
                    {content: <BuildActions id={r.id} download={r.download}/>, className: "u-align--center"}
                    ],
            }
        })

        return (
            <section>

                <Row>
                    <h3>{T('build-requests')}</h3>
                    <MainTable headers={[
                    {
                        content: T('name')
                    }, {
                        content: T('repo'), className: "col-large"
                    }, {
                        content: T('created'),
                    }, {
                        content: T('status'), className: "u-align--center col-small"
                    }, {
                        content: T('duration'), className: "col-medium u-align--right",
                    }, {
                        content: T('actions'), className: "u-align--center"
                    }]} rows={data} />
                </Row>
            </section>
        );
    }
}

export default BuildList;