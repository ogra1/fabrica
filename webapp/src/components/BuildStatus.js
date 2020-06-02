import React from 'react';


function getIcon(props) {
    if (props.status==='complete') {
        return <i className="p-icon--success back-one"></i>
    } else if (props.status==='failed') {
        return <i className="p-icon--error back-one"></i>
    } else {
        return <i className="p-icon--spinner back-one"></i>
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