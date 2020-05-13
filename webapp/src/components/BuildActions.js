import React from 'react';
import {T} from "./Utils";

function BuildActions(props) {
    return (
        <div>
            <Link href={'/builds/'+props.id}>{T('show')}</Link>
            <Link href={'/builds/'+props.id+'/download'}>{T('download')}</Link>
        </div>
    );
}

export default BuildActions;