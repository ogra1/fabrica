import React from 'react';
import {Spinner} from "@canonical/react-components";


function getIcon(props) {
    if (props.status==='complete') {
        return <i className="p-icon--success back-one"></i>
    } else if (props.status==='failed') {
        return <i className="p-icon--error back-one"></i>
    } else {
        return <Spinner />
    }
}


function BuildStatus(props) {
    return (
        <div className="u-align--center">
            {getIcon(props)}
        </div>
    );
}

export default BuildStatus;