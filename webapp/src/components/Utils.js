import Messages from './Messages'

export function T(message) {
    const msg = Messages[message] || message;
    return msg
}

export function formatError(data) {
    let message = T(data.code);
    if (data.message) {
        message += ': ' + data.message;
    }
    return message;
}

// URL is in the form:
//  /section
//  /section/sectionId
//  /section/sectionId/subsection
export function parseRoute() {
    const parts = window.location.pathname.split('/')

    switch (parts.length) {
        case 2:
            return {section: parts[1]}
        case 3:
            return {section: parts[1], sectionId: parts[2]}
        case 4:
            return {section: parts[1], sectionId: parts[2], subsection: parts[3]}
        default:
            return {}
    }
}