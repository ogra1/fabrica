import React from 'react';
import {T} from "./Utils";
import {Button, Link} from "@canonical/react-components";

function RepoActions(props) {
    return (
        <div>
            <Button data-key={props.id} onClick={props.onBuild}>{T('build')}</Button>
            <Link href="" title={T("delete")} onClick={props.onDelete}>
                <img className="action" src="/static/images/delete.svg" alt={T("delete")} data-key={props.id}/>
            </Link>
        </div>
    );
}

export default RepoActions;