import React from 'react';
import {T} from "./Utils";
import {Link} from "@canonical/react-components"

function BuildActions(props) {
    return (
        <div>
            <Link href={'/builds/'+props.id}  title={T("view")}>
                <img src="/static/images/show.svg" alt={T("view")}/>
            </Link>
            <Link href={'/v1/builds/'+props.id+'/download'} title={T("download")}>
                <img src="/static/images/download.svg" alt={T("download")}/>
            </Link>
        </div>
    );
}

export default BuildActions;