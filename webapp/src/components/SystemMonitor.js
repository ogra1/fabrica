import React from 'react';
import {MainTable} from "@canonical/react-components";
import {T} from "./Utils";

function SystemMonitor(props) {
    let data = [{columns: [
            {content: props.system.cpu.toFixed(2) + '%', className: "u-align--right"},
            {content: props.system.memory.toFixed(2) + '%', className: "u-align--right"},
            {content: props.system.disk.toFixed(2) + '%', className: "u-align--right"}]
    }]
    return (
        <MainTable
            className="col-medium u-float-right monitor"
            headers={[
                {content: T('cpu'), className: "u-align--right"},
                {content: T('memory'), className: "u-align--right"},
                {content: T('disk'), className: "u-align--right"},
            ]}
            rows={data}
        />
    );
}

export default SystemMonitor;