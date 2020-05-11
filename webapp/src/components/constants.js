const API_PREFIX = '/v1/'

function getBaseURL() {
    return window.location.protocol + '//' + window.location.hostname + ':' + window.location.port + API_PREFIX;
}

let Constants = {
    baseUrl: getBaseURL(),
}

export default Constants
