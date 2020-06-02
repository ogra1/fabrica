import React from 'react';
import {T} from "./Utils";
import {Row} from "@canonical/react-components";

function Footer(props) {
    return (
        <footer>
            <Row>
                <div>
                    <p className="col-6">{T('app-title')}</p>
                    <nav className="footer">
                        <ul className="p-inline-list--middot u-no-margin--bottom">
                            <li className="p-inline-list__item">
                                <a className="p-link--external"
                                   href="https://github.com/ogra1/fabrica">{T('contribute')}</a>
                            </li>
                            <li className="p-inline-list__item">
                                <a className="p-link--external"
                                   href="https://github.com/ogra1/fabrica/issues/new">{T('report-bug')}</a>
                            </li>
                        </ul>
                    </nav>
                </div>
            </Row>
        </footer>
    );
}

export default Footer;