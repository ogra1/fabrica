import React from 'react';
import {Link, Modal} from "@canonical/react-components";
import {T} from "./Utils";

function RepoDelete(props) {
    return (
        <Modal close={props.onCancel} title={T('confirm-delete')}>
            <p>
                {T('confirm-delete-repo-message') + props.message}
            </p>
            <div>
                {props.deleteBuilds ?
                    <Link href="" title={T("delete-builds")} onClick={props.onDeleteBuilds}>
                        <img className="action" src="/static/images/check-square.svg" alt={T("delete")} data-key={props.id}/>{T('delete-builds')}
                    </Link>
                    :
                    <Link href="" title={T("delete-builds")} onClick={props.onDeleteBuilds}>
                        <img className="action" src="/static/images/square.svg" alt={T("delete")} data-key={props.id}/>{T('delete-builds')}
                    </Link>
                }
            </div>
            <hr />
            <div className="u-align--right">
                <button onClick={props.onCancel} className="u-no-margin--bottom">
                    {T('cancel')}
                </button>
                <button className="p-button--negative u-no-margin--bottom" onClick={props.onConfirm} >
                    {T('delete')}
                </button>
            </div>
        </Modal>
    );
}

export default RepoDelete;