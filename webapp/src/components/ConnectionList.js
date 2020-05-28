import React from 'react';
import {T} from "./Utils";
import {Row, Notification} from "@canonical/react-components";

function ConnectionList(props) {
    return (
        <section className="col-12">
            <Row>
                <Notification type="caution" status={T('check-connections')}>
                    {props.connections.map(img => {
                        return (
                            <table className="col-12 notification-list">
                                <tr>
                                    <td className="col-large">{img.name}</td>
                                    <td className="col-small">
                                        {img.available ?
                                            <i className="p-icon--success"></i>
                                            :
                                            <i className="p-icon--error"></i>
                                        }
                                    </td>
                                </tr>
                            </table>
                        )
                    })}
                </Notification>
            </Row>
        </section>
    );
}

export default ConnectionList;