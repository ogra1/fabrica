import React from 'react';
import {MainTable} from "@canonical/react-components";

function SystemMonitor(props) {
    let data = [props.cpu, props.memory + ' Mb']
    return (
        <MainTable
            rows={data}
        />
    );
}

export default SystemMonitor;