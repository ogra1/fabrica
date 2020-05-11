import React, { Component } from 'react';
import {T} from "./Utils";


let links = [];

class Header extends Component {
    constructor(props) {
        super(props)
        this.state = {};
    }

    link(l) {
        if (this.props.sectionId) {
            // This is the secondary menu
            return '/' + this.props.section + '/' + this.props.sectionId + '/' + l;
        } else {
            return '/' + l;
        }
    }

    render() {
        return (
            <header id="navigation" class="p-navigation header-slim">
                <div className="p-navigation__banner row">
                    <div className="p-navigation__logo">
                        <div className="u-vertically-center">
                            <a href="/" className="p-navigation__link">
                                <img src="/static/images/fabrica.png" alt="ubuntu" />
                            </a>
                        </div>
                    </div>

                    <nav className="p-navigation__nav">
                        <span className="u-off-screen"><a href="#navigation">Jump to site</a></span>
                        <ul className="p-navigation__links" role="menu">
                            {links.map((l) => {
                                var active = '';
                                if ((this.props.section === l) || (this.props.subsection === l)) {
                                    active = ' active'
                                }
                                return (
                                    <li key={l} className={'p-navigation__link' + active} role="menuitem"><a href={this.link(l)}>{T(l)}</a></li>
                                )
                            })}
                        </ul>
                    </nav>
                </div>
            </header>
        );
    }
}

export default Header;