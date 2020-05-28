import React from 'react';
import {T} from "./Utils";
import {Row, Notification} from "@canonical/react-components";

function ImageList(props) {
    return (
        <section className="col-12">
            <Row>
                <Notification type="caution" status={T('loading-images')}>
                    {props.images.map(img => {
                        return (
                            <table className="col-12">
                                <tr>
                                    <td className="col-medium">{img.alias}</td>
                                    <td className="col-small">
                                        {img.available ?
                                            <i className="p-icon--success"></i>
                                            :
                                            <i className="p-icon--spinner"></i>
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

export default ImageList;