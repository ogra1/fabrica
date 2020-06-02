import React from 'react';
import {Row} from "@canonical/react-components";
import SystemMonitor from "./SystemMonitor";

function Footer(props) {
    return (
        <footer>
            <Row>
                <div>
                    <div>
                        <SystemMonitor system={props.system} />
                    </div>
                </div>
            </Row>
        </footer>
    );
}

export default Footer;